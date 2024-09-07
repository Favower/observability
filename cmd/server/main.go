package main

import (
	"flag"
	"fmt"
	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	// Определение флага -a для указания адреса HTTP-сервера
	address := flag.String("a", "localhost:8080", "адрес HTTP-сервера (по умолчанию localhost:8080)")

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
	r.POST("/update/:type/:name/:value", handlers.UpdateHandler(storage))

	// Запуск HTTP-сервера
	log.Printf("Запуск сервера на %s\n", *address)
	if err := r.Run(*address); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
