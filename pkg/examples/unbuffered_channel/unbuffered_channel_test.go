package unbuffered_channel

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunUnbufferedChannel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	value := RunUnbufferedChannel()

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
