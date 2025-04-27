package directed_channels

import (
	"fmt"
)

// sender передает данные в канал только для отправки
func sender(ch chan<- int) {
	// Канал только для отправки
	ch <- 42
}

// receiver получает данные из канала только для получения
func receiver(ch <-chan int) int {
	// Канал только для получения
	value := <-ch
	fmt.Println(value)
	return value
}

// RunDirectedChannels демонстрирует использование направленных каналов
func RunDirectedChannels() int {
	ch := make(chan int)

	go sender(ch)
	return receiver(ch)
}
