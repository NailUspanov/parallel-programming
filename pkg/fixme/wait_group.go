package fixme

import (
	"sync"
)

// ProcessData неправильно обрабатывает данные с использованием WaitGroup
// Функция должна обработать все элементы массива в параллельных горутинах
// и дождаться завершения всех горутин перед возвратом результата
func ProcessData(data []int) []int {
	var wg sync.WaitGroup
	result := make([]int, len(data))

	for i, val := range data {
		// Ошибка: отсутствует wg.Add перед запуском горутины
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
