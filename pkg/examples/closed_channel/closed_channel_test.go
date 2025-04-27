package closed_channel

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunClosedChannel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	values := RunClosedChannel()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые значения
	for i := 0; i < 5; i++ {
		if !strings.Contains(output, string(rune('0'+i))) {
			t.Errorf("Ожидаемое значение %d не найдено в выводе", i)
		}
	}

	// Проверяем сообщение о закрытии канала
	if !strings.Contains(output, "Канал закрыт, цикл завершен") {
		t.Error("Ожидаемое сообщение о закрытии канала не найдено в выводе")
	}

	// Проверяем возвращаемые значения
	if len(values) != 5 {
		t.Errorf("Ожидаемая длина слайса результатов 5, получено %d", len(values))
	}

	for i := 0; i < 5; i++ {
		if values[i] != i {
			t.Errorf("Ожидаемое значение values[%d] = %d, получено %d", i, i, values[i])
		}
	}
}
