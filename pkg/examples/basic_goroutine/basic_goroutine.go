package basic_goroutine

import (
	"fmt"
	"time"
)

// RunBasicGoroutine демонстрирует базовый пример использования горутин
func RunBasicGoroutine() {
	go func() {
		// Код, который будет выполняться в горутине
		fmt.Println("Это выполняется в горутине")
	}()

	// Код основной функции продолжает выполняться параллельно
	fmt.Println("Это выполняется в основной функции")

	// Даем горутине время выполниться
	time.Sleep(100 * time.Millisecond)
}

// RunMultipleGoroutines запускает несколько горутин и собирает их результаты
func RunMultipleGoroutines(count int) []int {
	resultChan := make(chan int, count)
	results := make([]int, 0, count)

	for i := 0; i < count; i++ {
		go func(id int) {
			// Имитация работы разной длительности
			time.Sleep(time.Duration(id*10) * time.Millisecond)
			fmt.Printf("Горутина %d завершила работу\n", id)
			resultChan <- id * id // Отправляем квадрат номера горутины
		}(i)
	}

	// Собираем результаты от всех горутин
	for i := 0; i < count; i++ {
		result := <-resultChan
		results = append(results, result)
	}

	fmt.Printf("Получены результаты от %d горутин\n", count)
	return results
}
