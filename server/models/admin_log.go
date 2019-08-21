package models

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/models/table"
	"server/modules/mongo"
	"time"
)

type AdminLog struct {
	Id    bson.ObjectId `bson:"_id" json:"id"`
	Admin string        `bson:"admin" json:"admin"`
	Type  string        `bson:"type" json:"type"`
	Info  string        `bson:"info" json:"info"`
	Time  string        `bson:"time" json:"time"`
}

func adminLogCol() *mgo.Collection {
	return mongo.DB().C("adminLogs")
}

func BuildAdminLogPageList(admin, Type string, params table.PageParams) (tableList *table.PageList, err error) {
	condition := bson.M{
		"admin": bson.M{"$regex": bson.RegEx{admin, "i"}},
		"type":  bson.M{"$regex": bson.RegEx{Type, "i"}},
	}

	projects := make([]*AdminLog, 0)

	return table.BuildColList(adminLogCol(), condition, params, &projects, BuildAdminLogCols())
}

func BuildAdminLogCols() []table.Col {
	return []table.Col{
		{
			Name:    "time",
			Type:    table.TypeString,
			Desc:    "Time",
			Feature: table.ColFeature(table.SortAble),
		},
		{
			Name:    "admin",
			Type:    table.TypeString,
			Desc:    "Admin",
			Feature: table.ColFeature(table.SearchAble, table.SortAble),
		},
		{
			Name:    "type",
			Type:    table.TypeString,
			Desc:    "Type",
			Feature: table.ColFeature(table.SearchAble, table.SortAble),
		},
		{
			Name:    "info",
			Type:    table.TypeString,
			Desc:    "Info",
			Feature: table.ColFeature(table.SortAble),
		},
	}
}

func RecordAdminLog(admin string, Type string, info string) {
	log := &AdminLog{
		Admin: admin,
		Type:  Type,
		Info:  info,
		Time:  time.Now().Format(time.RFC3339),
	}

	log.Id = bson.NewObjectId()

	if err := adminLogCol().Insert(log); err != nil {
		logrus.Warnf("RecordAdminLog with err: %v", err)
	}
}
