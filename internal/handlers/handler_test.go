package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/Favower/observability/internal/handlers"
	"github.com/Favower/observability/internal/storage"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
	}{
		{
			name:           "Valid gauge update",
			method:         "POST",
			url:            "/update/gauge/temperature/23.5",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Valid counter update",
			method:         "POST",
			url:            "/update/counter/hits/1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid metric type",
			method:         "POST",
			url:            "/update/invalidType/metric/10",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Incomplete URL",
			method:         "POST",
			url:            "/update/gauge/temperature",
			expectedStatus: http.StatusNotFound,
		},
	}

	// Создаем экземпляр хранилища для тестов
	memStorage := storage.NewMemStorage()

	// Создаем новый роутер Gin
	router := gin.New()

	// Регистрируем хендлер
	router.POST("/update/:type/:name/:value", handlers.UpdateHandler(memStorage))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем новый HTTP-запрос с использованием метода и URL из тестов
			req := httptest.NewRequest(tt.method, tt.url, nil)
			
			// Создаем рекордер для записи HTTP-ответа
			rr := httptest.NewRecorder()

			// Выполняем тестируемый хендлер
			router.ServeHTTP(rr, req)

			// Проверяем статус-код ответа
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v, want %v", status, tt.expectedStatus)
			}
		})
	}
}
