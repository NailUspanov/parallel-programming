package waitgroups

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunWaitGroups(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	results := RunWaitGroups()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые строки
	expectedOutputs := []string{
		"Горутина 1",
		"Горутина 2",
		"Все горутины завершены",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Ожидаемый вывод '%s' не найден", expected)
		}
	}

	// Проверяем возвращаемые значения
	if len(results) != 2 {
		t.Errorf("Ожидаемая длина слайса результатов 2, получено %d", len(results))
	}

	// Проверяем, что результаты включают оба ожидаемых значения
	hasGoroutine1 := false
	hasGoroutine2 := false

	for _, result := range results {
		if result == "Горутина 1" {
			hasGoroutine1 = true
		}
		if result == "Горутина 2" {
			hasGoroutine2 = true
		}
	}

	if !hasGoroutine1 {
		t.Error("Результат 'Горутина 1' не найден в возвращаемых значениях")
	}

	if !hasGoroutine2 {
		t.Error("Результат 'Горутина 2' не найден в возвращаемых значениях")
	}
}
