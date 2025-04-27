package waitgroups

import (
	"fmt"
	"sync"
)

// RunWaitGroups демонстрирует базовое использование WaitGroup
func RunWaitGroups() []string {
	var wg sync.WaitGroup
	results := make([]string, 0, 2)
	resultCh := make(chan string, 2)

	// Увеличиваем счетчик на 2
	wg.Add(2)

	go func() {
		// Уменьшаем счетчик при завершении
		defer wg.Done()
		resultCh <- "Горутина 1"
	}()

	go func() {
		defer wg.Done()
		resultCh <- "Горутина 2"
	}()

	// Запускаем горутину для сбора результатов
	go func() {
		// Ожидаем завершения всех горутин
		wg.Wait()
		// Закрываем канал результатов
		close(resultCh)
	}()

	// Собираем результаты
	for result := range resultCh {
		fmt.Println(result)
		results = append(results, result)
	}

	fmt.Println("Все горутины завершены")
	return results
}
