package race_condition

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// RunRaceConditionProblem демонстрирует проблему гонки данных
func RunRaceConditionProblem(incPerRoutine int) int {
	counter := 0
	var wg sync.WaitGroup
	numGoroutines := 2 * incPerRoutine // 2 горутины на каждый инкремент

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // Проблема: гонка данных
		}()
	}

	wg.Wait()
	fmt.Printf("Итоговое значение (с гонкой данных): %d\n", counter)
	fmt.Printf("Ожидаемое значение: %d\n", numGoroutines)
	// Результат будет непредсказуемым и почти всегда меньше ожидаемого
	return counter
}

// RunRaceConditionSolutionWithChannels демонстрирует решение проблемы гонки данных с помощью каналов
func RunRaceConditionSolutionWithChannels(incPerRoutine int) int {
	counter := 0
	numGoroutines := 2 * incPerRoutine
	ch := make(chan int, numGoroutines)
	var wg sync.WaitGroup

	fmt.Println("Решение с использованием каналов")
	for i := 0; i < numGoroutines; i++ {
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

	fmt.Printf("Итоговое значение (решение с каналами): %d\n", counter)
	// Результат всегда будет numGoroutines
	return counter
}

// RunRaceConditionSolutionWithMutex демонстрирует решение проблемы гонки данных с помощью мьютекса
func RunRaceConditionSolutionWithMutex(incPerRoutine int) int {
	counter := 0
	var wg sync.WaitGroup
	var mu sync.Mutex
	numGoroutines := 2 * incPerRoutine

	fmt.Println("Решение с использованием мьютекса")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++ // Безопасное увеличение с мьютексом
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Printf("Итоговое значение (решение с мьютексом): %d\n", counter)
	// Результат всегда будет numGoroutines
	return counter
}

// RunRaceConditionSolutionWithAtomic демонстрирует решение проблемы гонки данных с помощью атомарных операций
func RunRaceConditionSolutionWithAtomic(incPerRoutine int) int {
	var counter int64 = 0
	var wg sync.WaitGroup
	numGoroutines := 2 * incPerRoutine

	fmt.Println("Решение с использованием атомарных операций")
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // Атомарное увеличение счетчика
		}()
	}

	wg.Wait()
	fmt.Printf("Итоговое значение (с атомарными операциями): %d\n", counter)
	// Результат всегда будет numGoroutines
	return int(counter)
}

// RunAtomicCompareAndSwap демонстрирует использование атомарной операции CompareAndSwap
func RunAtomicCompareAndSwap() (bool, int64, int64) {
	var counter int64 = 42
	fmt.Printf("Начальное значение: %d\n", counter)

	// Атомарная операция сравнения и обмена
	// Меняем на 100, если текущее значение равно 42
	oldValue := counter
	newValue := int64(100)
	success := atomic.CompareAndSwapInt64(&counter, oldValue, newValue)

	if success {
		fmt.Printf("Успешная операция: %d -> %d\n", oldValue, newValue)
	} else {
		fmt.Printf("Неудачная операция: текущее значение не равно %d\n", oldValue)
	}

	fmt.Printf("Старое значение: %d\n", oldValue)
	fmt.Printf("Новое значение: %d\n", counter)

	return success, oldValue, counter
}

// RunAtomicLoad демонстрирует использование атомарных операций Load и Store
func RunAtomicLoad() int {
	var counter atomic.Int64
	counter.Store(100)
	var wg sync.WaitGroup

	fmt.Println("Запуск горутин для атомарного чтения/записи")

	// Запускаем горутины для одновременного чтения
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Атомарное чтение значения
			value := counter.Load()
			fmt.Printf("Горутина [%d]: Прочитано значение: %d\n", id, value)
		}(i)
	}

	// Запускаем горутину для изменения значения
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Атомарно сохраняем новое значение
		newValue := int64(200)
		counter.Store(newValue)
		fmt.Printf("Значение изменено на: %d\n", newValue)
	}()

	wg.Wait()
	finalValue := counter.Load()
	fmt.Printf("Итоговое значение: %d\n", finalValue)
	return int(finalValue)
}

// RunRaceConditionSolution демонстрирует все решения проблемы гонки данных
func RunRaceConditionSolution(incPerRoutine int) []int {
	fmt.Println("Запуск всех решений проблемы гонки данных")

	channelResult := RunRaceConditionSolutionWithChannels(incPerRoutine)
	mutexResult := RunRaceConditionSolutionWithMutex(incPerRoutine)
	atomicResult := RunRaceConditionSolutionWithAtomic(incPerRoutine)

	return []int{channelResult, mutexResult, atomicResult}
}
