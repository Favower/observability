package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	// Значение по умолчанию для адреса HTTP-сервера
	defaultAddress = "localhost:8080"
)

func main() {

	// Чтение переменной окружения ADDRESS (если существует)
	address := getEnv("ADDRESS", defaultAddress)

	// Переопределение адреса флагом командной строки
	flag.StringVar(&address, "a", address, "адрес HTTP-сервера (по умолчанию localhost:8080)")

	// Парсинг флагов
	flag.Parse()

	// Проверка на наличие неизвестных флагов
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Ошибка: неизвестные флаги: %v\n", flag.Args())
		os.Exit(1)
	}

	// Инициализация хранилища метрик
	storage := storage.NewMemStorage()

	// Инициализация роутера Gin
	r := gin.Default()

	// Маршруты
	r.GET("/value/:type/:name", handlers.GetMetricHandler(storage))
	r.POST("/update/:type/:name/:value", handlers.UpdateHandler(storage))
	r.GET("/", handlers.GetAllMetricsHandler(storage))
	r.POST("/update/", handlers.JSONUpdateMetricHandler(storage))
	r.POST("/value/", handlers.JSONGetMetricHandler(storage))

	// Запуск HTTP-сервера
	log.Printf("Запуск сервера на %s\n", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// Middleware для логирования запросов и ответов
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Выполняем запрос
		c.Next()

		// Вычисляем время выполнения запроса
		duration := time.Since(startTime)

		// Получаем статус ответа
		statusCode := c.Writer.Status()

		// Получаем размер ответа
		responseSize := c.Writer.Size()

		// Логируем информацию о запросе и ответе
		logger.Info("Request", 
			zap.String("method", c.Request.Method), 
			zap.String("uri", c.Request.RequestURI),
			zap.Duration("duration", duration),
			zap.Int("status", statusCode),
			zap.Int("response_size", responseSize),
		)
	}	
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 
