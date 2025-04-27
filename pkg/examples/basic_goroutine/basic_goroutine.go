package basic_goroutine

import (
	"fmt"
	"time"
)

// RunBasicGoroutine демонстрирует базовый пример использования горутин
func RunBasicGoroutine() {
	go func() {
		// Код, который будет выполняться в горутине
		fmt.Println("Это выполняется в горутине")
	}()

	// Код основной функции продолжает выполняться параллельно
	fmt.Println("Это выполняется в основной функции")

	// Даем горутине время выполниться
	time.Sleep(100 * time.Millisecond)
}
