package race_condition

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestRunRaceConditionProblem(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	counter := RunRaceConditionProblem(500)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие ключевых сообщений
	if !strings.Contains(output, "Ожидаемое значение: 1000") {
		t.Error("Ожидаемое сообщение не найдено")
	}

	// Проверяем, что значение счетчика меньше ожидаемого из-за гонки
	if counter == 1000 {
		t.Log("Внимание: Счетчик достиг ожидаемого значения 1000, " +
			"что необычно при наличии гонки. Возможно, это случайность.")
	}

	t.Logf("Значение счетчика с гонкой: %d", counter)

	// Проверяем несколько раз, чтобы увеличить шанс обнаружения гонки
	racingCounter := 0
	nonRacingCounter := 0

	for i := 0; i < 5; i++ {
		result := RunRaceConditionProblem(200)
		if result < 400 { // Если значение существенно ниже ожидаемого, вероятно наблюдаем гонку
			racingCounter++
		} else if result == 400 {
			nonRacingCounter++
		}
	}

	t.Logf("Из 5 запусков: %d с явной гонкой, %d без явной гонки",
		racingCounter, nonRacingCounter)

	// Не делаем этот тест проваливающимся, так как гонка может не проявиться всегда
}

func TestRunRaceConditionSolutionWithChannels(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	counter := RunRaceConditionSolutionWithChannels(500)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем результат
	if counter != 1000 {
		t.Errorf("Ожидаемое значение счетчика 1000, получено %d", counter)
	}
}

func TestRunRaceConditionSolutionWithMutex(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	counter := RunRaceConditionSolutionWithMutex(500)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем результат
	if counter != 1000 {
		t.Errorf("Ожидаемое значение счетчика 1000, получено %d", counter)
	}
}

func TestRunRaceConditionSolutionWithAtomic(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	counter := RunRaceConditionSolutionWithAtomic(500)

	// Восстанавливаем stdout
	w.Close()
	os.Stdout = oldStdout

	// Проверяем результат
	if counter != 1000 {
		t.Errorf("Ожидаемое значение счетчика 1000, получено %d", counter)
	}
}

func TestRunAtomicCompareAndSwap(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	success, oldValue, newValue := RunAtomicCompareAndSwap()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие ключевых сообщений
	expectedMessages := []string{
		"Начальное значение",
		"Успешная операция",
		"Старое значение",
		"Новое значение",
	}

	for _, msg := range expectedMessages {
		if !strings.Contains(output, msg) {
			t.Errorf("Ожидаемое сообщение не найдено: %s", msg)
		}
	}

	// Проверяем результаты
	if !success {
		t.Error("Операция должна быть успешной")
	}

	if oldValue == newValue {
		t.Error("Старое и новое значения не должны совпадать")
	}
}

func TestRunAtomicLoad(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	value := RunAtomicLoad()

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие ключевых сообщений
	if !strings.Contains(output, "Запуск горутин для атомарного чтения/записи") {
		t.Error("Ожидаемое сообщение о запуске горутин не найдено")
	}

	// Проверяем, что значение равно ожидаемому
	if value != 200 {
		t.Errorf("Ожидаемое значение 200, получено %d", value)
	}

	// Проверка чтения значений
	re := regexp.MustCompile(`Прочитано значение: (\d+)`)
	matches := re.FindAllStringSubmatch(output, -1)

	if len(matches) < 10 {
		t.Errorf("Ожидалось не менее 10 сообщений о чтении, получено %d", len(matches))
	}
}

func TestRunRaceConditionSolution(t *testing.T) {
	// Перенаправляем stdout для перехвата вывода
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Запускаем функцию
	results := RunRaceConditionSolution(200)

	// Восстанавливаем stdout и получаем вывод
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Проверяем наличие ключевых сообщений
	expectedSolutions := []string{
		"Решение с использованием каналов",
		"Решение с использованием мьютекса",
		"Решение с использованием атомарных операций",
	}

	for _, solution := range expectedSolutions {
		if !strings.Contains(output, solution) {
			t.Errorf("Ожидаемое решение не найдено: %s", solution)
		}
	}

	// Проверяем результаты
	if len(results) != 3 {
		t.Errorf("Ожидаемое количество результатов 3, получено %d", len(results))
	}

	for _, count := range results {
		if count != 400 {
			t.Errorf("Ожидаемое значение счетчика 400, получено %d", count)
		}
	}
}

func TestRaceConditionStressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем стресс-тест в коротком режиме")
	}

	// Параметры стресс-теста
	testCases := []struct {
		name          string
		goroutines    int
		incPerRoutine int
	}{
		{"Малая нагрузка", 10, 10},
		{"Средняя нагрузка", 100, 10},
		{"Высокая нагрузка", 1000, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Для каждого случая expected - это ожидаемое значение счетчика
			// Которое равно удвоенному количеству инкрементов на горутину (из-за 2 * incPerRoutine в функции)
			expected := 2 * tc.incPerRoutine

			// Тест решения с каналами
			t.Run("Каналы", func(t *testing.T) {
				counter := RunRaceConditionSolutionWithChannels(tc.incPerRoutine)
				if counter != expected {
					t.Errorf("Каналы: ожидаемое значение %d, получено %d", expected, counter)
				}
			})

			// Тест решения с мьютексом
			t.Run("Мьютекс", func(t *testing.T) {
				counter := RunRaceConditionSolutionWithMutex(tc.incPerRoutine)
				if counter != expected {
					t.Errorf("Мьютекс: ожидаемое значение %d, получено %d", expected, counter)
				}
			})

			// Тест решения с атомарными операциями
			t.Run("Атомарные операции", func(t *testing.T) {
				counter := RunRaceConditionSolutionWithAtomic(tc.incPerRoutine)
				if counter != expected {
					t.Errorf("Атомарные операции: ожидаемое значение %d, получено %d", expected, counter)
				}
			})

			// Проверка производительности разных решений
			t.Run("Сравнение производительности", func(t *testing.T) {
				// Сохраняем время выполнения каждого решения
				timings := make(map[string]time.Duration)

				// Тест проблемного кода (без защиты от гонок)
				start := time.Now()
				RunRaceConditionProblem(tc.incPerRoutine / 2)
				problem := time.Since(start)
				timings["Без защиты"] = problem

				// Решения
				start = time.Now()
				RunRaceConditionSolutionWithChannels(tc.incPerRoutine / 2)
				channels := time.Since(start)
				timings["Каналы"] = channels

				start = time.Now()
				RunRaceConditionSolutionWithMutex(tc.incPerRoutine / 2)
				mutex := time.Since(start)
				timings["Мьютекс"] = mutex

				start = time.Now()
				RunRaceConditionSolutionWithAtomic(tc.incPerRoutine / 2)
				atomicOp := time.Since(start)
				timings["Атомарные"] = atomicOp

				// Логируем результаты производительности для анализа
				for solution, timing := range timings {
					t.Logf("%s: время выполнения %v", solution, timing)
				}

				// В идеале атомарные операции должны быть быстрее мьютексов, а мьютексы быстрее каналов
				// Но в тестовой среде это может быть не так заметно
				t.Logf("Производительность решений (меньше лучше): Атомарные: %v, Мьютекс: %v, Каналы: %v",
					atomicOp, mutex, channels)
			})
		})
	}
}
