package race_condition

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunRaceConditionProblem(t *testing.T) {
	// Мы не можем точно протестировать гонку данных, но можем запустить пример
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	result := RunRaceConditionProblem()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемую строку
	if !strings.Contains(output, "Итоговое значение (с гонкой данных):") {
		t.Error("Ожидаемый вывод не найден")
	}

	// Мы не можем гарантировать точное значение из-за гонки данных,
	// но можем проверить, что оно в каком-то диапазоне
	if result < 0 || result > 1000 {
		t.Errorf("Ожидаемое значение в диапазоне [0, 1000], получено %d", result)
	}
}

func TestRunRaceConditionSolution(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	result := RunRaceConditionSolution()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемую строку
	if !strings.Contains(output, "Итоговое значение (решение с каналами):") {
		t.Error("Ожидаемый вывод не найден")
	}

	// Проверяем, что результат всегда 1000
	if result != 1000 {
		t.Errorf("Ожидаемое значение 1000, получено %d", result)
	}
}

func TestRunRaceConditionMutex(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию и получаем результат
	result := RunRaceConditionMutex()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемую строку
	if !strings.Contains(output, "Итоговое значение (решение с мьютексом):") {
		t.Error("Ожидаемый вывод не найден")
	}

	// Проверяем, что результат всегда 1000
	if result != 1000 {
		t.Errorf("Ожидаемое значение 1000, получено %d", result)
	}
}
