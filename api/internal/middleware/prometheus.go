package middleware

import (
	"api/internal/monitoring"
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware middleware для сбора метрик HTTP запросов
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропускаем метрики Prometheus
		if c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		start := time.Now()

		// Захватываем размер запроса
		var requestSize int64 = -1
		if c.Request.ContentLength > 0 {
			requestSize = c.Request.ContentLength
		} else if c.Request.Body != nil && c.Request.Method != "GET" && c.Request.Method != "HEAD" {
			// Сохраняем оригинальное тело
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
				requestSize = int64(len(bodyBytes))
				// Восстанавливаем тело запроса
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		// Создаем кастомный writer для захвата ответа
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			statusCode:     http.StatusOK, // Статус по умолчанию
		}
		c.Writer = writer

		c.Next()

		duration := time.Since(start)
		status := writer.statusCode

		// Получаем размер ответа
		responseSize := int64(writer.body.Len())
		if responseSize == 0 {
			responseSize = -1
		}

		// Используем путь из контекста Gin (более надежно)
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// Записываем метрики
		monitoring.ObserveHTTPRequest(
			c.Request.Method,
			path,
			status,
			duration,
			requestSize,
			responseSize,
		)
	}
}

// responseWriter кастомный ResponseWriter для захвата ответа
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
	size       int
}

// Write перехватывает запись тела ответа
func (w *responseWriter) Write(b []byte) (int, error) {
	if w.body != nil {
		w.body.Write(b)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

// WriteHeader перехватывает установку статус кода
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Status возвращает статус код (для совместимости с gin)
func (w *responseWriter) Status() int {
	return w.statusCode
}

// Size возвращает размер ответа
func (w *responseWriter) Size() int {
	return w.size
}

// Body возвращает тело ответа (для отладки)
func (w *responseWriter) Body() string {
	if w.body == nil {
		return ""
	}
	return w.body.String()
}
