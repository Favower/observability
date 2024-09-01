package metrics

import (
	"fmt"
	"net/http"
	"time"
	"runtime"
)

// MetricSender - интерфейс для отправки метрик.
type MetricSender interface {
	SendMetric(metricType, metricName string, value float64)
}

// HTTPMetricSender - реализация интерфейса MetricSender, отправляющая метрики через HTTP.
type HTTPMetricSender struct {
	ServerAddress string
}

func NewHTTPMetricSender(serverAddress string) *HTTPMetricSender {
	return &HTTPMetricSender{ServerAddress: serverAddress}
}

func (s *HTTPMetricSender) SendMetric(metricType, metricName string, value float64) {
	url := fmt.Sprintf("http://%s/update/%s/%s/%v", s.ServerAddress, metricType, metricName, value)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("Ошибка создания запроса: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Ошибка отправки метрики: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Сервер вернул некорректный статус: %v\n", resp.Status)
	}
}

type Collector struct {
	pollCount int64
}

func NewCollector() *Collector {
	return &Collector{}
}

func (c *Collector) CollectAndSendMetrics(sender MetricSender, pollInterval, reportInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		metrics := c.collectMetrics()
		for name, value := range metrics {
			metricType := "gauge"
			if name == "PollCount" {
				metricType = "counter"
			}
			sender.SendMetric(metricType, name, value)
		}
	}
}

func (c *Collector) collectMetrics() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.pollCount++

	return map[string]float64{
		"Alloc":        float64(memStats.Alloc),
		"TotalAlloc":   float64(memStats.TotalAlloc),
		"Sys":          float64(memStats.Sys),
		"Lookups":      float64(memStats.Lookups),
		"Mallocs":      float64(memStats.Mallocs),
		"Frees":        float64(memStats.Frees),
		"HeapAlloc":    float64(memStats.HeapAlloc),
		"HeapSys":      float64(memStats.HeapSys),
		"HeapIdle":     float64(memStats.HeapIdle),
		"HeapInuse":    float64(memStats.HeapInuse),
		"HeapReleased": float64(memStats.HeapReleased),
		"HeapObjects":  float64(memStats.HeapObjects),
		"StackInuse":   float64(memStats.StackInuse),
		"StackSys":     float64(memStats.StackSys),
		"MSpanInuse":   float64(memStats.MSpanInuse),
		"MSpanSys":     float64(memStats.MSpanSys),
		"MCacheInuse":  float64(memStats.MCacheInuse),
		"MCacheSys":    float64(memStats.MCacheSys),
		"GCSys":        float64(memStats.GCSys),
		"OtherSys":     float64(memStats.OtherSys),
		"NextGC":       float64(memStats.NextGC),
		"LastGC":       float64(memStats.LastGC),
		"PauseTotalNs": float64(memStats.PauseTotalNs),
		"NumGC":        float64(memStats.NumGC),
		"NumForcedGC":  float64(memStats.NumForcedGC),
		"GCCPUFraction": memStats.GCCPUFraction,
		"PollCount":    float64(c.pollCount),
	}
}
