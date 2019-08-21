package export

import (
	"archive/zip"
	"bytes"
	"io"
)

type ZipBuilder struct {
	buffer *bytes.Buffer
	writer *zip.Writer
}

func NewZipBuilder() *ZipBuilder {
	builder := &ZipBuilder{}
	builder.buffer = new(bytes.Buffer)
	builder.writer = zip.NewWriter(builder.buffer)
	return builder
}

func (b *ZipBuilder) GetFileWriter(filename string) (io.Writer, error) {
	return b.writer.Create(filename)
}

func (b *ZipBuilder) Write(w io.Writer) error {
	if err := b.writer.Close(); err != nil {
		return err
	}

	_, err := b.buffer.WriteTo(w)
	return err
}

func (b *ZipBuilder) Buffer() (*bytes.Buffer, error) {
	if err := b.writer.Close(); err != nil {
		return nil, err
	}

	return b.buffer, nil
}