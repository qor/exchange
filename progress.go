package exchange

type ImportProgress struct {
	Current uint
	Total   uint
	Cells   []ImportProgressCell
}

type ImportProgressCell struct {
	Header string
	Value  interface{}
	Error  error
}

type ExportProgress struct {
	Current uint
	Total   uint
	Value   interface{}
}
