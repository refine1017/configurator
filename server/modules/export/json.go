package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"server/models/table"
	"strconv"
	"strings"
)

const (
	JsonObject = "Object"
	JsonArr    = "Array"
)

type TableJsonBuilder struct {
	exportType string
	arr        []interface{}
	obj        map[string]interface{}
}

func NewTableJsonBuilder() *TableJsonBuilder {
	builder := &TableJsonBuilder{}
	return builder
}

func (b *TableJsonBuilder) SetObjValue(fields map[string]string, values map[string]interface{}) {
	b.obj = b.transfer(fields, values)
	b.exportType = JsonObject
}

func (b *TableJsonBuilder) AppendArrValue(fields map[string]string, values map[string]interface{}) {
	b.arr = append(b.arr, b.transfer(fields, values))
	b.exportType = JsonArr
}

func (b *TableJsonBuilder) transfer(fields map[string]string, values map[string]interface{}) map[string]interface{} {
	newValue := make(map[string]interface{})

	for field, _type := range fields {
		value := values[field]

		newValue[field] = b.strconv(value, _type, false)
	}

	return newValue
}

func (b *TableJsonBuilder) strconv(value interface{}, _type string, stringValue bool) interface{} {
	valueStr := fmt.Sprintf("%v", value)

	var arrFlag = false
	var arrType = table.TypeString

	switch _type {
	case table.TypeUint32, table.TypeUint64, table.TypeInt32, table.TypeInt64:
		var v int
		if stringValue {
			v, _ = strconv.Atoi(fmt.Sprintf("%v", value))
		} else {
			v, _ = strconv.Atoi(fmt.Sprintf("%.0f", value))
		}
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
	case table.TypeBoolArray:
		arrFlag = true
		arrType = table.TypeBool
	case table.TypeStringArray:
		arrFlag = true
		arrType = table.TypeString
	case table.TypeUint32Array:
		arrFlag = true
		arrType = table.TypeUint32
	case table.TypeUint64Array:
		arrFlag = true
		arrType = table.TypeUint64
	case table.TypeInt32Array:
		arrFlag = true
		arrType = table.TypeInt32
	case table.TypeInt64Array:
		arrFlag = true
		arrType = table.TypeInt64
	case table.TypeFloatArray:
		arrFlag = true
		arrType = table.TypeFloat
	case table.TypeDoubleArray:
		arrFlag = true
		arrType = table.TypeDouble
	default:
		return valueStr
	}

	if arrFlag {
		if value == nil || strings.TrimSpace(valueStr) == "" {
			return []int32{}
		} else {
			values := strings.Split(valueStr, ",")

			arr := make([]interface{}, len(values))

			for i, v := range values {
				arr[i] = b.strconv(v, arrType, true)
			}

			return arr
		}
	}

	return valueStr
}

func (b *TableJsonBuilder) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)

	if b.exportType == JsonArr {
		return encoder.Encode(b.arr)
	} else {
		return encoder.Encode(b.obj)
	}
}

func (b *TableJsonBuilder) Buffer() (*bytes.Buffer, error) {
	var buf = new(bytes.Buffer)
	return buf, b.Write(buf)
}
