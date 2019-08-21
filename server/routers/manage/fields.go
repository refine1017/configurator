package manage

import (
	"errors"
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type FieldQueryParams struct {
	table.PageParams
	EnvId     string `binding:"Required;MaxSize(255)"`
	ConfigIdx int    `binding:"Required;MaxSize(255)"`
	Name      string `binding:"Required;MaxSize(255)"`
}

type FieldAddInfo struct {
	Name  string `binding:"Required;MaxSize(255)"`
	Type  string `binding:"Required;MaxSize(255)"`
	Desc  string `binding:"Required;MaxSize(255)"`
	Index string `binding:"Required;MaxSize(255)"`
}

type FieldEditInfo struct {
	Type  string `binding:"Required;MaxSize(255)"`
	Desc  string `binding:"Required;MaxSize(255)"`
	Index string `binding:"Required;MaxSize(255)"`
}

func FetchFieldPageList(ctx *context.Context, params FieldQueryParams) {
	data, err := models.BuildFieldPageList(params.EnvId, params.ConfigIdx, params.Name, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildFieldPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateFieldRow(ctx *context.Context, info FieldAddInfo) {
	envId := ctx.Params(":envId")
	configIndex := ctx.ParamsInt(":configIndex")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if configIndex < 0 || configIndex >= len(env.Configs) {
		ctx.Refuse(code.PARAMS_ERROR, "configIndex", errors.New("config index is error"))
		return
	}

	config := env.Configs[configIndex]

	if err := env.CreateFieldRow(config, info.Name, info.Type, info.Desc, info.Index, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateFieldRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateFieldRow(ctx *context.Context, info FieldEditInfo) {
	envId := ctx.Params(":envId")
	configIndex := ctx.ParamsInt(":configIndex")
	index := ctx.ParamsInt(":index")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if configIndex < 0 || configIndex >= len(env.Configs) {
		ctx.Refuse(code.PARAMS_ERROR, "configIndex", errors.New("config index is error"))
		return
	}

	config := env.Configs[configIndex]

	if err := env.UpdateFieldRow(config, index, info.Type, info.Desc, info.Index, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateFieldRow", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteFieldRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	configIndex := ctx.ParamsInt(":configIndex")
	index := ctx.ParamsInt(":index")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	if configIndex < 0 || configIndex >= len(env.Configs) {
		ctx.Refuse(code.PARAMS_ERROR, "configIndex", errors.New("config index is error"))
		return
	}

	config := env.Configs[configIndex]

	if err := env.DeleteFieldRow(config, index, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateFieldRow", err)
		return
	}

	ctx.Ack(nil)
}
