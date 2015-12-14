package exchange

type Progress struct {
	Current uint
	Total   uint
	Cells   []struct {
		Value  interface{}
		Header bool
		Error  error
	}
}
