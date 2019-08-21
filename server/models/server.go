package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/mongo"
	"strings"
)

type Server struct {
	base.Record
	Id        bson.ObjectId `bson:"_id" json:"id"`
	Name      string        `bson:"name" json:"name"`
	ProjectId string        `bson:"project_id" json:"project_id"`
	Url       string        `bson:"url" json:"url"`
}

func serverCol() *mgo.Collection {
	return mongo.DB().C("servers")
}

func GetServerById(id string) (*Server, error) {
	var admin = &Server{}

	err := serverCol().FindId(bson.ObjectIdHex(id)).One(admin)

	return admin, err
}

func BuildServerPageList(projectId string, name string, params table.PageParams) (tableList *table.PageList, err error) {
	condition := bson.M{
		"name":       bson.M{"$regex": bson.RegEx{name, "i"}},
		"project_id": projectId,
	}

	projects := make([]*Server, 0)

	return table.BuildColList(serverCol(), condition, params, &projects, BuildServerCols())
}

func GetServersByProjectId(projectId string) ([]*Server, error) {
	condition := bson.M{
		"project_id": projectId,
	}

	var servers []*Server

	err := serverCol().Find(condition).All(&servers)

	return servers, err
}

func BuildServerCols() []table.Col {
	return base.AttachRecordCols([]table.Col{
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "Server name",
			Feature: table.ColFeature(table.EditAble, table.SearchAble, table.SortAble),
		},
		{
			Name:    "url",
			Type:    table.TypeText,
			Desc:    "Server urls",
			Feature: table.ColFeature(table.EditAble, table.SortAble),
		},
	})
}

func CreateServerRow(name, projectId, url string, admin string) error {
	name = strings.TrimSpace(name)
	url = strings.TrimSpace(url)

	pro := &Server{
		Name:      name,
		ProjectId: projectId,
		Url:       url,
	}

	return pro.Save(admin)
}

func DeleteServerRow(id string) error {
	return serverCol().RemoveId(bson.ObjectIdHex(id))
}

func (ser *Server) Save(admin string) error {
	if ser.Id == "" {
		ser.Id = bson.NewObjectId()
		ser.SetCreator(admin, fmt.Sprintf("create server %v", ser.Id.Hex()))
		return serverCol().Insert(ser)
	} else {
		ser.SetEditor(admin, fmt.Sprintf("update server %v", ser.Id.Hex()))
		return serverCol().UpdateId(ser.Id, ser)
	}
}

func (ser *Server) UpdateInfo(name, url string, admin string) error {
	ser.Name = name
	ser.Url = url
	return ser.Save(admin)
}
