package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Favower/observability/internal/client"
	"github.com/Favower/observability/internal/metrics"
)

func main() {
	// Определение флагов для адреса сервера, интервала опроса метрик и интервала отправки
	address := flag.String("a", "localhost:8080", "Адрес HTTP-сервера (по умолчанию localhost:8080)")
	pollInterval := flag.Int("p", 2, "Частота опроса метрик в секундах (по умолчанию 2 сек.)")
	reportInterval := flag.Int("r", 10, "Частота отправки метрик на сервер в секундах (по умолчанию 10 сек.)")

	// Парсинг флагов
	flag.Parse()

	// Проверка на неизвестные флаги
	if len(flag.Args()) > 0 {
		fmt.Fprintf(os.Stderr, "Ошибка: неизвестные флаги: %v\n", flag.Args())
		os.Exit(1)
	}

	// Преобразуем интервалы в тип `time.Duration`
	pollDuration := time.Duration(*pollInterval) * time.Second
	reportDuration := time.Duration(*reportInterval) * time.Second

	// Инициализация коллектора и отправителя метрик
	collector := metrics.NewCollector()
	sender := client.NewSender(fmt.Sprintf("http://%s", *address))

	// Запуск сбора и отправки метрик
	go collector.CollectAndSendMetrics(sender, pollDuration, reportDuration)

	// Бесконечный цикл, чтобы программа не завершалась
	select {}
}
