package solutions

import (
	"testing"
)

func TestFindMinMax(t *testing.T) {
	tests := []struct {
		name          string
		numbers       []int
		numGoroutines int
		wantMin       int
		wantMax       int
	}{
		{
			name:          "Empty slice",
			numbers:       []int{},
			numGoroutines: 3,
			wantMin:       0,
			wantMax:       0,
		},
		{
			name:          "Single element",
			numbers:       []int{5},
			numGoroutines: 2,
			wantMin:       5,
			wantMax:       5,
		},
		{
			name:          "Multiple elements",
			numbers:       []int{4, 2, 7, 1, 9},
			numGoroutines: 2,
			wantMin:       1,
			wantMax:       9,
		},
		{
			name:          "More goroutines than elements",
			numbers:       []int{3, 1, 4},
			numGoroutines: 5,
			wantMin:       1,
			wantMax:       4,
		},
		{
			name:          "Negative numbers",
			numbers:       []int{-5, -3, -8, -1, -6},
			numGoroutines: 3,
			wantMin:       -8,
			wantMax:       -1,
		},
		{
			name:          "Mixed positive and negative",
			numbers:       []int{-10, 5, 0, -3, 7, -8, 2},
			numGoroutines: 4,
			wantMin:       -10,
			wantMax:       7,
		},
		{
			name:          "Single goroutine",
			numbers:       []int{5, 2, 8, 1, 9, 3},
			numGoroutines: 1,
			wantMin:       1,
			wantMax:       9,
		},
		{
			name:          "Zero or negative goroutines",
			numbers:       []int{5, 2, 8, 1, 9, 3},
			numGoroutines: 0,
			wantMin:       1,
			wantMax:       9,
		},
		{
			name:          "Large dataset",
			numbers:       generateLargeDataset(1000),
			numGoroutines: 10,
			wantMin:       0,
			wantMax:       999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMin, gotMax := FindMinMax(tt.numbers, tt.numGoroutines)
			if gotMin != tt.wantMin {
				t.Errorf("FindMinMax() gotMin = %v, want %v", gotMin, tt.wantMin)
			}
			if gotMax != tt.wantMax {
				t.Errorf("FindMinMax() gotMax = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}

// generateLargeDataset создает большой массив чисел от 0 до size-1
func generateLargeDataset(size int) []int {
	result := make([]int, size)
	for i := 0; i < size; i++ {
		result[i] = i
	}
	return result
}
