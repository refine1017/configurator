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

var defaultJsonRow = `{"name":"JSON_TEST","desc":"json data"}`

type JsonRow struct {
	base.Record
	Id     bson.ObjectId          `bson:"_id" json:"id"`
	Values string                 `bson:"values" json:"-"`
	Data   string                 `bson:"data" json:"-"`
	Row    map[string]interface{} `bson:"-" json:"values"`
}

func (t *JsonRow) parse() {
	_ = json.Unmarshal([]byte(t.Values), &t.Row)
}

func GetConfigJsonList(database string, collect string) ([]*JsonRow, error) {
	var jsonList []*JsonRow

	err := configCol(database, collect).Find(nil).All(&jsonList)

	for _, data := range jsonList {
		data.parse()
	}

	return jsonList, err
}

func BuildConfigJsonPageList(envId string, collect string, params table.PageParams) (tableList *table.PageList, err error) {
	env, err := GetEnvironmentById(envId)
	if err != nil {
		return nil, err
	}

	condition := bson.M{}

	var list = make([]*JsonRow, 0)

	tableList, err = table.BuildColList(configCol(env.Database, collect), condition, params, &list, BuildConfigJsonCols())
	if err != nil {
		return
	}

	for _, item := range list {
		item.parse()
	}

	return
}

func BuildConfigJsonCols() []table.Col {
	return []table.Col{
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "name",
			Feature: table.ColFeature(table.EditAble, table.SearchAble),
		},
		{
			Name:    "desc",
			Type:    table.TypeText,
			Desc:    "description",
			Feature: table.ColFeature(table.EditAble),
		},
	}
}

func CreateConfigJsonRow(env *Environment, collect string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &JsonRow{}
	row.Id = bson.NewObjectId()
	row.SetCreator(admin, fmt.Sprintf("create json row %v", row.Id.Hex()))
	row.Values = values

	err := mgoCol.Insert(row)

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("create row %v", row.Id.Hex()), admin)

	return nil
}

func DeleteConfigJsonRow(env *Environment, collect string, id string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	err := mgoCol.RemoveId(bson.ObjectIdHex(id))

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("delete row %v", id), admin)

	return nil
}

func UpdateConfigJsonRow(env *Environment, collect string, id string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &JsonRow{}

	err := mgoCol.FindId(bson.ObjectIdHex(id)).One(row)
	if err != nil {
		return err
	}

	row.Values = values
	row.SetEditor(admin, fmt.Sprintf("update row %v", row.Id.Hex()))

	err = mgoCol.UpdateId(row.Id, row)
	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("update row %v", row.Id.Hex()), admin)

	return nil
}

func CopyConfigJson(fromEnv *Environment, fromCollect string, toEnv *Environment, toCollect string) error {
	fromCol := configCol(fromEnv.Database, fromCollect)

	var rows []*JsonRow

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

func GetConfigJsonRow(env *Environment, collect string, id string) (*JsonRow, error) {
	mgoCol := configCol(env.Database, collect)

	var row = &JsonRow{}

	err := mgoCol.FindId(bson.ObjectIdHex(id)).One(row)
	if err != nil {
		return nil, err
	}

	row.parse()

	return row, err
}

func GetConfigJsonRowByName(database string, collect string, name string) (*JsonRow, error) {
	mgoCol := configCol(database, collect)

	var row = &JsonRow{}

	err := mgoCol.Find(bson.M{"values": bson.M{"$regex": bson.RegEx{name, "i"}}}).One(row)
	if err != nil {
		return nil, err
	}

	row.parse()

	return row, err
}

func SetConfigJsonData(env *Environment, collect string, id string, data string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	err := mgoCol.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": base.AttachEditor(bson.M{
		"data": data,
	}, admin)})

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("update row %v json data", id), admin)

	return nil
}

func ExportJsonZipFile(env *Environment, collect string) (*bytes.Buffer, error) {
	data, err := GetConfigJsonList(env.Database, collect)
	if err != nil {
		return nil, err
	}

	builder := export.NewZipBuilder()

	for _, row := range data {
		name := row.Row["name"].(string) + ExtensionJson
		writer, err := builder.GetFileWriter(name)
		if err != nil {
			return nil, err
		}
		_, err = writer.Write([]byte(row.Data))
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportJsonProtoFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewProtoMessageBuilder(cfg.Name, ProtoPrefix, cfg.Desc)

	cols := BuildConfigJsonCols()

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc)
	}

	builder.AddField("string", "data", "json data")

	return builder.Buffer(), nil
}

func ExportJsonEntitasFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewEntitasCodeBuilder(cfg.Name, "", cfg.Desc, false)

	cols := BuildConfigJsonCols()

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc, "")
	}

	builder.AddField("string", "data", "json data", "")

	return builder.Buffer(), nil
}

func ExportJsonCSFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewCSCodeBuilder(cfg.Name, "", cfg.Desc, false)

	cols := BuildConfigJsonCols()

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc, "")
	}

	builder.AddField("string", "data", "json data", "")

	return builder.Buffer(), nil
}

func ExportJsonDir(env *Environment, cfg *Config, builder *export.ZipBuilder) error {
	data, err := GetConfigJsonList(env.Database, cfg.Name)
	if err != nil {
		return err
	}

	fileList := make(map[string]string)

	for _, row := range data {
		name := cfg.Name + "/" + row.Row["name"].(string) + ExtensionJson

		fileList[name] = row.Data
	}

	for name, content := range fileList {
		writer, err := builder.GetFileWriter(name)
		if err != nil {
			return err
		}

		_, err = writer.Write([]byte(content))
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	registerHook(FormatJson, HookConfigCreateAfter, func(config *Config, env *Environment, admin string) {
		err := CreateConfigJsonRow(env, config.Name, defaultJsonRow, admin)
		if err != nil {
			logrus.Error("json create default row with err: %v", err)
		}
	})
}
