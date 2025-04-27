package unbuffered_channel

import "fmt"

// RunUnbufferedChannel демонстрирует использование небуферизованного канала
func RunUnbufferedChannel() int {
	ch := make(chan int) // Небуферизованный канал

	go func() {
		// Отправка значения в канал блокирует горутину
		// до получения этого значения
		ch <- 42
	}()

	// Получение значения из канала блокирует главную горутину
	// до отправки значения
	value := <-ch
	fmt.Println(value) // Выведет: 42

	return value
}
