package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/export"
	"sort"
	"strconv"
)

type TableRow struct {
	base.Record
	Id     bson.ObjectId          `bson:"_id" json:"id"`
	Values string                 `bson:"values" json:"-"`
	Row    map[string]interface{} `bson:"-" json:"values"`
}

func (t *TableRow) parse() {
	err := json.Unmarshal([]byte(t.Values), &t.Row)
	if err != nil {
		logrus.Warnf("TableRow parse[%v] with err: %v", t.Values, err)
	}
}

var defaultTableFields = []table.Col{
	{Name: "id", Type: table.TypeUint32, Desc: "ID", Index: table.IndexPrimary},
}

var defaultTableRow = `{"id":1}`

func GetConfigTableList(database string, collect string) ([]*TableRow, error) {
	var tableList []*TableRow

	err := configCol(database, collect).Find(nil).All(&tableList)

	for _, data := range tableList {
		data.parse()
	}

	return tableList, err
}

func BuildConfigTablePageList(envId string, collect string, params table.PageParams) (tableList *table.PageList, err error) {
	env, err := GetEnvironmentById(envId)
	if err != nil {
		return nil, err
	}

	config := env.GetConfigByName(collect)
	if config == nil {
		return nil, errors.New("config not exists")
	}

	var list = make([]*TableRow, 0)

	limit := params.Limit
	page := params.Page
	sortField := params.Sort

	params.Limit = 0
	params.Page = 0
	params.Sort = ""

	tableList, err = table.BuildColList(configCol(env.Database, collect), nil, params, &list, BuildConfigTableCols(env, collect))
	if err != nil {
		return
	}

	for _, item := range list {
		item.parse()
	}

	if sortField != "" {
		var desc = false
		if sortField[0] == '-' {
			desc = true
			sortField = sortField[1:]
		}

		fieldType := config.GetFieldType(sortField)
		if fieldType != "" {
			sort.Sort(&TableRowSorter{
				Rows:      list,
				Field:     sortField,
				FieldType: fieldType,
				Desc:      desc,
			})
		}
	}

	if limit > 0 {
		start, end := table.CalculateLimitStartAndEnd(page, limit, len(list))
		tableList.Items = list[start:end]
	}

	return
}

func BuildConfigTableCols(env *Environment, collect string) []table.Col {
	for _, conf := range env.Configs {
		if conf.Name == collect {
			for i, field := range conf.Fields {
				field.Feature = table.FieldTypeFeature[field.Type]

				if field.Index != "" {
					field.Feature = table.AddFeature(field.Feature, table.SearchAble)
				}

				conf.Fields[i] = field
			}
			return conf.Fields
		}
	}

	return nil
}

func CreateConfigTableRow(env *Environment, collect string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &TableRow{}
	row.Id = bson.NewObjectId()
	row.SetCreator(admin, fmt.Sprintf("create table row %v", row.Id.Hex()))
	row.Values = values

	err := mgoCol.Insert(row)

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("create row %v", row.Id.Hex()), admin)

	return nil
}

func UpdateConfigTableRow(env *Environment, collect string, id string, values string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	row := &TableRow{}

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

func DeleteConfigTableRow(env *Environment, collect string, id string, admin string) error {
	mgoCol := configCol(env.Database, collect)

	err := mgoCol.RemoveId(bson.ObjectIdHex(id))

	if err != nil {
		return err
	}

	env.UpdateLog(collect, fmt.Sprintf("delete row %v", id), admin)

	return nil
}

func CopyConfigTable(fromEnv *Environment, fromCollect string, toEnv *Environment, toCollect string) error {
	fromCol := configCol(fromEnv.Database, fromCollect)

	var rows []*TableRow

	if err := fromCol.Find(nil).All(&rows); err != nil {
		return fmt.Errorf("find %v from %v with err: %v", fromCollect, fromEnv.Database, err)
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

func ExportTableJsonFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigTableList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	cols := BuildConfigTableCols(env, cfg.Name)

	fields := make(map[string]string)

	for _, col := range cols {
		fields[col.Name] = col.Type
	}

	builder := export.NewTableJsonBuilder()
	for _, item := range data {
		builder.AppendArrValue(fields, item.Row)
	}

	return builder.Buffer()
}

func ExportTableLuaFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigTableList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	cols := BuildConfigTableCols(env, cfg.Name)

	fields := make(map[string]string)

	for _, col := range cols {
		fields[col.Name] = col.Type
	}

	builder := export.NewTableLuaBuilder()
	for _, item := range data {
		builder.AppendArrValue(fields, item.Row)
	}

	return builder.Buffer()
}

func ExportTableExcelFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	data, err := GetConfigTableList(env.Database, cfg.Name)
	if err != nil {
		return nil, err
	}

	cols := BuildConfigTableCols(env, cfg.Name)

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

func ExportTableProtoFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewProtoMessageBuilder(cfg.Name, ProtoPrefix, cfg.Desc)

	cols := BuildConfigTableCols(env, cfg.Name)

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc)
	}

	return builder.Buffer(), nil
}

func ExportTableEntitasFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewEntitasCodeBuilder(cfg.Name, "", cfg.Desc, false)

	cols := BuildConfigTableCols(env, cfg.Name)

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc, col.Index)
	}

	return builder.Buffer(), nil
}

func ExportTableCSFile(env *Environment, cfg *Config) (*bytes.Buffer, error) {
	builder := export.NewCSCodeBuilder(cfg.Name, "", cfg.Desc, false)

	cols := BuildConfigTableCols(env, cfg.Name)

	for _, col := range cols {
		builder.AddField(col.Type, col.Name, col.Desc, col.Index)
	}

	return builder.Buffer(), nil
}

type TableRowSorter struct {
	Rows      []*TableRow
	Field     string
	FieldType string
	Desc      bool
}

func (s *TableRowSorter) Len() int      { return len(s.Rows) }
func (s *TableRowSorter) Swap(i, j int) { s.Rows[i], s.Rows[j] = s.Rows[j], s.Rows[i] }
func (s *TableRowSorter) Less(i, j int) bool {
	v1 := s.Rows[i].Row[s.Field]
	v2 := s.Rows[j].Row[s.Field]

	switch s.FieldType {
	case table.TypeInt32, table.TypeInt64, table.TypeUint32, table.TypeUint64:
		value1, _ := strconv.Atoi(fmt.Sprintf("%v", v1))
		value2, _ := strconv.Atoi(fmt.Sprintf("%v", v2))

		if s.Desc {
			return value1 > value2
		} else {
			return value1 < value2
		}

	case table.TypeFloat, table.TypeDouble:
		value1, _ := strconv.ParseFloat(fmt.Sprintf("%v", v1), 64)
		value2, _ := strconv.ParseFloat(fmt.Sprintf("%v", v2), 64)

		if s.Desc {
			return value1 > value2
		} else {
			return value1 < value2
		}

	default:
		value1 := fmt.Sprintf("%v", v1)
		value2 := fmt.Sprintf("%v", v2)

		if s.Desc {
			return value1 > value2
		} else {
			return value1 < value2
		}
	}
}

func init() {
	registerHook(FormatTable, HookConfigCreateBefore, func(config *Config, env *Environment, admin string) {
		config.Fields = defaultTableFields
	})
	registerHook(FormatTable, HookConfigCreateAfter, func(config *Config, env *Environment, admin string) {
		err := CreateConfigTableRow(env, config.Name, defaultTableRow, admin)
		if err != nil {
			logrus.Error("table create default row with err: %v", err)
		}
	})
}
