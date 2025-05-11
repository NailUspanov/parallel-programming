package directed_channels

import (
	"bytes"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
)

func TestRunDirectedChannels(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	value := RunDirectedChannels()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемое значение
	if !strings.Contains(output, "42") {
		t.Error("Ожидаемое значение 42 не найдено в выводе")
	}

	// Проверяем возвращаемое значение
	if value != 42 {
		t.Errorf("Ожидаемое возвращаемое значение 42, получено %d", value)
	}
}

func TestRunDataPipeline(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Входные данные
	numbers := []int{1, 2, 3, 4, 5}

	// Запускаем функцию и получаем результаты
	squaredValues, sumResult := RunDataPipeline(numbers)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем промежуточные результаты в выводе
	for _, n := range numbers {
		squareStr := "Промежуточный результат (квадрат): " + strconv.Itoa(n*n)
		if !strings.Contains(output, squareStr) {
			t.Errorf("Ожидаемый вывод '%s' не найден", squareStr)
		}
	}

	// Проверяем итоговый результат в выводе
	expectedSum := 55 // Сумма квадратов 1^2 + 2^2 + 3^2 + 4^2 + 5^2 = 1 + 4 + 9 + 16 + 25 = 55
	if !strings.Contains(output, "Итоговая сумма квадратов: "+strconv.Itoa(expectedSum)) {
		t.Errorf("Ожидаемый вывод 'Итоговая сумма квадратов: %d' не найден", expectedSum)
	}

	// Проверяем возвращенные промежуточные результаты
	if len(squaredValues) != len(numbers) {
		t.Errorf("Ожидаемое количество промежуточных результатов %d, получено %d", len(numbers), len(squaredValues))
	}

	// Сортируем промежуточные результаты для корректного сравнения
	sort.Ints(squaredValues)
	expectedSquares := []int{1, 4, 9, 16, 25}

	for i, expected := range expectedSquares {
		if i >= len(squaredValues) || squaredValues[i] != expected {
			t.Errorf("Ожидаемый промежуточный результат %d на позиции %d, получено %d",
				expected, i, squaredValues[i])
		}
	}

	// Проверяем итоговую сумму
	if sumResult != expectedSum {
		t.Errorf("Ожидаемая сумма квадратов %d, получено %d", expectedSum, sumResult)
	}
}

func TestSenderReceiver(t *testing.T) {
	// Проверяем работу функций sender и receiver отдельно
	ch := make(chan int)

	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем sender в горутине
	go sender(ch)

	// Запускаем receiver и получаем результат
	value := receiver(ch)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемое значение
	if !strings.Contains(output, "42") {
		t.Error("Ожидаемое значение 42 не найдено в выводе")
	}

	// Проверяем возвращаемое значение
	if value != 42 {
		t.Errorf("Ожидаемое возвращаемое значение 42, получено %d", value)
	}
}

func TestRunGeneratorExample(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Входные данные
	numbers := []int{1, 2, 3, 4, 5}

	// Запускаем функцию и получаем результаты
	results := RunGeneratorExample(numbers)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем результаты в выводе
	for _, n := range numbers {
		expectedStr := "Сгенерировано значение: " + strconv.Itoa(n)
		if !strings.Contains(output, expectedStr) {
			t.Errorf("Ожидаемый вывод '%s' не найден", expectedStr)
		}
	}

	// Проверяем возвращенные результаты
	if len(results) != len(numbers) {
		t.Errorf("Ожидаемое количество результатов %d, получено %d", len(numbers), len(results))
	}

	// Сортируем результаты для корректного сравнения
	sort.Ints(results)
	for i, expected := range numbers {
		if i >= len(results) || results[i] != expected {
			t.Errorf("Ожидаемый результат %d на позиции %d, получено %d",
				expected, i, results[i])
		}
	}
}

func TestRunTransformerExample(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Входные данные
	numbers := []int{1, 2, 3, 4, 5}

	// Запускаем функцию и получаем результаты
	results := RunTransformerExample(numbers)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем преобразованные результаты в выводе
	for _, n := range numbers {
		expectedStr := "Преобразованное значение: " + strconv.Itoa(n*n)
		if !strings.Contains(output, expectedStr) {
			t.Errorf("Ожидаемый вывод '%s' не найден", expectedStr)
		}
	}

	// Проверяем возвращенные результаты
	if len(results) != len(numbers) {
		t.Errorf("Ожидаемое количество результатов %d, получено %d", len(numbers), len(results))
	}

	// Сортируем результаты для корректного сравнения
	sort.Ints(results)
	expectedSquares := []int{1, 4, 9, 16, 25}
	for i, expected := range expectedSquares {
		if i >= len(results) || results[i] != expected {
			t.Errorf("Ожидаемый результат %d на позиции %d, получено %d",
				expected, i, results[i])
		}
	}
}

func TestRunSinkExample(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Входные данные
	numbers := []int{1, 2, 3, 4, 5}

	// Запускаем функцию и получаем результат
	result := RunSinkExample(numbers)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем результат в выводе
	expectedSum := 55 // Сумма квадратов 1^2 + 2^2 + 3^2 + 4^2 + 5^2 = 1 + 4 + 9 + 16 + 25 = 55
	expectedStr := "Результат агрегации: " + strconv.Itoa(expectedSum)
	if !strings.Contains(output, expectedStr) {
		t.Errorf("Ожидаемый вывод '%s' не найден", expectedStr)
	}

	// Проверяем возвращаемый результат
	if result != expectedSum {
		t.Errorf("Ожидаемый результат %d, получено %d", expectedSum, result)
	}
}
