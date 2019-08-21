package models

import (
	"bytes"
	"fmt"
	"server/modules/export"
	"server/modules/util"
	"strings"
)

type ConfigExporter func(env *Environment, cfg *Config) (*bytes.Buffer, error)

const (
	ExtensionJson  = ".json"
	ExtensionLua   = ".lua"
	ExtensionExcel = ".xlsx"
	ExtensionProto = ".proto"
	ExtensionZip   = ".zip"
	ExtensionCS    = ".cs"
)

const (
	ExportTypeJson    = "json"
	ExportTypeExcel   = "excel"
	ExportTypeProto   = "proto"
	ExportTypeEntitas = "entitas"
	ExportTypeCS      = "cs"
	ExportTypeLua     = "lua"
)

const ProtoPrefix = "Cfg"

var exporters = map[string]map[string]ConfigExporter{
	FormatTable: {
		ExportTypeJson:    ExportTableJsonFile,
		ExportTypeLua:     ExportTableLuaFile,
		ExportTypeExcel:   ExportTableExcelFile,
		ExportTypeProto:   ExportTableProtoFile,
		ExportTypeEntitas: ExportTableEntitasFile,
		ExportTypeCS:      ExportTableCSFile,
	},
	FormatKV: {
		ExportTypeJson:    ExportKVJsonFile,
		ExportTypeLua:     ExportKVLuaFile,
		ExportTypeExcel:   ExportKVExcelFile,
		ExportTypeProto:   ExportKVProtoFile,
		ExportTypeEntitas: ExportKVEntitasFile,
		ExportTypeCS:      ExportKVCSFile,
	},
	FormatJson: {
		ExportTypeProto:   ExportJsonProtoFile,
		ExportTypeEntitas: ExportJsonEntitasFile,
		ExportTypeCS:      ExportJsonCSFile,
	},
}

func ExportEnvJsonZip(env *Environment, config string) (*bytes.Buffer, error) {
	var buf *bytes.Buffer
	var err error

	builder := export.NewZipBuilder()

	config = strings.ToLower(config)

	for _, cfg := range env.Configs {
		if config != "all" && strings.ToLower(cfg.Name) != config {
			continue
		}

		if cfg.Format == FormatJson {
			if err := ExportJsonDir(env, cfg, builder); err != nil {
				return nil, err
			}
			continue
		}

		exporter := exporters[cfg.Format][ExportTypeJson]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		writer, err := builder.GetFileWriter(cfg.Name + ExtensionJson)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteTo(writer)
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportEnvLuaZip(env *Environment, config string) (*bytes.Buffer, error) {
	var buf *bytes.Buffer
	var err error

	builder := export.NewZipBuilder()

	config = strings.ToLower(config)

	for _, cfg := range env.Configs {
		if config != "all" && strings.ToLower(cfg.Name) != config {
			continue
		}

		if cfg.Format == FormatJson {
			if err := ExportJsonDir(env, cfg, builder); err != nil {
				return nil, err
			}
			continue
		}

		exporter := exporters[cfg.Format][ExportTypeLua]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		writer, err := builder.GetFileWriter(cfg.Name + ExtensionLua)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteTo(writer)
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportEnvExcelZip(env *Environment) (*bytes.Buffer, error) {
	var buf *bytes.Buffer
	var err error

	builder := export.NewZipBuilder()

	for _, cfg := range env.Configs {
		if cfg.Format == FormatJson {
			if err := ExportJsonDir(env, cfg, builder); err != nil {
				return nil, err
			}
			continue
		}

		exporter := exporters[cfg.Format][ExportTypeExcel]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		writer, err := builder.GetFileWriter(cfg.Name + ExtensionExcel)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteTo(writer)
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportEnvProto(env *Environment) (*bytes.Buffer, error) {
	var buffer = new(bytes.Buffer)

	var err error
	var buf *bytes.Buffer

	buffer.WriteString(`syntax = "proto3";
package pb;

`)

	builder := export.NewProtoMessageBuilder("All", ProtoPrefix, "All Configs")

	for _, cfg := range env.Configs {
		exporter := exporters[cfg.Format][ExportTypeProto]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		buffer.Write(buf.Bytes())
		buffer.WriteString("\n")

		switch cfg.Format {
		case FormatKV:
			builder.AddField(ProtoPrefix+util.CamelCase(cfg.Name), cfg.Name+"s", cfg.Name+" value")
		default:
			builder.AddMessageArray(ProtoPrefix+util.CamelCase(cfg.Name), cfg.Name+"s", cfg.Name+" list")
		}
	}

	buffer.Write(builder.Buffer().Bytes())

	return buffer, nil
}

func ExportEntitasProto(env *Environment) (*bytes.Buffer, error) {
	var err error
	var buf *bytes.Buffer

	builder := export.NewZipBuilder()

	for _, cfg := range env.Configs {
		exporter := exporters[cfg.Format][ExportTypeEntitas]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		filename := fmt.Sprintf("%sConfigComponent%s", util.CamelCase(cfg.Name), ExtensionCS)

		writer, err := builder.GetFileWriter(filename)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteTo(writer)
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}

func ExportCSProto(env *Environment) (*bytes.Buffer, error) {
	var err error
	var buf *bytes.Buffer

	builder := export.NewZipBuilder()

	for _, cfg := range env.Configs {
		exporter := exporters[cfg.Format][ExportTypeCS]
		buf, err = exporter(env, cfg)

		if err != nil {
			return nil, err
		}

		filename := fmt.Sprintf("%sConfig%s", util.CamelCase(cfg.Name), ExtensionCS)

		writer, err := builder.GetFileWriter(filename)
		if err != nil {
			return nil, err
		}

		_, err = buf.WriteTo(writer)
		if err != nil {
			return nil, err
		}
	}

	return builder.Buffer()
}
