package closed_channel

import "fmt"

// RunClosedChannel демонстрирует закрытие и перебор канала
func RunClosedChannel() []int {
	ch := make(chan int)
	results := make([]int, 0, 5)

	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
		}
		close(ch) // Закрытие канала
	}()

	// Перебор значений канала до его закрытия
	for value := range ch {
		fmt.Println(value)
		results = append(results, value)
	}
	// После закрытия канала цикл range завершится

	fmt.Println("Канал закрыт, цикл завершен")
	return results
}
