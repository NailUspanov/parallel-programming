package fixme

import (
	"reflect"
	"testing"
)

func TestProcessData(t *testing.T) {
	tests := []struct {
		name string
		data []int
		want []int
	}{
		{
			name: "Empty slice",
			data: []int{},
			want: []int{},
		},
		{
			name: "Single element",
			data: []int{5},
			want: []int{10},
		},
		{
			name: "Multiple elements",
			data: []int{1, 2, 3, 4, 5},
			want: []int{2, 4, 6, 8, 10},
		},
		{
			name: "Negative numbers",
			data: []int{-3, -2, -1, 0, 1},
			want: []int{-6, -4, -2, 0, 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ProcessData(tt.data)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProcessData() = %v, want %v", got, tt.want)
			}
		})
	}
}
