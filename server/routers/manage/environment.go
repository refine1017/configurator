package manage

import (
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type EnvironmentQueryParams struct {
	table.PageParams
	ProjectId string `binding:"Required;MaxSize(255)"`
	Name      string `binding:"Required;MaxSize(255)"`
}

type EnvironmentAddInfo struct {
	Name string `binding:"Required;MaxSize(255)"`
	Desc string `binding:"Required;MaxSize(255)"`
}

type EnvironmentEditInfo struct {
	Desc string `binding:"Required;MaxSize(255)"`
}

type EnvironmentCopyInfo struct {
	Name string `binding:"Required;MaxSize(255)"`
	Desc string `binding:"Required;MaxSize(255)"`
}

func FetchEnvironmentPageList(ctx *context.Context, params EnvironmentQueryParams) {
	data, err := models.BuildEnvironmentPageList(params.ProjectId, params.Name, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildEnvironmentPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateEnvironmentRow(ctx *context.Context, info EnvironmentAddInfo) {
	projectId := ctx.Params("projectId")

	project, err := models.GetProjectById(projectId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetProjectById", err)
		return
	}

	if err := models.CreateEnvironmentRow(project, info.Name, info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateEnvironmentRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateEnvironmentRow(ctx *context.Context, info EnvironmentEditInfo) {
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if err := env.UpdateInfo(info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateInfo", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteEnvironmentRow(ctx *context.Context) {
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if err := models.DeleteEnvironmentRow(env); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateEnvironmentRow", err)
		return
	}

	ctx.Ack(nil)
}

func CopyEnvironmentRow(ctx *context.Context, info EnvironmentCopyInfo) {
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	project, err := models.GetProjectById(env.ProjectId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetProjectById", err)
		return
	}

	if err := models.CopyEnvironmentRow(project, env, info.Name, info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CopyEnvironmentRow", err)
		return
	}

	ctx.Ack(nil)
}
