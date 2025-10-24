package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CacheHitsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hits_total",
			Help: "Total number of cache hits",
		},
		[]string{"type"},
	)

	CacheMissesTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_misses_total",
			Help: "Total number of cache misses",
		},
		[]string{"type"},
	)

	CacheInvalidationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_invalidations_total",
			Help: "Total number of cache invalidations",
		},
		[]string{"type"},
	)

	CacheSize = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cache_size_bytes",
			Help: "Size of cache in bytes",
		},
		[]string{"type"},
	)
)

// RecordCacheHit записывает попадание в кэш
func RecordCacheHit(cacheType string) {
	CacheHitsTotal.WithLabelValues(cacheType).Inc()
}

// RecordCacheMiss записывает промах кэша
func RecordCacheMiss(cacheType string) {
	CacheMissesTotal.WithLabelValues(cacheType).Inc()
}

// RecordCacheInvalidation записывает инвалидацию кэша
func RecordCacheInvalidation(cacheType string) {
	CacheInvalidationsTotal.WithLabelValues(cacheType).Inc()
}
