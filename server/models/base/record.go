package base

import (
	"gopkg.in/mgo.v2/bson"
	"server/models/table"
	"time"
)

type Record struct {
	Creator string `bson:"creator" json:"creator"`
	Created string `bson:"created" json:"created"`
	Editor  string `bson:"editor" json:"editor"`
	Updated string `bson:"updated" json:"updated"`
	Log     string `bson:"log" json:"log"`
}

var recordTableCols = []table.Col{
	{
		Name:    "created",
		Type:    "string",
		Desc:    "Create time",
		Feature: table.ColFeature(table.DisableExport),
	},
	{
		Name:    "updated",
		Type:    "string",
		Desc:    "Last update time",
		Feature: table.ColFeature(table.DisableExport),
	},
	{
		Name:    "creator",
		Type:    "string",
		Desc:    "Creator's name",
		Feature: table.ColFeature(table.DisableExport),
	},
	{
		Name:    "editor",
		Type:    "string",
		Desc:    "Last editor's name",
		Feature: table.ColFeature(table.DisableExport),
	},
}

func AttachRecordCols(cols []table.Col) []table.Col {
	return append(cols, recordTableCols...)
}

func AttachCreator(m bson.M, admin string) bson.M {
	m["created"] = time.Now().Format(time.RFC3339)
	m["creator"] = admin
	return AttachEditor(m, admin)
}

func AttachEditor(m bson.M, admin string) bson.M {
	m["updated"] = time.Now().Format(time.RFC3339)
	m["editor"] = admin
	return m
}

func (r *Record) SetCreator(admin string, log string) {
	r.Creator = admin
	r.Created = time.Now().Format(time.RFC3339)
	r.SetEditor(admin, log)
}

func (r *Record) SetEditor(admin string, log string) {
	r.Editor = admin
	r.Log = log
	r.Updated = time.Now().Format(time.RFC3339)
}

func (r *Record) SetLog(log string) {
	r.Log = log
}
