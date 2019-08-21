package export

import (
	"bytes"
	"fmt"
	"server/modules/util"
	"strings"
)

var csTypeMapping = map[string]string{
	"string":       "string",
	"text":         "string",
	"bool":         "bool",
	"float":        "float",
	"double":       "double",
	"uint32":       "uint",
	"uint64":       "ulong",
	"int64":        "long",
	"int32":        "int",
	"string_array": "List<string>",
	"bool_array":   "List<bool>",
	"uint32_array": "List<uint>",
	"uint64_array": "List<ulong>",
	"int64_array":  "List<long>",
	"int32_array":  "List<int>",
	"float_array":  "List<float>",
	"double_array": "List<double>",
}

type CSCodeBuilder struct {
	buf *bytes.Buffer
}

func NewCSCodeBuilder(name string, prefix string, desc string, unique bool) *CSCodeBuilder {
	builder := &CSCodeBuilder{}
	builder.buf = new(bytes.Buffer)
	builder.writeMessage(name, prefix, desc, unique)
	return builder
}

func (pb *CSCodeBuilder) writeMessage(name string, prefix string, desc string, unique bool) {
	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	pb.buf.WriteString(fmt.Sprintf(`using System.Collections.Generic;

// %s
namespace Config
{
	public class %s%sConfig : IConfig
	{
`, desc, prefix, util.CamelCase(name)))
}

func (pb *CSCodeBuilder) writeField(_type string, name string, desc string) {
	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	if name == "id" {
		pb.buf.WriteString("\t\tpublic uint id { get; private set; }\n")
	} else {
		pb.buf.WriteString(fmt.Sprintf("\t\tpublic %s %s; // %s\n", _type, name, desc))
	}

}

func (pb *CSCodeBuilder) AddField(_type string, name, desc, index string) {
	pb.writeField(csTypeMapping[_type], name, desc)
}

func (pb *CSCodeBuilder) Buffer() *bytes.Buffer {
	pb.buf.WriteString(`
	}
}
`)

	return pb.buf
}
