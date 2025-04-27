package buffered_channel

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunBufferedChannel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	values := RunBufferedChannel()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые значения
	expectedOutput := "Полученные значения: 1, 2"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Ожидаемый вывод '%s' не найден в '%s'", expectedOutput, output)
	}

	// Проверяем возвращаемые значения
	if len(values) != 2 || values[0] != 1 || values[1] != 2 {
		t.Errorf("Ожидаемые возвращаемые значения [1, 2], получено %v", values)
	}
}
