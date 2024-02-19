package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

func SetupInfoLogger() *log.Logger {
	// Устанавливаем формат и место вывода
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	return logger
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

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
