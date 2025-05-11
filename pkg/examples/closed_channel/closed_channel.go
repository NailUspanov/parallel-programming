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

// CheckChannelState демонстрирует различные способы проверки состояния канала
func CheckChannelState() map[string]bool {
	// Создаем канал
	ch := make(chan int, 1)
	results := make(map[string]bool)

	// 1. Проверка канала с помощью дополнительного флага ok
	ch <- 42
	value, isOpen := <-ch
	fmt.Printf("Значение: %d, канал открыт: %t\n", value, isOpen)
	results["Проверка канала с данными"] = isOpen

	// 2. Проверка закрытого канала
	close(ch)
	value, isOpen = <-ch
	fmt.Printf("Значение после закрытия: %d, канал открыт: %t\n", value, isOpen)
	results["Проверка закрытого канала"] = isOpen

	// 3. Безопасная отправка в канал с проверкой закрытия
	closedCh := make(chan int)
	close(closedCh)

	// Функция для безопасной отправки данных
	safeSend := func(ch chan int, value int) (success bool) {
		defer func() {
			// Перехватываем панику при отправке в закрытый канал
			if recover() != nil {
				fmt.Println("Невозможно отправить в закрытый канал!")
				success = false
			}
		}()

		ch <- value // Вызовет панику, если канал закрыт
		return true
	}

	// Проверяем отправку в закрытый канал
	sendResult := safeSend(closedCh, 100)
	fmt.Printf("Результат отправки в закрытый канал: %t\n", sendResult)
	results["Отправка в закрытый канал"] = sendResult

	// 4. Использование select для неблокирующей проверки
	selectCh := make(chan int, 1)
	selectCh <- 1

	// Неблокирующее чтение из канала
	select {
	case val, ok := <-selectCh:
		fmt.Printf("Прочитано значение %d из канала, канал открыт: %t\n", val, ok)
		results["Select неблокирующее чтение"] = ok
	default:
		fmt.Println("Канал пуст или закрыт")
		results["Select неблокирующее чтение"] = false
	}

	// Закрываем канал и проверяем снова
	close(selectCh)

	select {
	case val, ok := <-selectCh:
		fmt.Printf("Прочитано значение %d из закрытого канала, канал открыт: %t\n", val, ok)
		results["Select чтение из закрытого канала"] = ok
	default:
		fmt.Println("Канал пуст или закрыт")
		results["Select чтение из закрытого канала"] = false
	}

	fmt.Println("Проверки состояния канала завершены")
	return results
}
