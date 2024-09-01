package client

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Favower/observability/internal/metrics"
)

func TestHTTPMetricSender_SendMetric(t *testing.T) {
	// Создаем mock-сервер
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %v", r.Method)
		}
		if r.URL.Path != "/update/gauge/testMetric/123.45" {
			t.Errorf("Expected URL path /update/gauge/testMetric/123.45, got %v", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "text/plain" {
			t.Errorf("Expected Content-Type text/plain, got %v", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	})
	server := httptest.NewServer(handler)
	defer server.Close()

	sender := metrics.NewHTTPMetricSender(server.URL)
	sender.SendMetric("gauge", "testMetric", 123.45)
}
