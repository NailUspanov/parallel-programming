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

// generator создает канал и заполняет его числами
func generator(numbers ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range numbers {
			out <- n
		}
	}()
	return out
}

// square получает числа из входного канала, возводит в квадрат и отправляет в выходной
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			out <- n * n
		}
	}()
	return out
}

// sum складывает все числа из входного канала и отправляет результат в выходной
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		sum := 0
		for n := range in {
			sum += n
		}
		out <- sum
	}()
	return out
}

// RunDataPipeline демонстрирует использование пайплайна обработки данных
// с направленными каналами
func RunDataPipeline(numbers []int) ([]int, int) {
	// Создаем пайплайн обработки
	genCh := generator(numbers...) // Этап 1: Генерация чисел

	// Создаем промежуточный канал для перехвата квадратов
	// Используем буферизованный канал для хранения всех значений
	squareResults := make(chan int, len(numbers))

	// Создаем тройник: направляем данные от генератора в два канала
	splitCh1 := make(chan int)
	splitCh2 := make(chan int)
	go func() {
		defer close(splitCh1)
		defer close(splitCh2)
		for n := range genCh {
			splitCh1 <- n
			splitCh2 <- n
		}
	}()

	// Первый канал квадратов - для промежуточных результатов
	squareCh1 := square(splitCh1)

	// Второй канал квадратов - для суммирования
	squareCh2 := square(splitCh2)
	sumCh := sum(squareCh2)

	// Создаем слайс для хранения промежуточных результатов
	squaredValues := make([]int, 0, len(numbers))

	// Запускаем сбор и вывод промежуточных результатов
	go func() {
		for val := range squareCh1 {
			fmt.Printf("Промежуточный результат (квадрат): %d\n", val)
			squareResults <- val
		}
		close(squareResults)
	}()

	// Получаем итоговый результат
	result := <-sumCh
	fmt.Printf("Итоговая сумма квадратов: %d\n", result)

	// Собираем все промежуточные результаты
	for val := range squareResults {
		squaredValues = append(squaredValues, val)
	}

	return squaredValues, result
}

// RunGeneratorExample демонстрирует использование генератора (источника данных)
func RunGeneratorExample(numbers []int) []int {
	results := make([]int, 0, len(numbers))
	gen := generator(numbers...)

	for val := range gen {
		fmt.Printf("Сгенерировано значение: %d\n", val)
		results = append(results, val)
	}

	return results
}

// RunTransformerExample демонстрирует использование трансформера (преобразователя данных)
func RunTransformerExample(numbers []int) []int {
	results := make([]int, 0, len(numbers))
	gen := generator(numbers...)
	transformer := square(gen)

	for val := range transformer {
		fmt.Printf("Преобразованное значение: %d\n", val)
		results = append(results, val)
	}

	return results
}

// RunSinkExample демонстрирует использование приемника (агрегатора данных)
func RunSinkExample(numbers []int) int {
	gen := generator(numbers...)
	squareCh := square(gen)
	sumCh := sum(squareCh)

	result := <-sumCh
	fmt.Printf("Результат агрегации: %d\n", result)

	return result
}
