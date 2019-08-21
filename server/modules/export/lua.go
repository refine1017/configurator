package export

import (
	"bytes"
	"fmt"
	"server/models/table"
	"strconv"
	"strings"
)

type TableLuaBuilder struct {
	exportType string
	buf        *bytes.Buffer
}

func NewTableLuaBuilder() *TableLuaBuilder {
	builder := &TableLuaBuilder{}
	builder.buf = new(bytes.Buffer)

	builder.buf.WriteString(`local data = {
`)

	return builder
}

func (b *TableLuaBuilder) AppendKVValue(fields map[string]string, values map[string]interface{}) {
	for key, value := range values {
		b.buf.WriteString(fmt.Sprintf(`["%s"]={key="%s",value="%s"},`+"\n", key, key, value))
	}
}

func (b *TableLuaBuilder) AppendArrValue(fields map[string]string, values map[string]interface{}) {
	b.transfer(fields, values)
}

func (b *TableLuaBuilder) transfer(fields map[string]string, values map[string]interface{}) map[string]interface{} {
	newValue := make(map[string]interface{})

	if id, ok := values["id"]; ok {
		//v, _ := strconv.ParseInt(fmt.Sprintf(`%f`, id), 10, 64)
		b.buf.WriteString(fmt.Sprintf(`["%.0f"]={`, id))
	}

	for field, _type := range fields {
		value := values[field]

		b.buf.WriteString(fmt.Sprintf("%s=%v,", field, b.strconv(value, _type)))
	}

	b.buf.WriteString("},\n")

	return newValue
}

func (b *TableLuaBuilder) strconv(value interface{}, _type string) interface{} {
	valueStr := fmt.Sprintf(`%v`, value)

	switch _type {
	case table.TypeUint32, table.TypeUint64, table.TypeInt32, table.TypeInt64:
		v, _ := strconv.Atoi(fmt.Sprintf("%.0f", value))
		return v
	case table.TypeFloat:
		v, _ := strconv.ParseFloat(valueStr, 32)
		return v
	case table.TypeDouble:
		v, _ := strconv.ParseFloat(valueStr, 64)
		return v
	case table.TypeBool:
		v, _ := strconv.ParseBool(valueStr)
		return v
	case table.TypeBoolArray, table.TypeStringArray, table.TypeUint32Array, table.TypeUint64Array, table.TypeInt32Array, table.TypeInt64Array, table.TypeFloatArray, table.TypeDoubleArray:
		if value == nil || strings.TrimSpace(valueStr) == "" {
			return "{}"
		}

		return fmt.Sprintf(`{%v}`, valueStr)
	default:
		valueStr = strings.Replace(valueStr, "\n\r", "", -1)
		valueStr = strings.Replace(valueStr, "\n", "", -1)
		valueStr = strings.Replace(valueStr, "\r", "", -1)
		return fmt.Sprintf(`"%v"`, valueStr)
	}
}

func (b *TableLuaBuilder) Write() {
	b.buf.WriteString(`}return data`)
	return
}

func (b *TableLuaBuilder) Buffer() (*bytes.Buffer, error) {
	b.Write()
	return b.buf, nil
}
