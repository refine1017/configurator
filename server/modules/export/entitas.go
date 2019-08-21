package export

import (
	"bytes"
	"fmt"
	"server/modules/util"
	"strings"
)

var entitasTypeMapping = map[string]string{
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

const (
	IndexPrimary = "primary"
	IndexNormal  = "index"
)

type EntitasCodeBuilder struct {
	buf *bytes.Buffer
}

func NewEntitasCodeBuilder(name string, prefix string, desc string, unique bool) *EntitasCodeBuilder {
	builder := &EntitasCodeBuilder{}
	builder.buf = new(bytes.Buffer)
	builder.writeMessage(name, prefix, desc, unique)
	return builder
}

func (pb *EntitasCodeBuilder) writeMessage(name string, prefix string, desc string, unique bool) {
	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	var uniqueStr = ""
	if unique {
		uniqueStr = ", Unique"
	}

	pb.buf.WriteString(fmt.Sprintf(`using System.Collections.Generic;
using Entitas.CodeGeneration.Attributes;

// %s
namespace Entitas
{
	[Config%s]
	public class %s%sConfigComponent : IComponent
	{
`, desc, uniqueStr, prefix, util.CamelCase(name)))
}

func (pb *EntitasCodeBuilder) writeField(_type string, name string, desc string) {
	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	pb.buf.WriteString(fmt.Sprintf("\t\tpublic %s %s; // %s\n", _type, util.CamelCase(name), desc))
}

func (pb *EntitasCodeBuilder) AddField(_type string, name, desc, index string) {
	if index == IndexPrimary {
		pb.buf.WriteString("\t\t[PrimaryEntityIndex]\n")
	}
	if index == IndexNormal {
		pb.buf.WriteString("\t\t[EntityIndex]\n")
	}

	pb.writeField(entitasTypeMapping[_type], name, desc)
}

func (pb *EntitasCodeBuilder) Buffer() *bytes.Buffer {
	pb.buf.WriteString(`
	}
}
`)

	return pb.buf
}
