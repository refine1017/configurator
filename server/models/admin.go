package models

import (
	"crypto/sha256"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/pbkdf2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"server/models/base"
	"server/models/table"
	"server/modules/mongo"
	"server/modules/util"
	"time"
)

type Admin struct {
	base.Record
	Id            bson.ObjectId `bson:"_id" json:"id"`
	Username      string        `bson:"username" json:"username"`
	Password      string        `bson:"password" json:"-"`
	Name          string        `bson:"name" json:"name"`
	Avatar        string        `bson:"avatar" json:"avatar"`
	Introduction  string        `bson:"introduction" json:"introduction"`
	LastLoginTime string        `bson:"last_login_time" json:"last_login_time"`
	LastLoginIp   string        `bson:"last_login_ip" json:"last_login_ip"`
	Salt          string        `bson:"salt" json:"-"`
	Roles         []string      `bson:"-" json:"roles"`
}

func hashPassword(passwd, salt string) string {
	tempPasswd := pbkdf2.Key([]byte(passwd), []byte(salt), 10000, 50, sha256.New)
	return fmt.Sprintf("%x", tempPasswd)
}

func adminCol() *mgo.Collection {
	return mongo.DB().C("admins")
}

func GetAdminById(id string) (*Admin, error) {
	var admin = &Admin{}

	err := adminCol().FindId(bson.ObjectIdHex(id)).One(admin)

	return admin, err
}

func GetAdminByUsername(username string) (*Admin, error) {
	var admin = &Admin{}

	err := adminCol().Find(bson.M{
		"username": username,
	}).One(admin)

	return admin, err
}

func ValidAdminPassword(admin *Admin, password string) bool {
	return hashPassword(password, admin.Salt) == admin.Password
}

func BuildAdminPageList(username, name string, params table.PageParams) (tableList *table.PageList, err error) {
	condition := bson.M{
		"username": bson.M{"$regex": bson.RegEx{username, "i"}},
		"name":     bson.M{"$regex": bson.RegEx{name, "i"}},
	}

	projects := make([]*Admin, 0)

	return table.BuildColList(adminCol(), condition, params, &projects, BuildAdminCols())
}

func BuildAdminCols() []table.Col {
	return []table.Col{
		{
			Name:    "username",
			Type:    table.TypeString,
			Desc:    "admin username",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.SearchAble, table.SortAble),
		},
		{
			Name:    "password",
			Type:    table.TypeString,
			Desc:    "admin password",
			Feature: table.ColFeature(table.EditAble, table.DisableChange, table.DisableShow),
		},
		{
			Name:    "name",
			Type:    table.TypeString,
			Desc:    "admin name",
			Feature: table.ColFeature(table.EditAble, table.SearchAble, table.SortAble),
		},
		{
			Name:    "avatar",
			Type:    table.TypeString,
			Desc:    "admin avatar",
			Feature: table.ColFeature(table.EditAble),
		},
		{
			Name:    "introduction",
			Type:    table.TypeText,
			Desc:    "admin introduction",
			Feature: table.ColFeature(table.EditAble),
		},
		{
			Name:    "last_login_time",
			Type:    table.TypeText,
			Desc:    "LastLoginTime",
			Feature: table.ColFeature(table.DisableChange),
		},
		{
			Name:    "last_login_ip",
			Type:    table.TypeText,
			Desc:    "LastLoginIP",
			Feature: table.ColFeature(table.DisableChange),
		},
	}
}

func CreateAdminRow(username, password, name, avatar, introduction string, admin string) error {
	salt := util.HashGet(10)

	pro := &Admin{
		Username:     username,
		Password:     hashPassword(password, salt),
		Salt:         salt,
		Name:         name,
		Avatar:       avatar,
		Introduction: introduction,
	}

	return pro.Save(admin)
}

func DeleteAdminRow(id string) error {
	return adminCol().RemoveId(bson.ObjectIdHex(id))
}

func RecordAdminLoginInfo(admin *Admin, lastLoginIp string) {
	admin.LastLoginTime = time.Now().Format(time.RFC3339)
	admin.LastLoginIp = lastLoginIp
	if err := admin.Save(admin.Name); err != nil {
		logrus.Warnf("RecordAdminLoginInfo with err: %v", err)
	}
}

func (pro *Admin) Save(admin string) error {
	if pro.Id == "" {
		pro.Id = bson.NewObjectId()
		pro.SetCreator(admin, fmt.Sprintf("create admin %v", pro.Name))
		return adminCol().Insert(pro)
	} else {
		pro.SetEditor(admin, fmt.Sprintf("update admin %v", pro.Name))
		return adminCol().UpdateId(pro.Id, pro)
	}
}

func (pro *Admin) UpdateInfo(name, avatar, introduction string, admin string) error {
	pro.Name = name
	pro.Avatar = avatar
	pro.Introduction = introduction
	return pro.Save(admin)
}
