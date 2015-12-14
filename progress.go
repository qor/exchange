package exchange

type Progress struct {
	Current uint
	Total   uint
	Cells   []struct {
		Value interface{}
		Error error
	}
}

func (progress Progress) GetCurrent() uint {
	return progress.Current
}

func (progress Progress) GetTotal() uint {
	return progress.Current
}

func (progress Progress) GetCells() []struct {
	Value interface{}
	Error error
} {
	return progress.Cells
}
