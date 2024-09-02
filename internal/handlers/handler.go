package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Favower/observability/internal/storage"
)

func UpdateHandler(storage *storage.MemStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Извлечение информации из URL
		path := strings.TrimPrefix(r.URL.Path, "/update/")
		parts := strings.Split(path, "/")

		if len(parts) != 3 {
			http.Error(w, "Не найдено", http.StatusNotFound)
			return
		}

		metricType, metricName, metricValue := parts[0], parts[1], parts[2]

		if metricName == "" {
			http.Error(w, "Не найдено", http.StatusNotFound)
			return
		}

		switch metricType {
		case "gauge":
			// Преобразование значения в float64
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Неверный запрос", http.StatusBadRequest)
				return
			}
			// Обновление метрики типа Gauge, просто передаем значение value
			storage.UpdateGauge(metricName, value)
		case "counter":
			// Преобразование значения в int64
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Неверный запрос", http.StatusBadRequest)
				return
			}
			// Обновление метрики типа Counter, просто передаем значение value
			storage.UpdateCounter(metricName, value)
		default:
			http.Error(w, "Неверный запрос", http.StatusBadRequest)
			return
		}

		// Ответ при успешном обновлении метрики
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
