package logger

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const path = "orchestrator/log/log.txt"

func SetupLogger() *zap.Logger {
	// Настраиваем конфигурацию логгера

	config := zap.NewDevelopmentConfig()

	if _, err := os.Stat(path); err == nil {
		config.OutputPaths = append(config.OutputPaths, path)
	}

	// Уровень логирования
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	// Настраиваем логгер с конфигом
	logger, err := config.Build()
	if err != nil {
		fmt.Printf("Ошибка настройки логгера: %v\n", err)
	}

	return logger
}

func LoggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	// Middleware для логирования запросов и ответов
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропускаем запрос к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершаем логирование после того, как запрос выполнен
			duration := time.Since(start)
			logger.Info(
				"HTTP запрос",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
