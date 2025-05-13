package solutions

import (
	"sync"
)

// ProcessData правильно обрабатывает данные с использованием WaitGroup
// Функция обрабатывает все элементы массива в параллельных горутинах
// и дожидается завершения всех горутин перед возвратом результата
func ProcessData(data []int) []int {
	var wg sync.WaitGroup
	result := make([]int, len(data))

	for i, val := range data {
		// Добавляем счетчик WaitGroup перед запуском горутины
		wg.Add(1)
		go func(index int, value int) {
			// Обрабатываем данные (умножаем на 2)
			result[index] = value * 2
			wg.Done()
		}(i, val)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	return result
}
