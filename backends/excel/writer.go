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

	excelWriter, err := excel.getWriter()

	if err == nil {
		if excel.config.SheetName == "" {
			excel.config.SheetName = "Export Results"
		}

		if !excelWriter.GetSheetVisible(excel.config.SheetName) {
			excelWriter.NewSheet(excel.config.SheetName)
		}

		writer.Writer = excelWriter
	}

	return writer, err
}

// Writer CSV writer struct
type Writer struct {
	*Excel
	currentRow int
	sheetName  string
	context    *qor.Context
	metas      []*exchange.Meta
	Resource   *exchange.Resource
	Writer     *excelize.File
}

func toAxis(x, y int) string {
	return fmt.Sprintf("%v%v", excelize.ToAlphaString(x), y)
}

// WriteHeader write header
func (writer *Writer) WriteHeader() error {
	if !writer.Resource.Config.WithoutHeader {
		writer.currentRow++
		writer.Writer.InsertRow(writer.sheetName, writer.currentRow)
		for key, meta := range writer.metas {
			writer.Writer.InsertCol(writer.sheetName, excelize.ToAlphaString(key))
			writer.Writer.SetCellValue(writer.sheetName, toAxis(key, writer.currentRow), meta.Header)
		}
	}
	return nil
}

// WriteRow write row
func (writer *Writer) WriteRow(record interface{}) (*resource.MetaValues, error) {
	var metaValues resource.MetaValues
	writer.currentRow++

	for key, meta := range writer.metas {
		value := meta.GetFormattedValuer()(record, writer.context)
		metaValue := resource.MetaValue{
			Name:  meta.GetName(),
			Value: value,
		}

		metaValues.Values = append(metaValues.Values, &metaValue)
		writer.Writer.SetCellValue(writer.sheetName, toAxis(key, writer.currentRow), fmt.Sprint(value))
	}

	return &metaValues, nil
}

// Flush flush all changes
func (writer *Writer) Flush() {
	if writer.Excel.writer != nil {
		defer writer.Excel.writer.Close()
		writer.Writer.Write(writer.Excel.writer)
	}

	if writer.Excel.filename != "" {
		writer.Writer.SaveAs(writer.Excel.filename)
	}
}
