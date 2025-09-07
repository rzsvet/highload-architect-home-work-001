package monitoring

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTP метрики
	HttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	HttpRequestSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_size_bytes",
			Help: "Size of HTTP requests",
		},
		[]string{"method", "path"},
	)

	HttpResponseSize = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_response_size_bytes",
			Help: "Size of HTTP responses",
		},
		[]string{"method", "path", "status"},
	)

	// Бизнес метрики
	UsersRegistered = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "users_registered_total",
			Help: "Total number of registered users",
		},
	)

	UserLogins = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_logins_total",
			Help: "Total number of user logins",
		},
		[]string{"status"},
	)

	DatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "database_query_duration_seconds",
			Help:    "Duration of database queries",
			Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1, 2},
		},
		[]string{"operation", "success"},
	)

	// Системные метрики
	GoroutinesCount = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "goroutines_count",
			Help: "Current number of goroutines",
		},
	)

	MemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Memory usage in bytes",
		},
		[]string{"type"},
	)
)

// ObserveHTTPRequest записывает метрики HTTP запроса
func ObserveHTTPRequest(method, path string, status int, duration time.Duration, requestSize, responseSize int64) {
	statusStr := strconv.Itoa(status)

	HttpRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())
	HttpRequestsTotal.WithLabelValues(method, path, statusStr).Inc()
	HttpRequestSize.WithLabelValues(method, path).Observe(float64(requestSize))
	HttpResponseSize.WithLabelValues(method, path, statusStr).Observe(float64(responseSize))
}

// ObserveDatabaseQuery записывает метрики запроса к базе данных
func ObserveDatabaseQuery(operation string, success bool, duration time.Duration) {
	successStr := strconv.FormatBool(success)
	DatabaseQueryDuration.WithLabelValues(operation, successStr).Observe(duration.Seconds())
}

// RecordUserRegistration увеличивает счетчик зарегистрированных пользователей
func RecordUserRegistration() {
	UsersRegistered.Inc()
}

// RecordUserLogin увеличивает счетчик логинов
func RecordUserLogin(success bool) {
	status := "failed"
	if success {
		status = "success"
	}
	UserLogins.WithLabelValues(status).Inc()
}
