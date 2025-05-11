package basic_goroutine

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestRunBasicGoroutine(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	RunBasicGoroutine()

	// Даем время для вывода
	time.Sleep(200 * time.Millisecond)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые строки
	if !strings.Contains(output, "Это выполняется в основной функции") {
		t.Error("Ожидаемый вывод из основной функции не найден")
	}

	if !strings.Contains(output, "Это выполняется в горутине") {
		t.Error("Ожидаемый вывод из горутины не найден")
	}
}

func TestRunMultipleGoroutines(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с 5 горутинами
	results := RunMultipleGoroutines(5)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые строки
	for i := 0; i < 5; i++ {
		expectedOutput := fmt.Sprintf("Горутина %d завершила работу", i)
		if !strings.Contains(output, expectedOutput) {
			t.Errorf("Ожидаемый вывод '%s' не найден", expectedOutput)
		}
	}

	if !strings.Contains(output, "Получены результаты от 5 горутин") {
		t.Error("Ожидаемый итоговый вывод не найден")
	}

	// Проверяем, что возвращены результаты от всех горутин
	if len(results) != 5 {
		t.Errorf("Ожидаемая длина слайса результатов 5, получено %d", len(results))
	}

	// Проверяем, что результаты содержат квадраты чисел от 0 до 4
	expectedResults := map[int]bool{0: false, 1: false, 4: false, 9: false, 16: false}
	for _, r := range results {
		expectedResults[r] = true
	}

	for val, found := range expectedResults {
		if !found {
			t.Errorf("Ожидаемый результат %d не найден в возвращаемых значениях", val)
		}
	}
}
