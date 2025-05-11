package buffered_channel

import (
	"fmt"
	"sync"
	"time"
)

// RunBufferedChannel демонстрирует использование буферизованного канала
func RunBufferedChannel() []int {
	ch := make(chan int, 2) // Буферизованный канал с ёмкостью 2

	ch <- 1 // Не блокируется
	ch <- 2 // Не блокируется
	// ch <- 3 // Блокируется, пока из канала не будет прочитано значение

	value1 := <-ch // Получить первое значение (1)
	value2 := <-ch // Получить второе значение (2)

	fmt.Printf("Полученные значения: %d, %d\n", value1, value2)

	return []int{value1, value2}
}

// RunWorkerPool демонстрирует ограниченный пул горутин с буферизованным каналом
func RunWorkerPool(numWorkers, numJobs int) []int {
	jobs := make(chan int, numJobs)          // Буферизованный канал для заданий
	results := make(chan int, numJobs)       // Буферизованный канал для результатов
	processedJobs := make([]int, 0, numJobs) // Слайс для хранения результатов

	// Запускаем ограниченное количество горутин-обработчиков
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for job := range jobs {
				// Имитируем обработку задания
				fmt.Printf("Воркер %d обрабатывает задание %d\n", id, job)
				time.Sleep(time.Duration(50*id) * time.Millisecond)

				// Отправляем результат (квадрат номера задания)
				results <- job * job
				fmt.Printf("Воркер %d завершил задание %d\n", id, job)
			}
		}(w)
	}

	// Отправляем задания в канал
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs) // Сигнализируем, что больше заданий не будет

	// Ожидаем завершения всех воркеров и закрываем канал результатов
	go func() {
		wg.Wait()
		close(results)
	}()

	// Собираем результаты
	for result := range results {
		fmt.Printf("Получен результат: %d\n", result)
		processedJobs = append(processedJobs, result)
	}

	fmt.Printf("Все задания обработаны, получено %d результатов\n", len(processedJobs))
	return processedJobs
}
