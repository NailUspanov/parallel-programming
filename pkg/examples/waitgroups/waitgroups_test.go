package waitgroups

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
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

	// Проверяем корректность завершения вывода
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		t.Error("Вывод должен содержать минимум 3 строки")
	} else if lines[len(lines)-1] != "Все горутины завершены" {
		t.Error("Последняя строка вывода должна быть 'Все горутины завершены'")
	}
}

func TestRunDynamicWaitGroups(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с 2 начальными и 1 дополнительной горутиной
	results := RunDynamicWaitGroups(2, 1)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, содержит ли вывод ожидаемые строки
	if !strings.Contains(output, "Добавляем дополнительные горутины") {
		t.Error("Ожидаемое сообщение о добавлении горутин не найдено")
	}

	expectedTotal := fmt.Sprintf("Все горутины завершены, всего: %d", 2+1)
	if !strings.Contains(output, expectedTotal) {
		t.Error("Ожидаемое итоговое сообщение не найдено")
	}

	// Проверяем, что возвращены результаты от всех горутин
	if len(results) != 3 {
		t.Errorf("Ожидаемая длина слайса результатов 3, получено %d", len(results))
	}

	// Проверяем, что все результаты содержат ожидаемые горутины (от 0 до 2)
	goroutineResults := make(map[string]bool)
	for i := 0; i < 3; i++ {
		goroutineResults[fmt.Sprintf("Горутина %d", i)] = false
	}

	for _, result := range results {
		goroutineResults[result] = true
	}

	for goroutine, found := range goroutineResults {
		if !found {
			t.Errorf("Ожидаемый результат '%s' не найден в возвращаемых значениях", goroutine)
		}
	}

	// Проверка порядка выполнения
	reTimeOrder := regexp.MustCompile(`Горутина (\d+)`)
	matches := reTimeOrder.FindAllStringSubmatch(output, -1)
	var order []int

	for _, match := range matches {
		if len(match) > 1 {
			if id, err := strconv.Atoi(match[1]); err == nil {
				order = append(order, id)
			}
		}
	}

	// Проверка, что хотя бы начальные горутины запустились раньше дополнительных
	initialFound := false
	additionalFound := false

	for _, id := range order {
		if id < 2 { // Начальные горутины
			initialFound = true
		} else { // Дополнительные горутины
			if initialFound && !additionalFound {
				additionalFound = true
			}
		}
	}

	if !initialFound || !additionalFound {
		t.Log("Порядок выполнения горутин может быть произвольным, но обычно начальные горутины запускаются раньше")
	}
}

func TestRunDynamicWaitGroupsVariations(t *testing.T) {
	testCases := []struct {
		name       string
		initial    int
		additional int
		expected   int
	}{
		{"Без дополнительных горутин", 3, 0, 3},
		{"Равное число", 2, 2, 4},
		{"Больше дополнительных", 1, 4, 5},
		{"Много горутин", 5, 5, 10},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Перенаправляем stdout для перехвата вывода
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Засекаем время выполнения, чтобы проверить, что функция не зависает
			start := time.Now()

			// Запускаем функцию с указанными параметрами
			results := RunDynamicWaitGroups(tc.initial, tc.additional)

			// Проверяем, что функция выполнилась за разумное время
			duration := time.Since(start)
			if duration > 5*time.Second {
				t.Errorf("Функция выполнялась слишком долго: %v", duration)
			}

			// Восстанавливаем stdout и получаем вывод
			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			// Проверка итогового сообщения
			expectedTotal := fmt.Sprintf("Все горутины завершены, всего: %d", tc.expected)
			if !strings.Contains(output, expectedTotal) {
				t.Errorf("Ожидаемое итоговое сообщение '%s' не найдено", expectedTotal)
			}

			// Проверка количества результатов
			if len(results) != tc.expected {
				t.Errorf("Ожидаемая длина слайса результатов %d, получено %d", tc.expected, len(results))
			}

			// Проверка всех ожидаемых горутин
			expectedGoroutines := make(map[string]bool)
			for i := 0; i < tc.initial+tc.additional; i++ {
				expectedGoroutines[fmt.Sprintf("Горутина %d", i)] = false
			}

			for _, result := range results {
				if _, ok := expectedGoroutines[result]; ok {
					expectedGoroutines[result] = true
				} else {
					t.Errorf("Неожиданный результат: %s", result)
				}
			}

			// Проверка, что все ожидаемые горутины отработали
			missingGoroutines := 0
			for goroutine, found := range expectedGoroutines {
				if !found {
					missingGoroutines++
					t.Logf("Ожидаемый результат '%s' не найден", goroutine)
				}
			}

			if missingGoroutines > 0 {
				t.Errorf("Не найдено %d ожидаемых горутин из %d", missingGoroutines, tc.expected)
			}
		})
	}
}
