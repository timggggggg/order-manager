package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ordersCreated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orders_created_total",
		Help: "Total number of created orders",
	})

	ordersByStatus = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "orders_current_by_status",
			Help: "Current number of orders grouped by status",
		},
		[]string{"status"},
	)

	orderErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_errors_total",
			Help: "Total number of order processing errors",
		},
		[]string{"type"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"handler", "method"},
	)

	revenueTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "revenue_total",
			Help: "Total revenue from completed orders",
		},
	)
)

func IncrementOrdersCreated() {
	ordersCreated.Inc()
}

func UpdateOrderStatus(status string) {
	ordersByStatus.WithLabelValues(status).Inc()
}

func IncrementErrorCounter(errorType string) {
	orderErrors.WithLabelValues(errorType).Inc()
}

func AddToRevenue(amount float64) {
	revenueTotal.Add(amount)
}
