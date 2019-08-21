package manage

import (
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type ProjectQueryParams struct {
	table.PageParams
	Name string `binding:"Required;MaxSize(255)"`
}

type ProjectAddInfo struct {
	Name string `binding:"Required;MaxSize(255)"`
	Desc string `binding:"Required;MaxSize(255)"`
}

type ProjectEditInfo struct {
	Desc string `binding:"Required;MaxSize(255)"`
}

func FetchProjectPageList(ctx *context.Context, params ProjectQueryParams) {
	data, err := models.BuildProjectPageList(params.Name, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildProjectPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateProjectRow(ctx *context.Context, info ProjectAddInfo) {
	if err := models.CreateProjectRow(info.Name, info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateProjectRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateProjectRow(ctx *context.Context, info ProjectEditInfo) {
	id := ctx.Params(":id")

	pro, err := models.GetProjectById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if err := pro.UpdateInfo(info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateInfo", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteProjectRow(ctx *context.Context) {
	id := ctx.Params(":id")

	if err := models.DeleteProjectRow(id); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateProjectRow", err)
		return
	}

	ctx.Ack(nil)
}
