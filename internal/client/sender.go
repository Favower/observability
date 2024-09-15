package client

import (
	"fmt"
	"net/http"
)

type MetricSender interface {
	SendMetric(metricType, name string, value float64) error
}

type Sender struct {
	serverURL string
}

func NewSender(serverURL string) *Sender {
	return &Sender{serverURL: serverURL}
}

func (s *Sender) SendMetric(metricType, name string, value float64) error {
	// Формируем URL для отправки метрики
	url := fmt.Sprintf("%s/update/%s/%s/%f", s.serverURL, metricType, name, value)

	// Выполняем POST-запрос
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		return fmt.Errorf("failed to send metric %s: %v", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server error for %s: %s", name, resp.Status)
	}
	
	return nil
}
