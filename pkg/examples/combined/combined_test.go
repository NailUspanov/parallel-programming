package combined

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
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

	// Проверяем вывод на наличие сообщений о работе воркеров
	for w := 1; w <= 3; w++ {
		workerStr := fmt.Sprintf("Воркер %d", w)
		if !strings.Contains(output, workerStr) {
			t.Errorf("Ожидаемое сообщение от воркера %d не найдено", w)
		}
	}

	// Проверяем, что есть сообщения о результатах
	if !strings.Contains(output, "Результат:") {
		t.Error("Ожидаемые сообщения о результатах не найдены")
	}

	// Проверяем, что получено правильное количество результатов
	if len(results) != 10 {
		t.Errorf("Ожидаемое количество результатов 10, получено %d", len(results))
	}

	// Проверяем, что результаты содержат правильные значения
	expectedResults := map[int]bool{
		2: false, 4: false, 6: false, 8: false, 10: false,
		12: false, 14: false, 16: false, 18: false, 20: false,
	}

	for _, r := range results {
		if _, exists := expectedResults[r]; !exists {
			t.Errorf("Неожиданный результат: %d", r)
		} else {
			expectedResults[r] = true
		}
	}

	// Проверяем, что все ожидаемые результаты были получены
	for val, found := range expectedResults {
		if !found {
			t.Errorf("Ожидаемый результат %d не найден", val)
		}
	}
}

func TestRunContextCancellation(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с отменой через 300 мс (должно успеть выполнить только часть работы)
	results := RunContextCancellation(3, 10, 300*time.Millisecond)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что контекст был отменен
	if !strings.Contains(output, "Контекст отменен:") {
		t.Error("Ожидаемое сообщение об отмене контекста не найдено")
	}

	// Проверяем, что были отправлены задания
	if !strings.Contains(output, "Отправлено задание") {
		t.Error("Ожидаемые сообщения об отправке заданий не найдены")
	}

	// Проверяем, что некоторые задания были обработаны
	if !strings.Contains(output, "Воркер") {
		t.Error("Ожидаемые сообщения о работе воркеров не найдены")
	}

	// Учитывая время отмены, должны быть получены некоторые, но не все результаты
	t.Logf("Получено результатов: %d (из 10 возможных)", len(results))

	// Проверяем, что все полученные результаты соответствуют ожидаемому формату (j*10)
	for _, r := range results {
		if r%10 != 0 || r <= 0 || r > 100 {
			t.Errorf("Неожиданное значение результата: %d", r)
		}
	}

	// Проверяем, что контекст действительно привел к отмене некоторых заданий
	if len(results) == 10 {
		t.Error("Ожидалось, что будут отменены некоторые задания, но все были выполнены")
	}

	// Теперь запускаем с большим временем таймаута, чтобы убедиться, что все задания могут быть выполнены
	oldStdout = os.Stdout
	r, w, _ = os.Pipe()
	os.Stdout = w

	completeResults := RunContextCancellation(3, 5, 5*time.Second)

	w.Close()
	os.Stdout = oldStdout

	// Проверяем, что все 5 заданий были выполнены
	if len(completeResults) != 5 {
		t.Errorf("При достаточном таймауте ожидалось 5 результатов, получено %d", len(completeResults))
	}
}

func TestRunBasicWorkerPool(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с 4 воркерами и 8 заданиями
	results := RunBasicWorkerPool(4, 8)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем вывод на наличие сообщений о работе воркеров
	for w := 1; w <= 4; w++ {
		workerStr := fmt.Sprintf("Воркер %d", w)
		if !strings.Contains(output, workerStr) {
			t.Errorf("Ожидаемое сообщение от воркера %d не найдено", w)
		}
	}

	// Проверяем, что есть сообщения о результатах
	if !strings.Contains(output, "Результат:") {
		t.Error("Ожидаемые сообщения о результатах не найдены")
	}

	// Проверяем, что получено правильное количество результатов
	if len(results) != 8 {
		t.Errorf("Ожидаемое количество результатов 8, получено %d", len(results))
	}

	// Проверяем, что результаты содержат правильные значения (j*2)
	expectedResults := map[int]bool{
		2: false, 4: false, 6: false, 8: false,
		10: false, 12: false, 14: false, 16: false,
	}

	for _, r := range results {
		if _, exists := expectedResults[r]; !exists {
			t.Errorf("Неожиданный результат: %d", r)
		} else {
			expectedResults[r] = true
		}
	}

	// Проверяем, что все ожидаемые результаты были получены
	for val, found := range expectedResults {
		if !found {
			t.Errorf("Ожидаемый результат %d не найден", val)
		}
	}
}

func TestRunWorkerPoolWithTimeout(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с коротким таймаутом (должно успеть выполнить только часть работы)
	results := RunWorkerPoolWithTimeout(3, 10, 300*time.Millisecond)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что контекст был отменен
	if !strings.Contains(output, "Контекст отменен:") {
		t.Error("Ожидаемое сообщение об отмене контекста не найдено")
	}

	if !strings.Contains(output, "context deadline exceeded") {
		t.Error("Ожидаемое сообщение о превышении времени ожидания не найдено")
	}

	// Учитывая время отмены, должны быть получены некоторые, но не все результаты
	t.Logf("Получено результатов: %d (из 10 возможных) с таймаутом", len(results))

	// Все результаты должны соответствовать ожидаемому формату (j*10)
	for _, r := range results {
		if r%10 != 0 || r <= 0 || r > 100 {
			t.Errorf("Неожиданное значение результата: %d", r)
		}
	}
}

func TestRunWorkerPoolWithCancel(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию с отменой через 300 мс
	results := RunWorkerPoolWithCancel(3, 10, 300*time.Millisecond)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем, что контекст был отменен
	if !strings.Contains(output, "Время истекло, отменяем контекст") {
		t.Error("Ожидаемое сообщение об отмене контекста не найдено")
	}

	if !strings.Contains(output, "Контекст отменен: context canceled") {
		t.Error("Ожидаемое сообщение об отмене контекста не найдено")
	}

	// Учитывая время отмены, должны быть получены некоторые, но не все результаты
	t.Logf("Получено результатов: %d (из 10 возможных) с отменой", len(results))

	// Все результаты должны соответствовать ожидаемому формату (j*10)
	for _, r := range results {
		if r%10 != 0 || r <= 0 || r > 100 {
			t.Errorf("Неожиданное значение результата: %d", r)
		}
	}
}
