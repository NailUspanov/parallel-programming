package combined

import (
	"fmt"
	"sync"
	"time"
)

// worker обрабатывает задания из канала jobs и отправляет результаты в канал results
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Воркер %d начал обработку задания %d\n", id, j)
		time.Sleep(100 * time.Millisecond) // Имитация работы
		fmt.Printf("Воркер %d завершил обработку задания %d\n", id, j)
		results <- j * 2 // Отправляем результат
	}
}

// RunCombined демонстрирует комбинированное использование каналов и WaitGroups
func RunCombined() []int {
	var wg sync.WaitGroup
	jobs := make(chan int, 100)
	results := make(chan int, 100)
	collectedResults := make([]int, 0, 10)

	// Запускаем воркеров
	for w := 1; w <= 3; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Отправляем задачи
	for j := 1; j <= 10; j++ {
		jobs <- j
	}
	close(jobs) // Больше заданий не будет

	// Запускаем горутину для закрытия канала результатов
	// после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	for r := range results {
		fmt.Println("Результат:", r)
		collectedResults = append(collectedResults, r)
	}

	return collectedResults
}
