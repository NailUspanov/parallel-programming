package basic_goroutine

import (
	"bytes"
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
