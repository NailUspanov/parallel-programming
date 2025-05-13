package solutions

import (
	"sync"
)

// ProcessPipeline обрабатывает срез чисел через конвейер из трех стадий:
// 1. Вычисление квадратов чисел
// 2. Фильтрация (оставляем только четные числа)
// 3. Суммирование результатов
func ProcessPipeline(numbers []int) int {
	// Создаем каналы для передачи данных между стадиями
	stage1Chan := make(chan int) // Квадраты
	stage2Chan := make(chan int) // Отфильтрованные четные числа
	resultChan := make(chan int) // Результат суммирования

	// Объект WaitGroup для ожидания завершения всех стадий
	var wg sync.WaitGroup
	wg.Add(3) // Три стадии конвейера

	// Стадия 1: Вычисление квадратов чисел
	go func() {
		defer wg.Done()
		defer close(stage1Chan)

		for _, num := range numbers {
			stage1Chan <- num * num
		}
	}()

	// Стадия 2: Фильтрация (оставляем только четные числа)
	go func() {
		defer wg.Done()
		defer close(stage2Chan)

		for num := range stage1Chan {
			if num%2 == 0 {
				stage2Chan <- num
			}
		}
	}()

	// Стадия 3: Суммирование результатов
	go func() {
		defer wg.Done()

		sum := 0
		for num := range stage2Chan {
			sum += num
		}

		resultChan <- sum
		close(resultChan)
	}()

	// Запускаем отдельную горутину для ожидания завершения всех стадий
	// и закрытия канала результатов (если это еще не сделано)
	go func() {
		wg.Wait()
		// Проверяем, закрыт ли канал результатов
		select {
		case <-resultChan:
			// Канал уже закрыт или содержит значение
		default:
			// Канал еще открыт и пуст, значит что-то пошло не так
			// Закрываем канал с результатом 0
			close(resultChan)
		}
	}()

	// Получаем результат из канала
	// Если канал был закрыт без отправки результата, вернется нулевое значение
	result, ok := <-resultChan
	if !ok {
		// Канал был закрыт без отправки результата
		return 0
	}

	return result
}
