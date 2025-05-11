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

// RunPingPong демонстрирует пинг-понг между двумя горутинами
// с использованием небуферизованного канала
func RunPingPong(rounds int) []string {
	pingChan := make(chan string)
	pongChan := make(chan string)
	resultChan := make(chan string, rounds*2) // Канал для хранения результатов
	done := make(chan bool)
	results := make([]string, 0, rounds*2)

	// Горутина для игрока "Ping"
	go func() {
		for i := 0; i < rounds; i++ {
			// Отправляем "Ping"
			msg := fmt.Sprintf("Ping %d", i+1)
			pingChan <- msg
			resultChan <- msg

			// Получаем ответ от "Pong"
			response := <-pongChan
			fmt.Println("Игрок Ping получил:", response)
		}
		done <- true
	}()

	// Горутина для игрока "Pong"
	go func() {
		for i := 0; i < rounds; i++ {
			// Получаем "Ping"
			msg := <-pingChan
			fmt.Println("Игрок Pong получил:", msg)

			// Отправляем "Pong"
			response := fmt.Sprintf("Pong %d", i+1)
			pongChan <- response
			resultChan <- response
		}
	}()

	// Ожидаем завершения игры
	<-done
	close(resultChan)

	// Собираем результаты
	for msg := range resultChan {
		results = append(results, msg)
	}

	fmt.Println("Игра пинг-понг завершена")
	return results
}
