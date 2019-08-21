package manage

import (
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type ServerQueryParams struct {
	table.PageParams
	ProjectId string `binding:"Required;MaxSize(255)"`
	Name      string `binding:"Required;MaxSize(255)"`
}

type ServerAddInfo struct {
	ProjectId string `binding:"Required;MaxSize(255)"`
	Name      string `binding:"Required;MaxSize(255)"`
	Url       string `binding:"Required;MaxSize(255)"`
}

type ServerEditInfo struct {
	ProjectId string `binding:"Required;MaxSize(255)"`
	Name      string `binding:"Required;MaxSize(255)"`
	Url       string `binding:"Required;MaxSize(255)"`
}

func FetchServerPageList(ctx *context.Context, params ServerQueryParams) {
	data, err := models.BuildServerPageList(params.ProjectId, params.Name, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildServerPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateServerRow(ctx *context.Context, info ServerAddInfo) {
	projectId := ctx.Params("projectId")

	if err := models.CreateServerRow(info.Name, projectId, info.Url, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateServerRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateServerRow(ctx *context.Context, info ServerEditInfo) {
	id := ctx.Params(":id")

	pro, err := models.GetServerById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetServerById", err)
	}

	if err := pro.UpdateInfo(info.Name, info.Url, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateInfo", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteServerRow(ctx *context.Context) {
	id := ctx.Params(":id")

	if err := models.DeleteServerRow(id); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateServerRow", err)
		return
	}

	ctx.Ack(nil)
}
