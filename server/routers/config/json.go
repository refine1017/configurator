package config

import (
	"encoding/json"
	"server/models"
	"server/models/table"
	"server/modules/code"
	"server/modules/context"
	"strings"
)

type JsonQueryParams struct {
	table.PageParams
	EnvId   string `binding:"Required;MaxSize(255)"`
	Collect string `binding:"Required;MaxSize(255)"`
}

func FetchJsonPageList(ctx *context.Context, params JsonQueryParams) {
	data, err := models.BuildConfigJsonPageList(params.EnvId, params.Collect, params.PageParams)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "BuildJsonPageList", err)
		return
	}

	ctx.Ack(data)
}

func ExportConfigJsonZipFile(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	buf, err := models.ExportJsonZipFile(env, collect)
	if err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "ExportJsonZipFile", err)
		return
	}

	ctx.SetFileHeader(collect + models.ExtensionZip)

	if _, err := buf.WriteTo(ctx.Resp); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "buf.WriteTo", err)
		return
	}
}

func CreateConfigJsonRow(ctx *context.Context) {
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

	if err := models.CreateConfigJsonRow(env, collect, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "CreateConfigJsonRow", err)
		return
	}

	ctx.Ack(nil)
}

func UpdateConfigJsonRow(ctx *context.Context) {
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

	if err := models.UpdateConfigJsonRow(env, collect, id, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigJsonRow", err)
		return
	}

	ctx.Ack(nil)
}

func DeleteConfigJsonRow(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	if err := models.DeleteConfigJsonRow(env, collect, id, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "UpdateConfigJsonRow", err)
		return
	}

	ctx.Ack(nil)
}

func UploadConfigJson(ctx *context.Context) {
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

	var result struct {
		Filename string `json:"filename"`
		Data     string `json:"data"`
	}

	if err := json.Unmarshal(buf, &result); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "Unmarshal", err)
		return
	}

	fileInfo := strings.Split(result.Filename, ".")
	name := fileInfo[0]

	row, err := models.GetConfigJsonRowByName(env.Database, collect, name)
	if row == nil {
		values := make(map[string]interface{})
		values["name"] = name
		values["desc"] = ""

		buf, err := json.Marshal(values)
		if err != nil {
			ctx.Refuse(code.PARAMS_ERROR, "Marshal values", err)
			return
		}

		if err := models.CreateConfigJsonRow(env, collect, string(buf), ctx.GetAdminName()); err != nil {
			ctx.Refuse(code.PARAMS_ERROR, "CreateConfigJsonRow", err)
			return
		}
	}

	row, err = models.GetConfigJsonRowByName(env.Database, collect, name)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigJsonRowByName", err)
		return
	}

	if err := models.SetConfigJsonData(env, collect, row.Id.Hex(), result.Data, ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "SetConfigJsonData", err)
		return
	}

	ctx.Ack(nil)
}

func ExportConfigJsonData(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	row, err := models.GetConfigJsonRow(env, collect, id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigJsonRow", err)
		return
	}

	ctx.SetFileHeader(row.Row["name"].(string) + models.ExtensionJson)

	if _, err := ctx.Resp.Write([]byte(row.Data)); err != nil {
		ctx.Refuse(code.SYSTEM_ERROR, "Write", err)
		return
	}
}

func GetConfigJsonData(ctx *context.Context) {
	envId := ctx.Params(":envId")
	collect := ctx.Params(":collect")
	id := ctx.Params(":id")

	env, err := models.GetEnvironmentById(envId)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetEnvironmentById", err)
		return
	}

	row, err := models.GetConfigJsonRow(env, collect, id)
	if err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "GetConfigJsonRow", err)
		return
	}

	ctx.Ack(row.Data)
}

func SetConfigJsonData(ctx *context.Context) {
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

	if err := models.SetConfigJsonData(env, collect, id, string(buf), ctx.GetAdminName()); err != nil {
		ctx.Refuse(code.PARAMS_ERROR, "SetConfigJsonData", err)
		return
	}

	ctx.Ack(nil)
}
