package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Favower/observability/cmd/server/storage"
)

func TestUpdateHandler(t *testing.T) {
	// Создаем тестовый сервер и хранилище
	storage := storage.NewMemStorage()
	handler := UpdateHandler(storage)

	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid Gauge Metric",
			method:         "POST",
			url:            "/update/gauge/testMetric/123.45",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid Counter Metric",
			method:         "POST",
			url:            "/update/counter/testMetric/123",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid URL Format",
			method:         "POST",
			url:            "/update/gauge/testMetric",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid Metric Type",
			method:         "POST",
			url:            "/update/unknownMetric/testMetric/123",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Gauge Value",
			method:         "POST",
			url:            "/update/gauge/testMetric/invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid Counter Value",
			method:         "POST",
			url:            "/update/counter/testMetric/invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}
