package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

var namespace = "propaganda"

// Metrics provided through prometheus
var (
	bootTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "boot_time_seconds",
		Help:      "unix timestamp of when the service was started",
	})

	Up = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "up",
		Help:      "wether the service is up or not",
	})
	WebhooksReceived = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "received_total",
		Help:      "total number of received webhooks",
	})
	WebhooksBytesRead = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "bytes_read_total",
		Help:      "total number of incoming bytes",
	})
	WebhooksErrors = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "errors_total",
		Help:      "total number of webhooks errors",
	})
	WebhooksInvalid = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "invalid_total",
		Help:      "total number of invalid webhooks",
	}, []string{"reason"})
	WebhooksValid = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "webhooks",
		Name:      "valid_total",
		Help:      "total number of valid webhooks",
	}, []string{"project"})
	AnnouncementSuccesses = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "announcer",
		Name:      "success_total",
		Help:      "total number of announcement successes",
	}, []string{"project"})
	AnnouncementErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: "announcer",
		Name:      "errors_total",
		Help:      "total number of announcement errors",
	}, []string{"status"})
)

// Register registers all the metrics and sets the http handler
func Register(metricsPath string) {
	bootTime.Set(float64(time.Now().Unix()))
	Up.Set(0)

	prometheus.MustRegister(bootTime)
	prometheus.MustRegister(Up)
	prometheus.MustRegister(WebhooksReceived)
	prometheus.MustRegister(WebhooksBytesRead)
	prometheus.MustRegister(WebhooksErrors)
	prometheus.MustRegister(WebhooksInvalid)
	prometheus.MustRegister(WebhooksValid)
	prometheus.MustRegister(AnnouncementSuccesses)
	prometheus.MustRegister(AnnouncementErrors)

	http.Handle(metricsPath, prometheus.Handler())
}
