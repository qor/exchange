package csv

import (
	"errors"
	"io"
	"os"
)

// New initialize CSV backend, config is option, the last one will be used if there are more than one configs
func New(value interface{}, config ...Config) *CSV {
	csv := &CSV{}

	if f, ok := value.(string); ok {
		csv.filename = f
	} else {
		if r, ok := value.(io.ReadCloser); ok {
			csv.reader = r
		}

		if w, ok := value.(io.WriteCloser); ok {
			csv = &CSV{writer: w}
		}
	}

	for _, cfg := range config {
		csv.config = cfg
	}
	return csv
}

// Config CSV exchange config
type Config struct {
	TrimSpace bool
}

// CSV CSV struct
type CSV struct {
	config  Config
	records [][]string

	filename string
	reader   io.ReadCloser
	writer   io.WriteCloser
}

func (c CSV) getReader() (io.ReadCloser, error) {
	if c.reader != nil {
		return c.reader, nil
	} else if c.filename != "" {
		readerCloser, err := os.Open(c.filename)
		return readerCloser, err
	}

	return nil, errors.New("Nothing available to import")
}

func (c CSV) getWriter() (io.WriteCloser, error) {
	if c.writer != nil {
		return c.writer, nil
	} else if c.filename != "" {
		writerCloser, err := os.Open(c.filename)
		return writerCloser, err
	}

	return nil, errors.New("Nowhere to export")
}
