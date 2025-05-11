package fixme

import (
	"testing"
)

func TestParallelSum(t *testing.T) {
	tests := []struct {
		name          string
		arr           []int
		numGoroutines int
		want          int
	}{
		{
			name:          "Empty array",
			arr:           []int{},
			numGoroutines: 4,
			want:          0,
		},
		{
			name:          "Small array",
			arr:           []int{1, 2, 3, 4, 5},
			numGoroutines: 2,
			want:          15,
		},
		{
			name:          "Large array",
			arr:           []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			numGoroutines: 3,
			want:          55,
		},
		{
			name:          "More goroutines than elements",
			arr:           []int{1, 2, 3},
			numGoroutines: 5,
			want:          6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParallelSum(tt.arr, tt.numGoroutines); got != tt.want {
				t.Errorf("ParallelSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
