package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
)

func main() {
	storage := storage.NewMemStorage()

	r := gin.Default()

	// Маршрут для получения значения метрики
	r.GET("/value/:type/:name", handlers.GetMetricHandler(storage))

	// Маршрут для отображения всех метрик в HTML
	r.GET("/", handlers.GetAllMetricsHandler(storage))

	// Запуск сервера на порту 8080
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
