package client

import (
	"fmt"
	"net/http"
)

type MetricSender interface {
	SendMetric(metricType, name string, value float64)
}

type Sender struct {
	serverURL string
}

func NewSender(serverURL string) *Sender {
	return &Sender{serverURL: serverURL}
}

func (s *Sender) SendMetric(metricType, name string, value float64) {
	url := fmt.Sprintf("%s/update/%s/%s/%f", s.serverURL, metricType, name, value)
	resp, err := http.Post(url, "text/plain", nil)
	if err != nil {
		fmt.Printf("Failed to send metric %s: %v\n", name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error response from server for %s: %s\n", name, resp.Status)
	}
}
