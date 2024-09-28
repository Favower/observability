package handlers

import (
	"net/http"
	"strconv"

	"github.com/Favower/observability/internal/storage"
	"github.com/gin-gonic/gin"
)

// Определение констант для типов метрик
const (
	MetricTypeGauge   = "gauge"
	MetricTypeCounter = "counter"
)

// Хэндлер для получения значения метрики
func GetMetricHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		metricType := c.Param("type")
		metricName := c.Param("name")

		switch metricType {
		case MetricTypeGauge:
			if value, ok := storage.GetGauge(metricName); ok {
				c.String(http.StatusOK, strconv.FormatFloat(value, 'f', -1, 64))
			} else {
				c.String(http.StatusNotFound, "Metric not found")
			}
		case MetricTypeCounter:
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
			html += "<li>" + name + " (" + MetricTypeGauge + "): " + strconv.FormatFloat(value, 'f', -1, 64) + "</li>"
		}

		for name, value := range storage.Counters {
			html += "<li>" + name + " (" + MetricTypeCounter + "): " + strconv.FormatInt(value, 10) + "</li>"
		}

		html += "</ul></body></html>"

		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}

func UpdateHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получение параметров из URL
		metricType := c.Param("type")
		metricName := c.Param("name")
		metricValue := c.Param("value")

		// Проверка, что имя метрики не пустое
		if metricName == "" {
			c.String(http.StatusNotFound, "Не найдено")
			return
		}

		// Обработка разных типов метрик
		switch metricType {
		case MetricTypeGauge:
			// Преобразование значения в float64
			value, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "Неверный запрос")
				return
			}
			// Обновление метрики типа Gauge
			storage.UpdateGauge(metricName, value)
		case MetricTypeCounter:
			// Преобразование значения в int64
			value, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				c.String(http.StatusBadRequest, "Неверный запрос")
				return
			}
			// Обновление метрики типа Counter
			storage.UpdateCounter(metricName, value)
		default:
			c.String(http.StatusBadRequest, "Неверный тип метрики")
			return
		}

		// Успешный ответ
		c.String(http.StatusOK, "OK")
	}
}
