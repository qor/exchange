package csv

import (
	"bytes"
	stdcsv "encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/tealeg/xlsx"
)

func generateCSVFromXLSXFile(fileName string) (io.ReadCloser, error) {
	xlFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	if len(xlFile.Sheets) == 0 {
		return nil, errors.New("This XLSX file contains no sheets")
	}
	sheet := xlFile.Sheets[0]

	var buf bytes.Buffer
	csvWriter := stdcsv.NewWriter(&buf)

	var firstRowSize int
	err = sheet.ForEachRow(func(row *xlsx.Row) error {
		if row.Hidden {
			return nil
		}

		var record []string
		row.ForEachCell(func(cell *xlsx.Cell) error {
			record = append(record, cell.Value)
			return nil
		})

		if len(record) == 0 {
			return nil
		}

		if firstRowSize == 0 {
			firstRowSize = len(record)
		}

		if firstRowSize != len(record) {
			return errors.New(fmt.Sprintf("This XLSX file data is invalid,Header length:%d,This row length:%d,This row data:%s", firstRowSize, len(record), record))
		}

		err = csvWriter.Write(record)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		return nil, err
	}

	return ioutil.NopCloser(&buf), nil
}

func isXLSXFile(fileName string) bool {
	return filepath.Ext(fileName) == ".xlsx"
}
