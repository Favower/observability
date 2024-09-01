package main

import (
	"time"

	"github.com/Favower/observability/internal/metrics"
	"github.com/Favower/observability/internal/client"
)

func main() {
	pollInterval := 2 * time.Second
	reportInterval := 10 * time.Second

	collector := metrics.NewCollector()
	sender := client.NewSender("http://localhost:8080")

	go collector.CollectAndSendMetrics(sender, pollInterval, reportInterval)

	select {}
}
