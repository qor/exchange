package csv

import (
	"io"
	"os"
)

// New initialize CSV backend, config is option, the last one will be used if there are more than one configs
func New(filename string, config ...Config) *CSV {
	csv := &CSV{Filename: filename}
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
	Filename string
	records  [][]string
	config   Config
}

func (c CSV) getReader() (io.ReadCloser, error) {
	readerCloser, err := os.Open(c.Filename)
	return readerCloser, err
}
