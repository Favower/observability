package storage

import (
	"sync"
)

type MemStorage struct {
	mu      sync.RWMutex
	gauges  map[string]float64
	counters map[string]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:  make(map[string]float64),
		counters: make(map[string]int64),
	}
}

// Обновление значения метрики типа Gauge
func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}

// Обновление значения метрики типа Counter
func (m *MemStorage) UpdateCounter(name string, value int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += value
}
