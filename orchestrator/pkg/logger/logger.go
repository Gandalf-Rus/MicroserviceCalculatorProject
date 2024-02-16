package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

const path = "../log/log.txt"

func SetupInfoLogger() *log.Logger {
	// Создаем файл лога, если он не существует
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Ошибка создания файла лога: %v", err)
	}

	// Устанавливаем формат и место вывода
	logger := log.New(file, "", log.Ldate|log.Ltime)

	return logger
}

func ClearFileLog() {
	// Очищаем содержимое файла лога
	if _, err := os.Stat(path); err == nil {
		if err := os.Truncate(path, 0); err != nil {
			log.Fatalf("Ошибка очистки файла лога: %v", err)
		}
	}
}

func LoggingMiddleware(logger *log.Logger) func(next http.Handler) http.Handler {
	// Middleware для логирования запросов и ответов
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропускаем запрос к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершаем логирование после того, как запрос выполнен
			duration := time.Since(start)

			logger.Printf(
				"HTTP request: method=%s path=%s duration=%v\n",
				r.Method, r.URL.Path, duration,
			)
		})
	}
}
