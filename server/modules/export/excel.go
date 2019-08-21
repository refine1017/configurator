package export

import (
	"bytes"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
)

type TableExcelBuilder struct {
	file     *excelize.File
	sheet    string
	index    int
	fields   []string
}

func NewTableExcelBuilder() *TableExcelBuilder {
	builder := &TableExcelBuilder{}
	builder.init()
	return builder
}

func (b *TableExcelBuilder) init() {
	b.file = excelize.NewFile()
	b.sheet = "Sheet1"
	b.index = 1
}

func (b *TableExcelBuilder) SetFields(fields []string) error {
	for i, field := range fields {
		axis := fmt.Sprintf("%s%d", string([]rune{65 + int32(i)}), b.index)

		err := b.file.SetCellValue(b.sheet, axis, field)
		if err != nil {
			return err
		}
	}

	b.fields = fields

	return nil
}

func (b *TableExcelBuilder) AppendMapValues(values map[string]interface{}) error {
	b.index++

	var i = 0

	for _, key := range b.fields {
		axis := fmt.Sprintf("%s%d", string([]rune{65 + int32(i)}), b.index)

		err := b.file.SetCellValue(b.sheet, axis, values[key])
		if err != nil {
			return err
		}

		i++
	}

	return nil
}

func (b *TableExcelBuilder) Write(w io.Writer) error {
	return b.file.Write(w)
}

func (b *TableExcelBuilder) Buffer() (*bytes.Buffer, error) {
	var buf = new(bytes.Buffer)
	return buf, b.Write(buf)
}