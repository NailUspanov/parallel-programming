package buffered_channel

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"
)

func TestRunBufferedChannel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	results := RunBufferedChannel()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемую строку
	expectedOutput := "Полученные значения: 1, 2"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Ожидаемая строка '%s' не найдена в выводе", expectedOutput)
	}

	// Проверяем возвращаемые значения
	if len(results) != 2 || results[0] != 1 || results[1] != 2 {
		t.Errorf("Ожидаемые значения [1 2], получено %v", results)
	}
}

func TestRunWorkerPool(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с 3 воркерами и 5 заданиями
	results := RunWorkerPool(3, 5)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем вывод на наличие сообщений о работе воркеров
	for w := 1; w <= 3; w++ {
		workerPrefix := fmt.Sprintf("Воркер %d", w)
		if !strings.Contains(output, workerPrefix) {
			t.Errorf("Ожидаемый вывод содержащий '%s' не найден", workerPrefix)
		}
	}

	// Проверяем наличие сообщения о завершении всех заданий
	if !strings.Contains(output, "Все задания обработаны, получено 5 результатов") {
		t.Error("Ожидаемое сообщение о завершении не найдено")
	}

	// Проверяем количество результатов
	if len(results) != 5 {
		t.Errorf("Ожидаемое количество результатов 5, получено %d", len(results))
	}

	// Проверяем, что результаты содержат квадраты чисел от 1 до 5
	expectedResults := []int{1, 4, 9, 16, 25}

	// Сортируем результаты для правильного сравнения
	sortedResults := make([]int, len(results))
	copy(sortedResults, results)
	sort.Ints(sortedResults)

	for i, expected := range expectedResults {
		if i >= len(sortedResults) || sortedResults[i] != expected {
			t.Errorf("Ожидаемый результат %d не найден или не на правильной позиции", expected)
		}
	}
}
