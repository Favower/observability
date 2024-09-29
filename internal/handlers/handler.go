package handlers

import (
	"encoding/json"
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

var metric storage.MetricsForJson

// UpdateMetricHandler принимает метрики в формате JSON и обновляет их
func JsonUpdateMetricHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Чтение и парсинг JSON
		if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// Проверка типа метрики
		switch metric.MType {
		case MetricTypeGauge:
			if metric.Value == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing value for gauge"})
				return
			}
			storage.UpdateGauge(metric.ID, *metric.Value)

		case MetricTypeCounter:
			if metric.Delta == nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Missing delta for counter"})
				return
			}
			storage.UpdateCounter(metric.ID, *metric.Delta)

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		// Отправка обновленной метрики в ответе
		c.JSON(http.StatusOK, metric)
	}
}

// JsonGetMetricHandler принимает запросы на получение метрик в формате JSON и возвращает их значения
func JsonGetMetricHandler(storage *storage.MemStorage) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Чтение и парсинг JSON
		if err := json.NewDecoder(c.Request.Body).Decode(&metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// Получение значения метрики в зависимости от типа
		switch metric.MType {
		case "gauge":
			if value, ok := storage.GetGauge(metric.ID); ok {
				metric.Value = &value // Присваиваем значение в поле Value
				metric.Delta = nil    // Убираем Delta, так как это не применимо для gauge
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
				return
			}
		case "counter":
			if value, ok := storage.GetCounter(metric.ID); ok {
				metric.Value = nil    // Убираем Value, так как это не применимо для counter
				metric.Delta = &value // Присваиваем значение в поле Delta
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "Metric not found"})
				return
			}
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid metric type"})
			return
		}

		// Отправка ответа с актуальными значениями метрики
		c.JSON(http.StatusOK, metric)
	}
}
