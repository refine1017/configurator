package system

import (
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type AdminQueryParams struct {
	table.PageParams
	Username string `binding:"Required;MaxSize(255)"`
	Name     string `binding:"Required;MaxSize(255)"`
}

type AdminAddInfo struct {
	Username     string `binding:"Required;MaxSize(255)"`
	Password     string `binding:"Required;MaxSize(255)"`
	Name         string `binding:"Required;MaxSize(255)"`
	Avatar       string `binding:"Required;MaxSize(255)"`
	Introduction string `binding:"Required;MaxSize(255)"`
}

type AdminEditInfo struct {
	Name         string `binding:"Required;MaxSize(255)"`
	Avatar       string `binding:"Required;MaxSize(255)"`
	Introduction string `binding:"Required;MaxSize(255)"`
}

type AdminLogQueryParams struct {
	table.PageParams
	Admin string `binding:"Required;MaxSize(255)"`
	Type  string `binding:"Required;MaxSize(255)"`
}

func FetchAdminPageList(ctx *context.Context, params AdminQueryParams) {
	data, err := models.BuildAdminPageList(params.Username, params.Name, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildAdminPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateAdminRow(ctx *context.Context, info AdminAddInfo) {
	if err := models.CreateAdminRow(info.Username, info.Password, info.Name, info.Avatar, info.Introduction, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateAdminRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateAdminRow(ctx *context.Context, info AdminEditInfo) {
	id := ctx.Params(":id")

	pro, err := models.GetAdminById(id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if err := pro.UpdateInfo(info.Name, info.Avatar, info.Introduction, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateInfo", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteAdminRow(ctx *context.Context) {
	id := ctx.Params(":id")

	if err := models.DeleteAdminRow(id); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateAdminRow", err)
		return
	}

	ctx.Ack(nil)
}

func FetchAdminLogPageList(ctx *context.Context, params AdminLogQueryParams) {
	data, err := models.BuildAdminLogPageList(params.Admin, params.Type, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildAdminLogPageList", err)
		return
	}

	ctx.Ack(data)
}
