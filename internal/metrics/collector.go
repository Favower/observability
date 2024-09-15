package metrics

import (
	"fmt"
	"net/http"
	"time"
	"runtime"
	"math/rand"
)

// MetricSender - интерфейс для отправки метрик.
type MetricSender interface {
	SendMetric(metricType, metricName string, value float64) error
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
		// Сбор метрик
		metrics := c.collectMetrics()

		// Отправка каждой метрики
		for name, value := range metrics {
			metricType := "gauge"
			if name == "PollCount" {
				metricType = "counter"
			}

			// Попытка отправки метрики и обработка возможной ошибки
			if err := sender.SendMetric(metricType, name, value); err != nil {
				// Логирование ошибки отправки метрики
				fmt.Printf("Ошибка при отправке метрики %s (%s): %v\n", name, metricType, err)
			}
		}
	}
}


func (c *Collector) collectMetrics() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.pollCount++

	return map[string]float64{
		"Alloc":        float64(memStats.Alloc),
		"BuckHashSys":	float64(memStats.BuckHashSys),
		"Frees":        float64(memStats.Frees),
		"GCCPUFraction": memStats.GCCPUFraction,
		"GCSys":        float64(memStats.GCSys),
		"HeapAlloc":    float64(memStats.HeapAlloc),
		"HeapIdle":     float64(memStats.HeapIdle),
		"HeapInuse":    float64(memStats.HeapInuse),
		"HeapObjects":  float64(memStats.HeapObjects),
		"HeapReleased": float64(memStats.HeapReleased),
		"HeapSys":      float64(memStats.HeapSys),
		"LastGC":       float64(memStats.LastGC),
		"Lookups":      float64(memStats.Lookups),
		"MCacheInuse":  float64(memStats.MCacheInuse),
		"MCacheSys":    float64(memStats.MCacheSys),
		"MSpanInuse":   float64(memStats.MSpanInuse),
		"MSpanSys":     float64(memStats.MSpanSys),
		"NextGC":       float64(memStats.NextGC),
		"NumForcedGC":  float64(memStats.NumForcedGC),
		"NumGC":        float64(memStats.NumGC),
		"OtherSys":     float64(memStats.OtherSys),
		"PauseTotalNs": float64(memStats.PauseTotalNs),		
		"StackInuse":   float64(memStats.StackInuse),
		"StackSys":     float64(memStats.StackSys),
		"Sys":          float64(memStats.Sys),
		"TotalAlloc":   float64(memStats.TotalAlloc),
		"PollCount":    float64(c.pollCount),
		"RandomValue":	rand.Float64(),
	}
}
