package metrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

// MetricSender - интерфейс для отправки метрик.
type MetricSender interface {
	SendMetric(metricType, metricName string, value float64) error
}

// HTTPMetricSender - реализация интерфейса MetricSender, отправляющая метрики через HTTP.
type HTTPMetricSender struct {
	ServerAddress string
}

// Metrics - структура для передачи метрик в формате JSON.
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

// NewHTTPMetricSender - конструктор для HTTPMetricSender.
func NewHTTPMetricSender(serverAddress string) *HTTPMetricSender {
	return &HTTPMetricSender{ServerAddress: serverAddress}
}

// SendMetric - отправка метрики на сервер.
func (s *HTTPMetricSender) SendMetric(metricType, metricName string, value float64) error {
	url := fmt.Sprintf("http://%s/value/", s.ServerAddress)
	reqBody := Metrics{ID: metricName, MType: metricType}

	if metricType == "gauge" {
		reqBody.Value = &value
	} else if metricType == "counter" {
		delta := int64(value) // Приведение к int64
		reqBody.Delta = &delta
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("ошибка сериализации метрики: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("сшибка отправки метрики: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул некорректный статус: %v", resp.Status)
	}
	return nil
}

// Collector - структура для сбора метрик.
type Collector struct {
	pollCount int64
}

// NewCollector - конструктор для Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// CollectAndSendMetrics - сбор и отправка метрик.
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

// collectMetrics - сбор метрик системы.
func (c *Collector) collectMetrics() map[string]float64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	c.pollCount++

	return map[string]float64{
		"Alloc":        float64(memStats.Alloc),
		"BuckHashSys":  float64(memStats.BuckHashSys),
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
		"RandomValue":  rand.Float64(),
	}
}
