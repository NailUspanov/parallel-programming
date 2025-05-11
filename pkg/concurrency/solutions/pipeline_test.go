package solutions

import (
	"testing"
)

func TestProcessPipeline(t *testing.T) {
	tests := []struct {
		name    string
		numbers []int
		want    int
	}{
		{
			name:    "Empty slice",
			numbers: []int{},
			want:    0,
		},
		{
			name:    "No even squares",
			numbers: []int{1, 3, 5, 7},
			want:    0, // Квадраты: 1, 9, 25, 49 - все нечетные
		},
		{
			name:    "All even squares",
			numbers: []int{2, 4, 6, 8},
			want:    120, // Квадраты: 4, 16, 36, 64 - все четные, сумма = 120
		},
		{
			name:    "Mixed even and odd squares",
			numbers: []int{1, 2, 3, 4, 5},
			want:    20, // Квадраты: 1, 4, 9, 16, 25 - четные: 4, 16, сумма = 20
		},
		{
			name:    "Negative numbers",
			numbers: []int{-1, -2, -3, -4},
			want:    20, // Квадраты: 1, 4, 9, 16 - четные: 4, 16, сумма = 20
		},
		{
			name:    "Zero included",
			numbers: []int{0, 1, 2, 3},
			want:    4, // Квадраты: 0, 1, 4, 9 - четные: 0, 4, сумма = 4
		},
		{
			name:    "Larger dataset",
			numbers: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			want:    220, // Четные квадраты: 4, 16, 36, 64, 100, сумма = 220
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProcessPipeline(tt.numbers); got != tt.want {
				t.Errorf("ProcessPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}
