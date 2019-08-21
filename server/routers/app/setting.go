package app

import (
	"server/models"
	"server/models/table"
	"server/modules/context"
)

func Setting(ctx *context.Context) {
	var setting struct {
		Features    map[string]uint32 `json:"features"`
		FieldTypes  []string          `json:"field_types"`
		FormatTypes []string          `json:"format_types"`
	}

	setting.Features = table.FeatureSetting()
	setting.FieldTypes = table.FieldTypeList()
	setting.FormatTypes = models.FormatSetting()

	ctx.Ack(setting)
}
