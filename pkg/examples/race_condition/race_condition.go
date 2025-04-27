package race_condition

import (
	"fmt"
	"sync"
)

// RunRaceConditionProblem демонстрирует проблему гонки данных
func RunRaceConditionProblem() int {
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // Проблема: гонка данных
		}()
	}

	wg.Wait()
	fmt.Println("Итоговое значение (с гонкой данных):", counter)
	// Результат будет непредсказуемым и почти всегда меньше 1000
	return counter
}

// RunRaceConditionSolution демонстрирует решение проблемы гонки данных с помощью каналов
func RunRaceConditionSolution() int {
	counter := 0
	ch := make(chan int, 1000)
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- 1 // Отправляем сигнал увеличения
		}()
	}

	// Дожидаемся завершения горутин
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Получаем все сигналы и увеличиваем счетчик
	for _ = range ch {
		counter++
	}

	fmt.Println("Итоговое значение (решение с каналами):", counter)
	// Результат всегда будет 1000
	return counter
}

// RunRaceConditionMutex демонстрирует решение проблемы гонки данных с помощью мьютекса
func RunRaceConditionMutex() int {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // Безопасное увеличение с мьютексом
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println("Итоговое значение (решение с мьютексом):", counter)
	// Результат всегда будет 1000
	return counter
}
