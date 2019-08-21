package models

import (
	"errors"
	"server/models/table"
	"strings"
)

func BuildFieldPageList(envId string, configIdx int, name string, params table.PageParams) (tableList *table.PageList, err error) {
	env, err := GetEnvironmentById(envId)
	if err != nil {
		return nil, err
	}

	if configIdx < 0 || configIdx >= len(env.Configs) {
		return nil, errors.New("config index is error")
	}

	config := env.Configs[configIdx]

	var list = make([]interface{}, 0, len(config.Fields))
	for idx, item := range config.Fields {
		item.Idx = idx
		if strings.Contains(item.Name, name) {
			list = append(list, item)
		}
	}

	return table.BuildArrayList(list, params, table.FieldCols())
}

func (env *Environment) CreateFieldRow(config *Config, name string, _type string, desc string, index string, admin string) error {
	name = strings.TrimSpace(name)
	desc = strings.TrimSpace(desc)

	field := table.Col{
		Name:  name,
		Type:  _type,
		Desc:  desc,
		Index: index,
	}

	config.Fields = append(config.Fields, field)

	if err := env.Save(admin); err != nil {
		return err
	}

	return nil
}

func (env *Environment) UpdateFieldRow(config *Config, idx int, _type string, desc string, index string, admin string) error {
	var field *table.Col

	for i, f := range config.Fields {
		if i == idx {
			field = &f
			break
		}
	}

	if field == nil {
		return errors.New("field not found")
	}

	field.Type = _type
	field.Desc = desc
	field.Index = index

	config.Fields[idx] = *field

	if err := env.Save(admin); err != nil {
		return err
	}

	return nil
}

func (env *Environment) DeleteFieldRow(config *Config, index int, admin string) error {
	if index < 0 || index >= len(config.Fields) {
		return errors.New("index is error")
	}

	config.Fields = append(config.Fields[:index], config.Fields[index+1:]...)

	if err := env.Save(admin); err != nil {
		return err
	}

	return nil
}
