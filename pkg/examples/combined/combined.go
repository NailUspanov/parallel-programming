package combined

import (
	"context"
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

// RunBasicWorkerPool демонстрирует базовый рабочий пул с использованием каналов и WaitGroups
func RunBasicWorkerPool(numWorkers, numJobs int) []int {
	var wg sync.WaitGroup
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	collectedResults := make([]int, 0, numJobs)

	// Запускаем воркеров
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	// Отправляем задачи
	for j := 1; j <= numJobs; j++ {
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

// contextWorker обрабатывает задания с поддержкой отмены через контекст
func contextWorker(ctx context.Context, id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// Контекст был отменен, завершаем работу
			fmt.Printf("Воркер %d завершается из-за отмены\n", id)
			return
		case j, ok := <-jobs:
			if !ok {
				// Канал заданий закрыт, завершаем работу
				fmt.Printf("Воркер %d завершается, заданий больше нет\n", id)
				return
			}

			// Имитация работы с проверкой отмены
			fmt.Printf("Воркер %d начал обработку задания %d\n", id, j)

			// Разбиваем обработку на части, чтобы можно было отменить
			// в процессе выполнения
			for i := 0; i < 5; i++ {
				select {
				case <-ctx.Done():
					fmt.Printf("Воркер %d: задание %d отменено\n", id, j)
					return
				default:
					time.Sleep(50 * time.Millisecond)
				}
			}

			// Проверяем, не был ли контекст отменен во время обработки
			select {
			case <-ctx.Done():
				fmt.Printf("Воркер %d: задание %d отменено\n", id, j)
				return
			default:
				// Задание выполнено успешно
				fmt.Printf("Воркер %d завершил обработку задания %d\n", id, j)
				select {
				case results <- j * 10: // Отправляем результат
				case <-ctx.Done(): // Но может быть отменено и здесь
					return
				}
			}
		}
	}
}

// RunWorkerPoolWithTimeout демонстрирует рабочий пул с таймаутом через контекст
func RunWorkerPoolWithTimeout(numWorkers, numJobs int, timeout time.Duration) []int {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel() // Гарантируем вызов отмены при выходе из функции

	return runContextPool(ctx, numWorkers, numJobs)
}

// RunWorkerPoolWithCancel демонстрирует рабочий пул с отменой через контекст
func RunWorkerPoolWithCancel(numWorkers, numJobs int, cancelAfter time.Duration) []int {
	// Создаем контекст с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())

	// Запускаем горутину, которая отменит контекст через указанное время
	go func() {
		time.Sleep(cancelAfter)
		fmt.Println("Время истекло, отменяем контекст")
		cancel()
	}()

	defer cancel() // Гарантируем вызов отмены при выходе из функции
	return runContextPool(ctx, numWorkers, numJobs)
}

// runContextPool содержит общую логику для запуска пула с контекстом
func runContextPool(ctx context.Context, numWorkers, numJobs int) []int {
	var wg sync.WaitGroup
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)
	collectedResults := make([]int, 0, numJobs)

	// Запускаем воркеров с поддержкой контекста
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go contextWorker(ctx, w, jobs, results, &wg)
	}

	// Запускаем горутину для отправки заданий
	go func() {
		for j := 1; j <= numJobs; j++ {
			select {
			case <-ctx.Done():
				// Контекст отменен, прекращаем отправку заданий
				fmt.Println("Отмена отправки заданий")
				return
			case jobs <- j:
				fmt.Printf("Отправлено задание %d\n", j)
				time.Sleep(100 * time.Millisecond) // Имитация задержки между заданиями
			}
		}
		close(jobs) // Закрываем канал, когда все задания отправлены
		fmt.Println("Все задания отправлены")
	}()

	// Запускаем горутину для сбора результатов
	go func() {
		// Ожидаем завершения всех воркеров
		wg.Wait()
		// Закрываем канал результатов
		close(results)
		fmt.Println("Все воркеры завершены")
	}()

	// Собираем результаты до отмены контекста или до закрытия канала
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Контекст отменен:", ctx.Err())
			return collectedResults // Возвращаем собранные результаты
		case r, ok := <-results:
			if !ok {
				// Канал результатов закрыт, все воркеры завершились
				fmt.Println("Сбор результатов завершен")
				return collectedResults
			}
			fmt.Printf("Получен результат: %d\n", r)
			collectedResults = append(collectedResults, r)
		}
	}
}

// RunCombined демонстрирует комбинированное использование каналов и WaitGroups
// Оставлен для обратной совместимости
func RunCombined() []int {
	return RunBasicWorkerPool(3, 10)
}

// RunContextCancellation демонстрирует использование контекста для отмены горутин
// Оставлен для обратной совместимости
func RunContextCancellation(numWorkers, numJobs int, cancelAfter time.Duration) []int {
	return RunWorkerPoolWithTimeout(numWorkers, numJobs, cancelAfter)
}
