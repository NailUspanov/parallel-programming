package closed_channel

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRunClosedChannel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	results := RunClosedChannel()

	// Возвращаем stdout и читаем перехваченный вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие цифр от 0 до 4 в выводе
	for i := 0; i < 5; i++ {
		if !strings.Contains(output, fmt.Sprintf("%d", i)) {
			t.Errorf("Ожидаемое значение %d не найдено в выводе", i)
		}
	}

	// Проверяем, что было сообщение о закрытии канала
	if !strings.Contains(output, "Канал закрыт, цикл завершен") {
		t.Error("Ожидаемое сообщение о закрытии канала не найдено")
	}

	// Проверяем возвращаемые значения
	if len(results) != 5 {
		t.Errorf("Ожидаемая длина слайса результатов 5, получено %d", len(results))
	}

	for i, val := range results {
		if val != i {
			t.Errorf("Ожидаемое значение %d на позиции %d, получено %d", i, i, val)
		}
	}
}

func TestCheckChannelState(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результаты
	results := CheckChannelState()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем ожидаемые сообщения в выводе
	expectedOutputs := []string{
		"Значение: 42, канал открыт: true",
		"Значение после закрытия: 0, канал открыт: false",
		"Невозможно отправить в закрытый канал!",
		"Результат отправки в закрытый канал: false",
		"Прочитано значение 1 из канала, канал открыт: true",
		"Прочитано значение 0 из закрытого канала, канал открыт: false",
		"Проверки состояния канала завершены",
	}

	for _, expected := range expectedOutputs {
		if !strings.Contains(output, expected) {
			t.Errorf("Ожидаемое сообщение '%s' не найдено в выводе", expected)
		}
	}

	// Проверяем возвращаемые результаты
	expectedResults := map[string]bool{
		"Проверка канала с данными":         true,
		"Проверка закрытого канала":         false,
		"Отправка в закрытый канал":         false,
		"Select неблокирующее чтение":       true,
		"Select чтение из закрытого канала": false,
	}

	if len(results) != len(expectedResults) {
		t.Errorf("Ожидаемое количество результатов %d, получено %d", len(expectedResults), len(results))
	}

	for key, expected := range expectedResults {
		if actual, exists := results[key]; !exists || actual != expected {
			if !exists {
				t.Errorf("Ожидаемый результат для '%s' не найден", key)
			} else {
				t.Errorf("Для '%s' ожидается %t, получено %t", key, expected, actual)
			}
		}
	}
}
