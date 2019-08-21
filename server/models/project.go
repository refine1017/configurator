package models

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/mongo"
	"strings"
)

type Project struct {
	base.Record
	Id   bson.ObjectId `bson:"_id" json:"id"`
	Name string        `bson:"name" json:"name"`
	Desc string        `bson:"desc" json:"desc"`
}

func projectCol() *mgo.Collection {
	return mongo.DB().C("projects")
}

func GetProjectById(id string) (*Project, error) {
	if id == "" {
		return nil, errors.New("project id is empty")
	}

	var project = &Project{}

	err := projectCol().Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(project)

	return project, err
}

func GetProjectList() ([]*Project, error) {
	var tableList []*Project

	err := projectCol().Find(nil).All(&tableList)

	return tableList, err
}

func BuildProjectPageList(name string, params table.PageParams) (tableList *table.PageList, err error) {
	condition := bson.M{"name": bson.M{
		"$regex": bson.RegEx{name, "i"},
	}}

	projects := make([]*Project, 0)

	return table.BuildColList(projectCol(), condition, params, &projects, BuildProjectCols())
}

func BuildProjectCols() []table.Col {
	return base.AttachRecordCols([]table.Col{
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "Project name",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.SearchAble, table.SortAble),
		},
		{
			Name:    "desc",
			Type:    table.TypeText,
			Desc:    "Project description",
			Feature: table.ColFeature(table.EditAble),
		},
	})
}

func CreateProjectRow(name string, desc string, admin string) error {
	name = strings.TrimSpace(name)
	desc = strings.TrimSpace(desc)

	pro := &Project{
		Name: name,
		Desc: desc,
	}

	return pro.Save(admin)
}

func DeleteProjectRow(id string) error {
	return projectCol().RemoveId(bson.ObjectIdHex(id))
}

func (pro *Project) Save(admin string) error {
	if pro.Id == "" {
		pro.Id = bson.NewObjectId()
		pro.SetCreator(admin, fmt.Sprintf("create project %v", pro.Id.Hex()))
		return projectCol().Insert(pro)
	} else {
		pro.SetEditor(admin, fmt.Sprintf("update project %v", pro.Id.Hex()))
		return projectCol().UpdateId(pro.Id, pro)
	}
}

func (pro *Project) UpdateInfo(desc string, admin string) error {
	pro.Desc = desc
	return pro.Save(admin)
}
