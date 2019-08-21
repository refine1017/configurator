package config

import (
	"encoding/json"
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
)

type KVQueryParams struct {
	table.PageParams
	EnvId   string `binding:"Required;MaxSize(255)"`
	Collect string `binding:"Required;MaxSize(255)"`
}

func FetchKVPageList(ctx *context.Context, params KVQueryParams) {
	data, err := models.BuildConfigKVPageList(params.EnvId, params.Collect, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildKVPageList", err)
		return
	}

	ctx.Ack(data)
}

func ExportConfigKVJsonFile(ctx *context.Context) {
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

	buf, err := models.ExportKVJsonFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportJsonZipFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionJson)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func ExportConfigKVLuaFile(ctx *context.Context) {
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

	buf, err := models.ExportKVLuaFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportJsonZipFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionLua)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func ExportConfigKVExcelFile(ctx *context.Context) {
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

	buf, err := models.ExportKVExcelFile(env, cfg)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportKVJsonFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionExcel)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func CreateConfigKVRow(ctx *context.Context) {
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

	if err := models.CreateConfigKVRow(env, collect, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateConfigKVRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateConfigKVRow(ctx *context.Context) {
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

	if err := models.UpdateConfigKVRow(env, collect, id, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigKVRow", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteConfigKVRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := models.DeleteConfigKVRow(env, collect, id, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigKVRow", err)
		return
	}

	ctx.Ack(nil)
}

func UploadConfigKV(ctx *context.Context) {
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
		ctx.Refuse(code.PARAMS_ERROR, "ClearConfigKVList", err)
		return
	}

	for _, data := range dataList {
		item, _ := json.Marshal(data)
		if err := models.CreateConfigKVRow(env, collect, string(item), ctx.GetAdminName()); err != nil {
			ctx.Refuse(code.PARAMS_ERROR, "CreateConfigKVRow", err)
			return
		}
	}

	ctx.Ack(nil)
}
