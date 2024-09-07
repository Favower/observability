package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Favower/observability/internal/client"
	"github.com/Favower/observability/internal/metrics"
)

func main() {
	// Определение флагов для адреса сервера, интервала опроса метрик и интервала отправки
	defaultAddress := "localhost:8080"
	defaultPollInterval := 2    // в секундах
	defaultReportInterval := 10 // в секундах

	// Чтение переменных окружения с приоритетом
	address := getEnv("ADDRESS", defaultAddress)
	pollInterval := getEnvAsInt("POLL_INTERVAL", defaultPollInterval)
	reportInterval := getEnvAsInt("REPORT_INTERVAL", defaultReportInterval)

	// Переопределение флагами командной строки
	flag.StringVar(&address, "a", address, "Адрес HTTP-сервера (по умолчанию localhost:8080)")
	flag.IntVar(&pollInterval, "p", pollInterval, "Частота опроса метрик в секундах (по умолчанию 2 сек.)")
	flag.IntVar(&reportInterval, "r", reportInterval, "Частота отправки метрик на сервер в секундах (по умолчанию 10 сек.)")

	// Парсинг флагов
	flag.Parse()

	// Проверка на неизвестные флаги
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Ошибка: неизвестные флаги: %v\n", flag.Args())
		os.Exit(1)
	}

	// Преобразуем интервалы в тип `time.Duration`
	pollDuration := time.Duration(pollInterval) * time.Second
	reportDuration := time.Duration(reportInterval) * time.Second

	// Инициализация коллектора и отправителя метрик
	collector := metrics.NewCollector()
	sender := client.NewSender(fmt.Sprintf("http://%s", address))

	// Запуск сбора и отправки метрик
	go collector.CollectAndSendMetrics(sender, pollDuration, reportDuration)

	// Бесконечный цикл, чтобы программа не завершалась
	select {}
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt возвращает значение переменной окружения как int или значение по умолчанию
func getEnvAsInt(name string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(name); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
