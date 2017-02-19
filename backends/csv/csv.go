package csv

// New initialize CSV backend, config is option, the last one will be used if there are more than one configs
func New(filename string, config ...Config) *CSV {
	csv := &CSV{Filename: filename}
	for _, cfg := range config {
		csv.config = cfg
	}
	return csv
}

type Config struct {
	TrimSpace bool
}

type CSV struct {
	Filename string
	records  [][]string
	config   Config
}
