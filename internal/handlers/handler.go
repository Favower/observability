package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/Favower/observability/internal/storage"
	"github.com/gin-gonic/gin"
)

// Хэндлер для получения значения метрики
func GetMetricHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")

		switch metricType {
		case "gauge":
			if value, ok := storage.GetGauge(metricName); ok {
				c.String(http.StatusOK, strconv.FormatFloat(value, 'f', -1, 64))
			} else {
				c.String(http.StatusNotFound, "Metric not found")
			}
		case "counter":
			if value, ok := storage.GetCounter(metricName); ok {
				c.String(http.StatusOK, strconv.FormatInt(value, 10))
			} else {
				c.String(http.StatusNotFound, "Metric not found")
			}
		default:
			c.String(http.StatusBadRequest, "Invalid metric type")
		}
	}
}

// Хэндлер для отображения всех метрик в HTML
func GetAllMetricsHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		storage.Mu.RLock()
		defer storage.Mu.RUnlock()

		html := "<html><body><h1>Metrics</h1><ul>"

		for name, value := range storage.Gauges {
			html += "<li>" + name + " (gauge): " + strconv.FormatFloat(value, 'f', -1, 64) + "</li>"
		}

		for name, value := range storage.Counters {
			html += "<li>" + name + " (counter): " + strconv.FormatInt(value, 10) + "</li>"
		}

		html += "</ul></body></html>"

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

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
