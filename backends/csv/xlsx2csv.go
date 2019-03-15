package csv

import (
	"bytes"
	stdcsv "encoding/csv"
	"errors"
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

	for _, row := range sheet.Rows {
		if row.Hidden {
			continue
		}

		var record = make([]string, len(row.Cells))
		for i, cell := range row.Cells {
			record[i] = cell.Value
		}
		err = csvWriter.Write(record)
		if err != nil {
			return nil, err
		}
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
