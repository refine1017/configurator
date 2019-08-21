package manage

import (
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type ConfigQueryParams struct {
	table.PageParams
	EnvId  string `binding:"Required;MaxSize(255)"`
	Name   string `binding:"Required;MaxSize(255)"`
	Format string `binding:"Required;MaxSize(255)"`
}

type ConfigAddInfo struct {
	Name   string `binding:"Required;MaxSize(255)"`
	Format string `binding:"Required;MaxSize(255)"`
	Desc   string `binding:"Required;MaxSize(255)"`
}

type ConfigEditInfo struct {
	Desc string `binding:"Required;MaxSize(255)"`
}

func FetchConfigPageList(ctx *context.Context, params ConfigQueryParams) {
	data, err := models.BuildConfigPageList(params.EnvId, params.Name, params.Format, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildConfigPageList", err)
		return
	}

	ctx.Ack(data)
}

func CreateConfigRow(ctx *context.Context, info ConfigAddInfo) {
	envId := ctx.Params(":envId")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := env.CreateConfigRow(info.Name, info.Format, info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateConfigRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateConfigRow(ctx *context.Context, info ConfigEditInfo) {
	envId := ctx.Params(":envId")
	index := ctx.ParamsInt(":index")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := env.UpdateConfigRow(index, info.Desc, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigRow", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteConfigRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	index := ctx.ParamsInt(":index")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := env.DeleteConfigRow(index, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigRow", err)
		return
	}

	ctx.Ack(nil)
}

func PushConfig(ctx *context.Context) {
	envId := ctx.Params(":envId")
	config := ctx.Params(":config")
	server := ctx.Params(":server")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	data := models.PushConfig(env, config, server)

	ctx.Ack(data)
}

func MergeConfigInfo(ctx *context.Context) {
	envId := ctx.Params(":envId")
	config := ctx.Params(":config")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "Params Error", err)
		return
	}

	activities, mergeError, err := models.MergeConfigInfo(env, config)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "System Error", err)
		return
	}

	var result struct {
		Activities []*models.ConfigActivity `json:"activities"`
		MergeError string                   `json:"merge_error"`
	}

	result.Activities = activities
	result.MergeError = mergeError

	ctx.Ack(result)
}

func MergeConfig(ctx *context.Context) {
	envId := ctx.Params(":envId")
	config := ctx.Params(":config")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	err = models.MergeConfig(env, config, ctx.GetAdminName())
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "MergeConfig", err)
		return
	}

	ctx.Ack(nil)
}

func ExportConfigJson(ctx *context.Context) {
	envId := ctx.Params(":envId")
	config := ctx.Params(":config")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := models.ExportEnvJsonZip(env, config)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEnvJsonZip", err)
		return
	}

	ctx.SetFileHeader(env.Name + "_json" + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}

func ExportConfigLua(ctx *context.Context) {
	envId := ctx.Params(":envId")
	config := ctx.Params(":config")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := models.ExportEnvLuaZip(env, config)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEnvJsonZip", err)
		return
	}

	ctx.SetFileHeader(env.Name + "_lua" + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}

func ExportConfigExcel(ctx *context.Context) {
	envId := ctx.Params(":envId")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := models.ExportEnvExcelZip(env)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEnvExcelZip", err)
		return
	}

	ctx.SetFileHeader(env.Name + "_excel" + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}

func ExportConfigProto(ctx *context.Context) {
	envId := ctx.Params(":envId")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	buf, err := models.ExportEnvProto(env)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEnvProto", err)
		return
	}

	ctx.SetFileHeader(env.Name + models.ExtensionProto)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}

func ExportConfigEntitas(ctx *context.Context) {
	envId := ctx.Params(":envId")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	buf, err := models.ExportEntitasProto(env)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEntitasProto", err)
		return
	}

	ctx.SetFileHeader(env.Name + "_cs" + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}

func ExportConfigCS(ctx *context.Context) {
	envId := ctx.Params(":envId")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
	}

	buf, err := models.ExportCSProto(env)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportEntitasProto", err)
		return
	}

	ctx.SetFileHeader(env.Name + "_cs" + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "WriteTo", err)
		return
	}
}
