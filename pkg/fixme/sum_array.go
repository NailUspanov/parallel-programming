package fixme

import (
	"sync"
)

// ParallelSum неправильно суммирует элементы массива, используя горутины и каналы
// Функция должна разделить массив на части, обработать каждую часть в отдельной горутине,
// собрать результаты через канал и вернуть общую сумму
func ParallelSum(arr []int, numGoroutines int) int {
	if len(arr) == 0 {
		return 0
	}

	// Ошибка: канал закрывается слишком рано
	results := make(chan int, numGoroutines)
	var wg sync.WaitGroup

	// Вычисляем размер части для каждой горутины
	chunkSize := (len(arr) + numGoroutines - 1) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()

			// Вычисляем границы части для текущей горутины
			start := goroutineID * chunkSize
			end := start + chunkSize
			if end > len(arr) {
				end = len(arr)
			}

			// Если начало за пределами массива, ничего не делаем
			if start >= len(arr) {
				return
			}

			// Суммируем элементы в части
			sum := 0
			for j := start; j < end; j++ {
				sum += arr[j]
			}

			// Отправляем результат в канал
			results <- sum
		}(i)
	}

	// Ошибка: канал закрывается до завершения всех горутин
	close(results)

	// Ждем завершения всех горутин
	wg.Wait()

	// Суммируем результаты из канала
	totalSum := 0
	for sum := range results {
		totalSum += sum
	}

	return totalSum
}
