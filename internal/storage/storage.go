package storage

import (
	"sync"
)

// MemStorage структура для хранения метрик
type MemStorage struct {
	Mu       sync.RWMutex
	Gauges   map[string]float64
	Counters map[string]int64
}

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
 }

// NewMemStorage возвращает новый экземпляр MemStorage
func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]int64),
	}
}

// GetGauge возвращает значение метрики типа Gauge
func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	value, ok := m.Gauges[name]
	return value, ok
}

// GetCounter возвращает значение метрики типа Counter
func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.Mu.RLock()
	defer m.Mu.RUnlock()
	value, ok := m.Counters[name]
	return value, ok
}


// Обновление значения метрики типа Gauge
func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Gauges[name] = value
}

// Обновление значения метрики типа Counter
func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.Mu.Lock()
	defer m.Mu.Unlock()
	m.Counters[name] += value
}
