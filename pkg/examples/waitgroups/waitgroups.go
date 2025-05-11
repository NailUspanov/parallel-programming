package waitgroups

import (
	"fmt"
	"sync"
	"time"
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

// RunDynamicWaitGroups демонстрирует динамическое добавление горутин в WaitGroup
func RunDynamicWaitGroups(initial, additional int) []string {
	var wg sync.WaitGroup
	results := make([]string, 0, initial+additional)
	resultCh := make(chan string, initial+additional)
	additionalDone := make(chan bool)

	// Функция запуска горутины с указанным ID
	startWorker := func(id int, delay time.Duration) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(delay)
			resultCh <- fmt.Sprintf("Горутина %d", id)
		}()
	}

	// Запускаем начальные горутины
	for i := 0; i < initial; i++ {
		startWorker(i, time.Duration(10*i)*time.Millisecond)
	}

	// Запускаем дополнительную горутину, которая добавит еще горутин
	go func() {
		// Даем время для запуска первых горутин
		time.Sleep(50 * time.Millisecond)
		fmt.Println("Добавляем дополнительные горутины")

		// Динамически добавляем еще горутин
		for i := 0; i < additional; i++ {
			startWorker(initial+i, time.Duration(5*i)*time.Millisecond)
		}
		additionalDone <- true
	}()

	// Ждем, пока все дополнительные горутины будут добавлены
	<-additionalDone

	// Запускаем горутину для ожидания и закрытия канала
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Собираем результаты
	for result := range resultCh {
		fmt.Println(result)
		results = append(results, result)
	}

	fmt.Println("Все горутины завершены, всего:", len(results))
	return results
}
