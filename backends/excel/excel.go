package excel

import (
	"errors"
	"io"
	"os"
)

// New new excel backend
func New(value interface{}, config ...Config) *Excel {
	excel := &Excel{filename: name}

	if f, ok := value.(string); ok {
		excel.filename = f
	} else {
		if r, ok := value.(io.ReadCloser); ok {
			excel.reader = r
		}
	}

	for _, cfg := range config {
		excel.config = cfg
	}
	return excel
}

// Config excel config
type Config struct {
}

// Excel excel struct
type Excel struct {
	filename string
	reader   io.ReadCloser
}

func (excel *Excel) getReader() (io.ReadCloser, error) {
	if excel.reader != nil {
		return excel.reader, nil
	} else if excel.filename != "" {
		readerCloser, err := os.Open(excel.filename)
		return readerCloser, err
	}

	return nil, errors.New("Nothing available to import")
}
