package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Collector is the interface a collector has to implement.
type Collector interface {
	Update(ch chan<- prometheus.Metric) error // Get new metrics and expose them via prometheus registry.
}
