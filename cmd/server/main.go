package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
	"github.com/gin-gonic/gin"
)

func main() {
	// Значение по умолчанию для адреса HTTP-сервера
	defaultAddress := "localhost:8080"

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
	r.GET("/", handlers.GetAllMetricsHandler(storage))
	r.PUT("/update/:type/:name/:value", handlers.UpdateHandler(storage))

	// Запуск HTTP-сервера
	log.Printf("Запуск сервера на %s\n", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 
