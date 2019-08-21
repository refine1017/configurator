package export

import (
	"bytes"
	"fmt"
	"server/modules/util"
	"strings"
)

var protoTypeMapping = map[string]string{
	"fields":       "string",
	"string":       "string",
	"text":         "string",
	"bool":         "bool",
	"float":        "float",
	"double":       "double",
	"uint32":       "uint32",
	"uint64":       "uint64",
	"int64":        "int64",
	"int32":        "int32",
	"string_array": "repeated string",
	"bool_array":   "repeated bool",
	"uint32_array": "repeated uint32",
	"uint64_array": "repeated uint64",
	"int64_array":  "repeated int64",
	"int32_array":  "repeated int32",
	"float_array":  "repeated float",
	"double_array": "repeated double",
}

type ProtoMessageBuilder struct {
	buf   *bytes.Buffer
	index int
}

func NewProtoMessageBuilder(name string, prefix string, desc string) *ProtoMessageBuilder {
	builder := &ProtoMessageBuilder{}
	builder.buf = new(bytes.Buffer)
	builder.writeMessage(name, prefix, desc)
	return builder
}

func (pb *ProtoMessageBuilder) writeMessage(name string, prefix string, desc string) {
	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	pb.buf.WriteString(fmt.Sprintf("// %s\n", desc))
	pb.buf.WriteString(fmt.Sprintf("message %s%s {\n", prefix, util.CamelCase(name)))
}

func (pb *ProtoMessageBuilder) writeField(_type string, name string, desc string) {
	pb.index++

	desc = strings.Replace(desc, "\r\n", ";", -1)
	desc = strings.Replace(desc, "\n", ";", -1)

	pb.buf.WriteString(fmt.Sprintf("\t%s %s = %d; // %s\n", _type, name, pb.index, desc))
}

func (pb *ProtoMessageBuilder) AddField(_type string, name, desc string) {
	t, found := protoTypeMapping[_type]
	if found {
		pb.writeField(t, name, desc)
	} else {
		pb.writeField(_type, name, desc)
	}
}

func (pb *ProtoMessageBuilder) AddMessageArray(_type string, name, desc string) {
	pb.writeField("repeated "+_type, name, desc)
}

func (pb *ProtoMessageBuilder) Buffer() *bytes.Buffer {
	pb.buf.WriteString("}\n")

	return pb.buf
}
