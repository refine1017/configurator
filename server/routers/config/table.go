package config

import (
	"encoding/json"
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type TableQueryParams struct {
	table.PageParams
	EnvId   string `binding:"Required;MaxSize(255)"`
	Collect string `binding:"Required;MaxSize(255)"`
}

func FetchTablePageList(ctx *context.Context, params TableQueryParams) {
	data, err := models.BuildConfigTablePageList(params.EnvId, params.Collect, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildTablePageList", err)
		return
	}

	ctx.Ack(data)
}

func ExportConfigTableJsonFile(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	cfg := env.GetConfigByName(collect)
	if cfg == nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigByName is empty", nil)
		return
	}

	buf, err := models.ExportTableJsonFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportKVJsonFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionJson)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func ExportConfigTableLuaFile(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	cfg := env.GetConfigByName(collect)
	if cfg == nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigByName is empty", nil)
		return
	}

	buf, err := models.ExportTableLuaFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportKVJsonFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionLua)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func ExportConfigTableExcelFile(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	cfg := env.GetConfigByName(collect)
	if cfg == nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigByName is empty", nil)
		return
	}

	buf, err := models.ExportTableExcelFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportTableExcelFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionExcel)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func CreateConfigTableRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "ctx.Req.Body().Bytes", err)
		return
	}

	if err := models.CreateConfigTableRow(env, collect, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateConfigTableRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateConfigTableRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "ctx.Req.Body().Bytes", err)
		return
	}

	if err := models.UpdateConfigTableRow(env, collect, id, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigTableRow", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteConfigTableRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := models.DeleteConfigTableRow(env, collect, id, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigTableRow", err)
		return
	}

	ctx.Ack(nil)
}

func UploadConfigTable(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := ctx.Req.Body().Bytes()
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "ctx.Req.Body().Bytes", err)
		return
	}

	var dataList []interface{}

	if err := json.Unmarshal(buf, &dataList); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "Unmarshal", err)
		return
	}

	if err := models.ClearConfigData(env.Database, collect); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "ClearConfigData", err)
		return
	}

	for _, data := range dataList {
		item, _ := json.Marshal(data)
		if err := models.CreateConfigTableRow(env, collect, string(item), ctx.GetAdminName()); err != nil {
			ctx.Refuse(code.PARAMS_ERROR, "CreateConfigTableRow", err)
			return
		}
	}

	ctx.Ack(nil)
}
