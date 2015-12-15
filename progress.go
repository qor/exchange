package exchange

type Progress struct {
	Current uint
	Total   uint
	Value   interface{}
	Cells   []Cell
}

type Cell struct {
	Header string
	Value  interface{}
	Error  error
}
