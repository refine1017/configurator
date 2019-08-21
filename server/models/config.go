package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"server/models/base"
	"server/models/table"
	"server/modules/mongo"
	"strings"
)

const (
	FormatAll   = "all"
	FormatTable = "table"
	FormatKV    = "kv"
	FormatJson  = "json"
)

var formatList = []string{
	FormatTable,
	FormatKV,
	FormatJson,
}

type HookEvent int

type HookFunc func(config *Config, env *Environment, admin string)

const (
	HookConfigCreateBefore HookEvent = iota
	HookConfigCreateAfter
	HookConfigUpdateBefore
	HookConfigUpdateAfter
	HookConfigDeleteAfter
)

var configHooks = make(map[string]map[HookEvent][]HookFunc)

func registerHook(format string, event HookEvent, hook HookFunc) {
	if _, found := configHooks[format]; !found {
		configHooks[format] = make(map[HookEvent][]HookFunc)
	}

	configHooks[format][event] = append(configHooks[format][event], hook)
}

func doHook(format string, event HookEvent, config *Config, env *Environment, admin string) {
	hooks := configHooks[FormatAll][event]
	hooks = append(hooks, configHooks[format][event]...)
	for _, hook := range hooks {
		hook(config, env, admin)
	}
}

type Config struct {
	base.Record
	Index  int         `bson:"-" json:"index"` // just for edit by
	Name   string      `bson:"name" json:"name"`
	Format string      `bson:"format" json:"format"`
	Desc   string      `bson:"desc" json:"desc"`
	Fields []table.Col `bson:"fields" json:"fields"`
}

func (c *Config) GetFieldType(name string) string {
	for _, field := range c.Fields {
		if field.Name == name {
			return field.Type
		}
	}
	return ""
}

func configCol(database string, collect string) *mgo.Collection {
	db := mongo.Session.DB(database)
	return db.C(collect)
}

func BuildConfigPageList(envId string, name, format string, params table.PageParams) (tableList *table.PageList, err error) {
	environment, err := GetEnvironmentById(envId)
	if err != nil {
		return nil, err
	}

	var list = make([]interface{}, 0, len(environment.Configs))
	for index, item := range environment.Configs {
		item.Index = index
		if strings.Contains(item.Name, name) && strings.Contains(item.Format, format) {
			list = append(list, item)
		}
	}

	return table.BuildArrayList(list, params, BuildConfigCols())
}

func BuildConfigCols() []table.Col {
	return base.AttachRecordCols([]table.Col{
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "Config name",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.SearchAble, table.SortAble),
		},
		{
			Name:    "format",
			Type:    "format",
			Desc:    "Config format",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.SearchAble, table.SortAble),
		},
		{
			Name:    "desc",
			Type:    table.TypeText,
			Desc:    "Config description",
			Feature: table.ColFeature(table.EditAble),
		},
	})
}

func ClearConfigData(database string, collect string) error {
	mgoCol := configCol(database, collect)

	return mgoCol.DropCollection()
}

func (env *Environment) CreateConfigRow(name string, format string, desc string, admin string) error {
	name = strings.TrimSpace(name)
	desc = strings.TrimSpace(desc)

	config := &Config{
		Name:   name,
		Format: format,
		Desc:   desc,
	}

	doHook(format, HookConfigCreateBefore, config, env, admin)

	env.Configs = append(env.Configs, config)

	if err := env.Save(admin); err != nil {
		return err
	}

	doHook(format, HookConfigCreateAfter, config, env, admin)

	return nil
}

func (env *Environment) CopyConfigRow(oldConfig *Config, admin string) error {
	env.Configs = append(env.Configs, oldConfig)

	if err := env.Save(admin); err != nil {
		return err
	}

	return nil
}

func (env *Environment) UpdateConfigRow(index int, desc string, admin string) error {
	var config *Config

	for i, c := range env.Configs {
		if i == index {
			config = c
			break
		}
	}

	if config == nil {
		return errors.New("config not found")
	}

	config.Desc = desc
	doHook(config.Format, HookConfigUpdateBefore, config, env, admin)

	if err := env.Save(admin); err != nil {
		return err
	}

	doHook(config.Format, HookConfigUpdateAfter, config, env, admin)

	return nil
}

func (env *Environment) DeleteConfigRow(index int, admin string) error {
	if index < 0 || index >= len(env.Configs) {
		return errors.New("index is error")
	}

	config := env.Configs[index]

	env.Configs = append(env.Configs[:index], env.Configs[index+1:]...)

	if err := env.Save(admin); err != nil {
		return err
	}

	doHook(config.Format, HookConfigDeleteAfter, config, env, admin)

	return nil
}

func FormatSetting() []string {
	return formatList
}

func init() {
	registerHook(FormatAll, HookConfigCreateBefore, func(config *Config, env *Environment, admin string) {
		config.SetCreator(admin, fmt.Sprintf("create config %v", config.Name))
	})
	registerHook(FormatAll, HookConfigUpdateBefore, func(config *Config, env *Environment, admin string) {
		config.SetEditor(admin, fmt.Sprintf("update config %v", config.Name))
	})
	registerHook(FormatAll, HookConfigDeleteAfter, func(config *Config, env *Environment, admin string) {
		err := configCol(env.Database, config.Name).DropCollection()
		if err != nil {
			logrus.Warnf("Drop delete config collection with err: %v", err)
		}
	})
}
