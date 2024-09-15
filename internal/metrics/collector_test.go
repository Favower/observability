package metrics

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// MockMetricSender - Mock-реализация MetricSender для тестов
type MockMetricSender struct {
	ReceivedMetrics map[string]float64
	mu              sync.Mutex
}

// NewMockMetricSender создает новый MockMetricSender
func NewMockMetricSender() *MockMetricSender {
	return &MockMetricSender{
		ReceivedMetrics: make(map[string]float64),
	}
}

func (m *MockMetricSender) SendMetric(metricType, name string, value float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ReceivedMetrics[name] = value
	return nil
}

func TestCollector_CollectAndSendMetrics(t *testing.T) {
	mockSender := NewMockMetricSender()
	collector := NewCollector()

	// Запускаем сбор метрик и отправку в отдельной горутине
	go collector.CollectAndSendMetrics(mockSender, 1*time.Second, 10*time.Second)

	// Ждем немного, чтобы метрики были собраны и отправлены
	time.Sleep(2 * time.Second)

	// Проверяем, что метрики были собраны
	mockSender.mu.Lock()
	defer mockSender.mu.Unlock()

	assert.NotEmpty(t, mockSender.ReceivedMetrics, "Expected metrics to be collected and sent")

	// Проверяем, что среди собранных метрик есть "Alloc"
	_, exists := mockSender.ReceivedMetrics["Alloc"]
	assert.True(t, exists, "Expected metric 'Alloc' to be present")

	// Проверяем, что метрика PollCount тоже присутствует
	_, exists = mockSender.ReceivedMetrics["PollCount"]
	assert.True(t, exists, "Expected metric 'PollCount' to be present")
}
