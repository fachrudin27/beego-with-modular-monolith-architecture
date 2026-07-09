package shared

import (
	"strconv"
	"time"

	beego "github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
)

func PrometheusMiddleware() beego.HandleFunc {
	return func(ctx *beecontext.Context) {
		start := time.Now()

		path, _ := ctx.Input.GetData("RouterPattern").(string)
		if path == "" {
			path = ctx.Request.URL.Path
		}

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(ctx.ResponseWriter.Status)

		httpRequestsTotal.WithLabelValues(ctx.Request.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(ctx.Request.Method, path).Observe(duration)
	}
}
