package models

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/mongo"
	"strings"
)

type Environment struct {
	base.Record
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	Desc      string        `bson:"desc" json:"desc"`
	ProjectId string        `bson:"project_id" json:"project_id"`
	Database  string        `bson:"database" json:"database"`
	Copy      string        `bson:"copy" json:"copy"`
	Configs   []*Config     `bson:"configs" json:"configs"`
}

func envCol() *mgo.Collection {
	return mongo.DB().C("environments")
}

func GetEnvironmentById(id string) (*Environment, error) {
	if id == "" {
		return nil, errors.New("environment id is empty")
	}

	var env = &Environment{}

	err := envCol().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(env)

	return env, err
}

func GetEnvironmentListByProjectId(projectId string) ([]*Environment, error) {
	var tableList []*Environment

	err := envCol().Find(bson.M{"project_id": projectId}).All(&tableList)

	return tableList, err
}

func BuildEnvironmentPageList(projectId string, name string, params table.PageParams) (tableList *table.PageList, err error) {
	condition := bson.M{
		"name":       bson.M{"$regex": bson.RegEx{name, "i"}},
		"project_id": projectId,
	}

	var environments = make([]*Environment, 0)

	return table.BuildColList(envCol(), condition, params, &environments, BuildEnvironmentCols())
}

func BuildEnvironmentCols() []table.Col {
	return base.AttachRecordCols([]table.Col{
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "Environment name",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.SearchAble, table.SortAble),
		},
		{
			Name:    "desc",
			Type:    table.TypeText,
			Desc:    "Environment description",
			Feature: table.ColFeature(table.EditAble, table.SortAble),
		},
		{
			Name: "database",
			Type: "string",
			Desc: "Environment database",
		},
	})
}

func CreateEnvironmentRow(project *Project, name string, desc string, admin string) error {
	name = strings.TrimSpace(name)
	desc = strings.TrimSpace(desc)

	database := GenerateDBName(project.Name, name)

	env := &Environment{
		Name:      name,
		Desc:      desc,
		ProjectId: project.Id.Hex(),
		Database:  database,
	}

	return env.Save(admin)
}

func DeleteEnvironmentRow(env *Environment) error {
	err := envCol().RemoveId(env.Id)
	if err != nil {
		return err
	}

	if err := mongo.Session.DB(env.Database).DropDatabase(); err != nil {
		logrus.Warnf("Drop delete environment database with err: %v", err)
	}

	return nil
}

func CopyEnvironmentRow(project *Project, oldEnv *Environment, name string, desc string, admin string) (err error) {
	name = strings.TrimSpace(name)
	desc = strings.TrimSpace(desc)

	var env *Environment

	defer func() {
		if err != nil && env != nil {
			_ = DeleteEnvironmentRow(env)
		}
	}()

	database := GenerateDBName(project.Name, name)

	env = &Environment{
		Name:      name,
		Desc:      desc,
		ProjectId: project.Id.Hex(),
		Database:  database,
		Copy:      oldEnv.Id.Hex(),
	}

	if addEnv := env.Save(admin); addEnv != nil {
		return fmt.Errorf("Create Environment with err: %v", addEnv)
	}

	for _, config := range oldEnv.Configs {
		config.SetCreator(admin, "copy config")
		env.Configs = append(env.Configs, config)

		switch config.Format {
		case FormatTable:
			if err = CopyConfigTable(oldEnv, config.Name, env, config.Name); err != nil {
				return fmt.Errorf("CopyConfigTable:%v with err: %v", config.Name, err)
			}
		case FormatKV:
			if err = CopyConfigKV(oldEnv, config.Name, env, config.Name); err != nil {
				return fmt.Errorf("CopyConfigKV:%v with err: %v", config.Name, err)
			}
		case FormatJson:
			if err = CopyConfigJson(oldEnv, config.Name, env, config.Name); err != nil {
				return fmt.Errorf("CopyConfigJson:%v with err: %v", config.Name, err)
			}
		}
	}

	return env.Save(admin)
}

func GenerateDBName(projectName, envName string) string {
	projectName = strings.Replace(projectName, ".", "_", -1)
	projectName = strings.Replace(projectName, " ", "_", -1)
	envName = strings.Replace(envName, ".", "_", -1)
	envName = strings.Replace(envName, " ", "_", -1)

	return fmt.Sprintf("db_%s_%s", strings.ToLower(projectName), strings.ToLower(envName))
}

func (env *Environment) Save(admin string) error {
	if env.Id == "" {
		env.Id = bson.NewObjectId()
		env.SetCreator(admin, fmt.Sprintf("create environment %v", env.Id.Hex()))
		return envCol().Insert(env)
	} else {
		env.SetEditor(admin, fmt.Sprintf("update environment %v", env.Id.Hex()))
		return envCol().UpdateId(env.Id, env)
	}
}

func (env *Environment) UpdateInfo(desc string, admin string) error {
	env.Desc = desc

	return env.Save(admin)
}

func (env *Environment) UpdateLog(config string, log string, admin string) {
	for _, cfg := range env.Configs {
		if cfg.Name == config {
			cfg.SetEditor(admin, log)
		}
	}

	_ = env.Save(admin)
}

func (env *Environment) GetConfigByName(name string) *Config {
	for _, cfg := range env.Configs {
		if cfg.Name == name {
			return cfg
		}
	}

	return nil
}
