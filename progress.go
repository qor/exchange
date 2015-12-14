package exchange

type Progress struct {
	Current uint
	Total   uint
	Cells   []struct {
		Value interface{}
		Error error
	}
}
