package solutions

import (
	"math"
	"sync"
)

// MinMax представляет пару минимального и максимального значений
type MinMax struct {
	Min int
	Max int
}

// FindMinMax находит минимальный и максимальный элементы в срезе
// используя указанное количество горутин для параллельной обработки
func FindMinMax(numbers []int, numGoroutines int) (min, max int) {
	if len(numbers) == 0 {
		return 0, 0
	}

	if numGoroutines <= 0 {
		numGoroutines = 1
	}

	// Ограничиваем количество горутин размером среза
	if numGoroutines > len(numbers) {
		numGoroutines = len(numbers)
	}

	// Вычисляем размер части для каждой горутины
	chunkSize := (len(numbers) + numGoroutines - 1) / numGoroutines

	// Канал для сбора результатов
	resultChan := make(chan MinMax, numGoroutines)
	var wg sync.WaitGroup

	// Запускаем горутины для поиска локальных минимумов и максимумов
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Вычисляем границы части для текущей горутины
			start := goroutineID * chunkSize
			end := start + chunkSize
			if end > len(numbers) {
				end = len(numbers)
			}

			// Если начало за пределами массива, ничего не делаем
			if start >= len(numbers) {
				return
			}

			// Инициализируем локальные минимум и максимум первым элементом части
			localMin := numbers[start]
			localMax := numbers[start]

			// Находим локальные минимум и максимум в части
			for j := start + 1; j < end; j++ {
				if numbers[j] < localMin {
					localMin = numbers[j]
				}
				if numbers[j] > localMax {
					localMax = numbers[j]
				}
			}

			// Отправляем результат в канал
			resultChan <- MinMax{Min: localMin, Max: localMax}
		}(i)
	}

	// Ждем завершения всех горутин в отдельной горутине и закрываем канал результатов
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Инициализируем глобальные минимум и максимум
	globalMin := math.MaxInt32
	globalMax := math.MinInt32

	// Обрабатываем результаты из канала
	for result := range resultChan {
		if result.Min < globalMin {
			globalMin = result.Min
		}
		if result.Max > globalMax {
			globalMax = result.Max
		}
	}

	return globalMin, globalMax
}
