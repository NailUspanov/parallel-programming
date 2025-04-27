package buffered_channel

import "fmt"

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
