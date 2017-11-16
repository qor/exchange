package excel

import "io"

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

// Excel excel struct
type Excel struct {
	filename string
	reader   io.ReadCloser
}

// Config excel config
type Config struct {
}
