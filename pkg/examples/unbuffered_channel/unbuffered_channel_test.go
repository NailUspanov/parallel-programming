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

	// Запускаем функцию
	value := RunUnbufferedChannel()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые данные
	if !strings.Contains(output, "42") {
		t.Error("Ожидаемое значение 42 не найдено в выводе")
	}

	// Проверяем возвращаемое значение
	if value != 42 {
		t.Errorf("Ожидаемое возвращаемое значение 42, получено %d", value)
	}
}

func TestRunPingPong(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с 3 раундами
	results := RunPingPong(3)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые строки
	expectedPhrases := []string{
		"Игрок Pong получил: Ping",
		"Игрок Ping получил: Pong",
		"Игра пинг-понг завершена",
	}

	for _, phrase := range expectedPhrases {
		if !strings.Contains(output, phrase) {
			t.Errorf("Ожидаемая фраза не найдена: %s", phrase)
		}
	}

	// Проверяем количество пинг-понгов
	pingCount := 0
	pongCount := 0

	for _, result := range results {
		if strings.HasPrefix(result, "Ping") {
			pingCount++
		} else if strings.HasPrefix(result, "Pong") {
			pongCount++
		}
	}

	if pingCount != 3 {
		t.Errorf("Ожидалось 3 пинга, получено %d", pingCount)
	}

	if pongCount != 3 {
		t.Errorf("Ожидалось 3 понга, получено %d", pongCount)
	}
}
