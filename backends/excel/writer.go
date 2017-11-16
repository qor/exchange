package excel

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/qor/exchange"
	"github.com/qor/qor"
	"github.com/qor/qor/resource"
	"github.com/qor/roles"
)

// NewWriter new csv writer
func (excel *Excel) NewWriter(res *exchange.Resource, context *qor.Context) (exchange.Writer, error) {
	writer := &Writer{Excel: excel, Resource: res, context: context}

	var metas []*exchange.Meta
	for _, meta := range res.Metas {
		if meta.HasPermission(roles.Read, context) {
			metas = append(metas, meta)
		}
	}
	writer.metas = metas

	excelWriter := excelize.NewFile()

	writer.Writer = excelWriter

	if excel.config.SheetName != "" {
	}

	return writer, nil
}

// Writer CSV writer struct
type Writer struct {
	*Excel
	context  *qor.Context
	Resource *exchange.Resource
	Writer   *excelize.File
	metas    []*exchange.Meta
}

// WriteHeader write header
func (writer *Writer) WriteHeader() error {
	if !writer.Resource.Config.WithoutHeader {
		var results []string
		for _, meta := range writer.metas {
			results = append(results, meta.Header)
		}
		writer.Writer.InsertRow()
		writer.Writer.SetCellValue()
		writer.Writer.Write(results)
	}
	return nil
}

// WriteRow write row
func (writer *Writer) WriteRow(record interface{}) (*resource.MetaValues, error) {
	var metaValues resource.MetaValues
	var results []string

	for _, meta := range writer.metas {
		value := meta.GetFormattedValuer()(record, writer.context)
		metaValue := resource.MetaValue{
			Name:  meta.GetName(),
			Value: value,
		}

		metaValues.Values = append(metaValues.Values, &metaValue)
		results = append(results, fmt.Sprint(value))
	}

	return &metaValues, writer.Writer.Write(results)
}

// Flush flush all changes
func (writer *Writer) Flush() {
	// FIXME write to writer
	writer.Writer.Write(writer.writer)
}
