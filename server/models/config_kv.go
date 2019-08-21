package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/export"
)

const (
	KVFieldKey   = "key"
	KVFieldType  = "type"
	KVFieldValue = "value"
	KVFieldDesc  = "desc"
)

type KVRow struct {
	base.Record
	Id     bson.ObjectId          `bson:"_id" json:"id"`
	Values string                 `bson:"values" json:"-"`
	Row    map[string]interface{} `bson:"-" json:"values"`
}

func (t *KVRow) parse() {
	err := json.Unmarshal([]byte(t.Values), &t.Row)
	if err != nil {
		logrus.Warnf("KVRow parse[%v] with err: %v", t.Values, err)
	}
}

var defaultKVRow = `{"key":"TEST_KEY","value":"TEST_VALUE","desc":"test data"}`

func GetConfigKVList(database string, collect string) ([]*KVRow, error) {
	var tableList []*KVRow

	err := configCol(database, collect).Find(nil).All(&tableList)

	for _, data := range tableList {
		data.parse()
	}

	return tableList, err
}

func BuildConfigKVPageList(envId string, collect string, params table.PageParams) (tableList *table.PageList, err error) {
	env, err := GetEnvironmentById(envId)
	if err != nil {
		return nil, err
	}

	condition := bson.M{}

	var list = make([]*TableRow, 0)

	tableList, err = table.BuildColList(configCol(env.Database, collect), condition, params, &list, BuildConfigKVCols())
	if err != nil {
		return
	}

	for _, item := range list {
		item.parse()
	}

	return
}

func BuildConfigKVCols() []table.Col {
	return []table.Col{
		{
			Name:    "key",
			Type:    table.TypeString,
			Desc:    "key",
			Feature: table.ColFeature(table.EditAble, table.SearchAble),
		},
		{
			Name:    "type",
			Type:    table.TypeFields,
			Desc:    "value's type",
			Feature: table.ColFeature(table.EditAble),
		},
		{
			Name:    "value",
			Type:    table.TypeString,
			Desc:    "value",
			Feature: table.ColFeature(table.EditAble),
		},
		{
			Name:    "desc",
			Type:    table.TypeText,
			Desc:    "description",
			Feature: table.ColFeature(table.EditAble),
		},
	}
}

func CreateConfigKVRow(env *Environment, collect string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &KVRow{}
	row.Id = bson.NewObjectId()
	row.SetCreator(admin, fmt.Sprintf("create kv row %v", row.Id.Hex()))
	row.Values = values

	err := mgoCol.Insert(row)

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("create row %v", row.Id.Hex()), admin)

	return nil
}

func UpdateConfigKVRow(env *Environment, collect string, id string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &KVRow{}

	err := mgoCol.FindId(bson.ObjectIdHex(id)).One(row)
	if err != nil {
		return err
	}

	row.Values = values
	row.SetEditor(admin, "update")

	err = mgoCol.UpdateId(row.Id, row)
	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("update row %v", row.Id.Hex()), admin)

	return nil
}

func DeleteConfigKVRow(env *Environment, collect string, id string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	err := mgoCol.RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("delete row %v", id), admin)

	return nil
}

func CopyConfigKV(fromEnv *Environment, fromCollect string, toEnv *Environment, toCollect string) error {
	fromCol := configCol(fromEnv.Database, fromCollect)

	var rows []*KVRow

	if err := fromCol.Find(nil).All(&rows); err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	mgoCol := configCol(toEnv.Database, toCollect)

	var newRows = make([]interface{}, 0, len(rows))
	for _, row := range rows {
		newRows = append(newRows, row)
	}

	return mgoCol.Insert(newRows...)
}

func ExportKVJsonFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	var fields = make(map[string]string)
	var values = make(map[string]interface{})

	for _, item := range data {
		_type, ok := item.Row[KVFieldType].(string)
		if !ok {
			_type = "string"
		}

		key, ok := item.Row[KVFieldKey].(string)
		if !ok {
			key = "Empty"
		}

		value, ok := item.Row[KVFieldValue].(string)
		if !ok {
			value = ""
		}

		fields[key] = _type
		values[key] = value
	}

	builder := export.NewTableJsonBuilder()
	builder.SetObjValue(fields, values)

	return builder.Buffer()
}

func ExportKVLuaFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	var fields = make(map[string]string)
	var values = make(map[string]interface{})

	for _, item := range data {
		_type, ok := item.Row[KVFieldType].(string)
		if !ok {
			_type = "string"
		}

		key, ok := item.Row[KVFieldKey].(string)
		if !ok {
			key = "Empty"
		}

		value, ok := item.Row[KVFieldValue].(string)
		if !ok {
			value = ""
		}

		fields[key] = _type
		values[key] = value
	}

	builder := export.NewTableLuaBuilder()
	builder.AppendKVValue(fields, values)

	return builder.Buffer()
}

func ExportKVExcelFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	cols := BuildConfigKVCols()

	var fields []string
	for _, col := range cols {
		fields = append(fields, col.Name)
	}

	builder := export.NewTableExcelBuilder()
	if err := builder.SetFields(fields); err != nil {
		return nil, err
	}
	for _, item := range data {
		if err := builder.AppendMapValues(item.Row); err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportKVProtoFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewProtoMessageBuilder(cfg.Name, ProtoPrefix, cfg.Desc)

	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	for _, col := range data {
		_type, ok := col.Row[KVFieldType].(string)
		if !ok {
			_type = "string"
		}
		key, ok := col.Row[KVFieldKey].(string)
		if !ok {
			key = "Empty"
		}
		desc, ok := col.Row[KVFieldDesc].(string)
		if !ok {
			desc = ""
		}

		builder.AddField(_type, key, desc)
	}

	return builder.Buffer(), nil
}

func ExportKVEntitasFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewEntitasCodeBuilder(cfg.Name, "", cfg.Desc, true)

	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	for _, col := range data {
		_type, ok := col.Row[KVFieldType].(string)
		if !ok {
			_type = "string"
		}
		key, ok := col.Row[KVFieldKey].(string)
		if !ok {
			key = "Empty"
		}
		desc, ok := col.Row[KVFieldDesc].(string)
		if !ok {
			desc = ""
		}

		builder.AddField(_type, key, desc, "")
	}

	return builder.Buffer(), nil
}

func ExportKVCSFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewCSCodeBuilder(cfg.Name, "", cfg.Desc, true)

	data, err := GetConfigKVList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	for _, col := range data {
		_type, ok := col.Row[KVFieldType].(string)
		if !ok {
			_type = "string"
		}
		key, ok := col.Row[KVFieldKey].(string)
		if !ok {
			key = "Empty"
		}
		desc, ok := col.Row[KVFieldDesc].(string)
		if !ok {
			desc = ""
		}

		builder.AddField(_type, key, desc, "")
	}

	return builder.Buffer(), nil
}

func init() {
	registerHook(FormatKV, HookConfigCreateAfter, func(config *Config, env *Environment, admin string) {
		err := CreateConfigKVRow(env, config.Name, defaultKVRow, admin)
		if err != nil {
			logrus.Error("kv create default row with err: %v", err)
		}
	})
}
