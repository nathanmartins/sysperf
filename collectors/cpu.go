package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CPUCollector struct{}

func (c CPUCollector) Collect(ch chan<- prometheus.Metric) {
	panic("TODO")
}

func (c CPUCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "info"),
		"CPU information from /proc/cpuinfo.",
		[]string{"package", "core", "cpu", "vendor", "family", "model", "model_name", "microcode", "stepping", "cachesize"}, nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "flag_info"),
		"The `flags` field of CPU information from /proc/cpuinfo taken from the first core.",
		[]string{"flag"}, nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "bug_info"),
		"The `bugs` field of CPU information from /proc/cpuinfo taken from the first core.",
		[]string{"bug"}, nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "guest_seconds_total"),
		"Seconds the CPUs spent in guests (VMs) for each mode.",
		[]string{"cpu", "mode"}, nil,
	)
	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "core_throttles_total"),
		"Number of times this CPU core has been throttled.",
		[]string{"package", "core"}, nil,
	)

	ch <- prometheus.NewDesc(
		prometheus.BuildFQName("sysperf", "cpu", "package_throttles_total"),
		"Number of times this CPU package has been throttled.",
		[]string{"package"}, nil,
	)
}
