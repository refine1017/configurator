package user

import (
	"errors"
	"github.com/sirupsen/logrus"
	"server/models"
	"server/modules/auth"
	"server/modules/code"
	"server/modules/context"
)

// SignInForm form for signing in with user/password
type SignInForm struct {
	UserName string `binding:"Required;MaxSize(255)"`
	Password string `binding:"Required;MaxSize(255)"`
	Remember bool
}

func Login(ctx *context.Context, form SignInForm) {
	admin, err := models.GetAdminByUsername(form.UserName)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetAdminByUsername", err)
		return
	}

	if !models.ValidAdminPassword(admin, form.Password) {
		ctx.Refuse(code.PARAMS_ERROR, "ValidAdminPassword", errors.New("ValidAdminPassword fail"))
		return
	}

	models.RecordAdminLoginInfo(admin, ctx.Req.RemoteAddr)

	auth.SetAdmin(ctx.Session, admin)

	var Data struct {
		Token string `json:"token"`
	}

	Data.Token = "admin-token"

	logrus.Infof("admin[%v-%v] login success", ctx.Session.Get("adminName"), ctx.Session.Get("adminId"))

	ctx.Ack(Data)
}

func Logout(ctx *context.Context) {
	auth.DelAdmin(ctx.Session)

	ctx.Ack("success")
}

func Info(ctx *context.Context) {
	admin := ctx.GetAdmin()
	admin.Roles = []string{"admin"}

	ctx.Ack(admin)
}

func Projects(ctx *context.Context) {
	projects, err := getUserProjects()
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "getUserProjects", err)
		return
	}

	ctx.Ack(projects)
}

type ProjectInfo struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Envs    []*EnvInfo    `json:"envs"`
	Servers []*ServerInfo `json:"servers"`
}

type ServerInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type EnvInfo struct {
	Id      string        `json:"id"`
	Name    string        `json:"name"`
	Configs []*ConfigInfo `json:"configs"`
}

type ConfigInfo struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

func getUserProjects() ([]*ProjectInfo, error) {
	var userProjects []*ProjectInfo

	projects, err := models.GetProjectList()
	if err != nil {
		return userProjects, err
	}

	for _, project := range projects {
		userProject := &ProjectInfo{
			Id:   project.Id.Hex(),
			Name: project.Name,
		}

		servers, err := models.GetServersByProjectId(userProject.Id)
		if err != nil {
			return userProjects, err
		}

		for _, server := range servers {
			userProject.Servers = append(userProject.Servers, &ServerInfo{
				Id:   server.Id.Hex(),
				Name: server.Name,
			})
		}

		envs, err := models.GetEnvironmentListByProjectId(userProject.Id)
		if err != nil {
			return userProjects, err
		}

		for _, env := range envs {
			userEnv := &EnvInfo{
				Id:      env.Id.Hex(),
				Name:    env.Name,
				Configs: make([]*ConfigInfo, 0),
			}

			for _, conf := range env.Configs {
				userEnv.Configs = append(userEnv.Configs, &ConfigInfo{
					Name:   conf.Name,
					Format: conf.Format,
				})
			}

			userProject.Envs = append(userProject.Envs, userEnv)
		}

		userProjects = append(userProjects, userProject)
	}

	return userProjects, nil
}
