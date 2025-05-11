package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/basic_goroutine"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/buffered_channel"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/closed_channel"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/combined"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/directed_channels"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/race_condition"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/unbuffered_channel"
	"github.com/nailuspanov/sstu-projects/parallel-programming/pkg/examples/waitgroups"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run main.go [имя_примера]")
		fmt.Println("Доступные примеры:")
		fmt.Println("  all - запустить все примеры")
		fmt.Println("  goroutine - базовый пример горутин")
		fmt.Println("  unbuffered - небуферизованный канал")
		fmt.Println("  buffered - буферизованный канал")
		fmt.Println("  directed - направленные каналы")
		fmt.Println("  closed - закрытие канала")
		fmt.Println("  waitgroups - группы ожидания")
		fmt.Println("  race - пример гонки данных и решения")
		fmt.Println("  combined - комбинированное использование каналов и групп ожидания")
		return
	}

	example := os.Args[1]

	fmt.Println("=== Запуск примера:", example, "===")
	fmt.Println()

	switch example {
	case "all":
		runAll()
	case "goroutine":
		basic_goroutine.RunBasicGoroutine()
	case "unbuffered":
		result := unbuffered_channel.RunUnbufferedChannel()
		fmt.Println("Результат:", result)
	case "buffered":
		results := buffered_channel.RunBufferedChannel()
		fmt.Println("Результаты:", results)
	case "directed":
		result := directed_channels.RunDirectedChannels()
		fmt.Println("Результат:", result)
	case "closed":
		results := closed_channel.RunClosedChannel()
		fmt.Println("Результаты:", results)
	case "waitgroups":
		results := waitgroups.RunWaitGroups()
		fmt.Println("Результаты:", results)
	case "race":
		fmt.Println("1. Демонстрация проблемы гонки данных:")
		race_condition.RunRaceConditionProblem(500)

		fmt.Println("\n2. Решения проблемы гонки данных:")
		race_condition.RunRaceConditionSolution(500)
	case "combined":
		results := combined.RunBasicWorkerPool(4, 10)
		fmt.Println("Результаты:", results)

		fmt.Println("\nПример с таймаутом:")
		resultsTimeout := combined.RunWorkerPoolWithTimeout(4, 20, 500*time.Millisecond)
		fmt.Println("Результаты с таймаутом:", resultsTimeout)
	default:
		fmt.Println("Неизвестный пример:", example)
	}
}

func runAll() {
	examples := []struct {
		name string
		fn   func()
	}{
		{"Базовый пример горутин", func() { basic_goroutine.RunBasicGoroutine() }},
		{"Небуферизованный канал", func() { unbuffered_channel.RunUnbufferedChannel() }},
		{"Буферизованный канал", func() { buffered_channel.RunBufferedChannel() }},
		{"Направленные каналы", func() { directed_channels.RunDirectedChannels() }},
		{"Закрытие канала", func() { closed_channel.RunClosedChannel() }},
		{"Группы ожидания", func() { waitgroups.RunWaitGroups() }},
		{"Проблема гонки данных", func() { race_condition.RunRaceConditionProblem(500) }},
		{"Решение гонки данных", func() { race_condition.RunRaceConditionSolution(500) }},
		{"Комбинированный пример", func() { combined.RunBasicWorkerPool(4, 10) }},
		{"Пример с контекстом и таймаутом", func() {
			combined.RunWorkerPoolWithTimeout(4, 20, 500*time.Millisecond)
		}},
	}

	var wg sync.WaitGroup

	for _, e := range examples {
		wg.Add(1)
		go func(name string, fn func()) {
			defer wg.Done()

			fmt.Println("\n=== Запуск примера:", name, "===")
			fn()
			fmt.Println("=== Пример", name, "завершен ===")

			// Пауза между примерами для лучшей читаемости вывода
			time.Sleep(100 * time.Millisecond)
		}(e.name, e.fn)

		// Небольшая пауза между запусками для лучшей читаемости вывода
		time.Sleep(500 * time.Millisecond)
	}

	wg.Wait()
	fmt.Println("\nВсе примеры выполнены")
}
