package combined

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestRunCombined(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	results := RunCombined()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие сообщений от воркеров
	for i := 1; i <= 3; i++ { // 3 воркера
		workerStarted := false
		workerFinished := false

		for j := 1; j <= 10; j++ { // 10 заданий
			if strings.Contains(output,
				fmt.Sprintf("Воркер %d начал обработку задания %d", i, j)) {
				workerStarted = true
			}
			if strings.Contains(output,
				fmt.Sprintf("Воркер %d завершил обработку задания %d", i, j)) {
				workerFinished = true
			}
		}

		if !workerStarted || !workerFinished {
			t.Errorf("Не все сообщения от воркера %d найдены", i)
		}
	}

	// Проверяем наличие результатов
	for i := 1; i <= 10; i++ {
		expectedResult := i * 2
		resultFound := false

		for _, result := range results {
			if result == expectedResult {
				resultFound = true
				break
			}
		}

		if !resultFound {
			t.Errorf("Ожидаемый результат %d не найден", expectedResult)
		}
	}

	// Проверяем количество результатов
	if len(results) != 10 {
		t.Errorf("Ожидаемое количество результатов 10, получено %d", len(results))
	}
}
