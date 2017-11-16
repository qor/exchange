package excel

import (
	"errors"
	"io"
	"os"
)

// New new excel backend
func New(value interface{}, config ...*Config) *Excel {
	excel := &Excel{}

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
	TrimSpace bool
}

// Excel excel struct
type Excel struct {
	filename string
	reader   io.ReadCloser
	config   *Config
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
