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

// UpdateHandler обрабатывает POST-запросы для обновления метрик
func UpdateHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Извлечение информации из URL
		path := strings.TrimPrefix(c.Request.URL.Path, "/update/")
		parts := strings.Split(path, "/")

		if len(parts) != 3 {
			c.String(http.StatusNotFound, "Не найдено")
			return
		}

		metricType, metricName, metricValue := parts[0], parts[1], parts[2]

		if metricName == "" {
			c.String(http.StatusNotFound, "Не найдено")
			return
		}

		switch metricType {
		case "gauge":
			// Преобразование значения в float64
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "Неверный запрос")
				return
			}
			// Обновление метрики типа Gauge, просто передаем значение value
			storage.UpdateGauge(metricName, value)
		case "counter":
			// Преобразование значения в int64
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "Неверный запрос")
				return
			}
			// Обновление метрики типа Counter, просто передаем значение value
			storage.UpdateCounter(metricName, value)
		default:
			c.String(http.StatusBadRequest, "Неверный запрос")
			return
		}

		// Ответ при успешном обновлении метрики
		c.String(http.StatusOK, "OK")
	}
}
