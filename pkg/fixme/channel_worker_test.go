package fixme

import (
	"sort"
	"testing"
)

func TestWorker(t *testing.T) {
	tests := []struct {
		name       string
		jobs       []int
		numWorkers int
	}{
		{
			name:       "Empty jobs",
			jobs:       []int{},
			numWorkers: 3,
		},
		{
			name:       "Single job",
			jobs:       []int{5},
			numWorkers: 2,
		},
		{
			name:       "Multiple jobs",
			jobs:       []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			numWorkers: 3,
		},
		{
			name:       "One worker",
			jobs:       []int{1, 2, 3, 4, 5},
			numWorkers: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем каналы для заданий и результатов
			jobs := make(chan int, len(tt.jobs))
			results := make(chan int, len(tt.jobs))

			// Запускаем функцию Worker
			Worker(jobs, results, tt.numWorkers)

			// Отправляем задания в канал
			for _, job := range tt.jobs {
				jobs <- job
			}
			close(jobs)

			// Собираем результаты
			var receivedResults []int
			for i := 0; i < len(tt.jobs); i++ {
				select {
				case result := <-results:
					receivedResults = append(receivedResults, result)
				default:
					// Если канал пуст, а должны были получить еще результаты
					if len(receivedResults) < len(tt.jobs) {
						t.Fatalf("Received only %d results, expected %d", len(receivedResults), len(tt.jobs))
					}
				}
			}

			// Проверяем, что все результаты получены
			if len(receivedResults) != len(tt.jobs) {
				t.Errorf("Expected %d results, got %d", len(tt.jobs), len(receivedResults))
			}

			// Вычисляем ожидаемые результаты (квадраты чисел)
			var expectedResults []int
			for _, job := range tt.jobs {
				expectedResults = append(expectedResults, job*job)
			}

			// Сортируем результаты для сравнения (так как порядок может быть разным из-за параллельности)
			sort.Ints(receivedResults)
			sort.Ints(expectedResults)

			// Проверяем, что результаты совпадают
			for i, expected := range expectedResults {
				if i >= len(receivedResults) {
					t.Errorf("Missing result at index %d", i)
					continue
				}
				if receivedResults[i] != expected {
					t.Errorf("Result at index %d: got %d, want %d", i, receivedResults[i], expected)
				}
			}
		})
	}
}
