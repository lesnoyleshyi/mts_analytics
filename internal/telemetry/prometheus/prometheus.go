package prometheus

import "github.com/prometheus/client_golang/prometheus"

var (
	errTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "errCount",
			Help: "Counter of all service errors",
		})

	httpTotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "httpRequestsCounter",
			Help: "Counter of http requests with labels(method, statusCode)",
		},
		[]string{"status_code", "method"})

	kafkaTailLength = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kafka_tail_length",
			Help: "Count of unprocessed kafka messages",
		})
)

func Initialise() {
	prometheus.MustRegister(errTotal)
	prometheus.MustRegister(httpTotalRequests)
	prometheus.MustRegister(kafkaTailLength)
}
