package excel

import "testing"

func TestToAxis(t *testing.T) {
	results := map[string][]int{
		"A1":  []int{1, 1},
		"A5":  []int{1, 5},
		"Z1":  []int{26, 1},
		"Z10": []int{26, 10},
		"AA1": []int{27, 1},
		"AZ9": []int{52, 9},
		"BA8": []int{53, 8},
	}

	for key, value := range results {
		if axis := toAxis(value[0], value[1]); axis != key {
			t.Errorf("%v's axis should be %v, but got %v", value, key, axis)
		}
	}
}
