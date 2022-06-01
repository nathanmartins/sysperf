package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// CPU
	cpuUsageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_usage",
		Help: "Current usage of CPU resource",
	})
	cpuLatencyHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cpu_latency",
		},
		[]string{"command", "hostname"},
	)
	cpuLatencySpentGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spent_cpu_latency",
		},
		[]string{"command", "hostname"},
	)
	// Memory
	memLatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_latency",
		},
		[]string{"command", "hostname"},
	)
	agentTime = 5 * time.Second
)

func init() {
	// CPU
	prometheus.MustRegister(cpuUsageGauge)
	prometheus.MustRegister(cpuLatencySpentGauge)
	prometheus.MustRegister(cpuLatencyHistogram)

	// Memory
	prometheus.MustRegister(memLatencyGauge)
}

func main() {

	hostnameBytes, _ := os.ReadFile("/etc/hostname")
	hostname := string(hostnameBytes)
	hostname = strings.TrimSuffix(hostname, "\n")

	go func() {

		for {
			sample, err := SampleCPUUsage(agentTime)
			if err != nil {
				log.Fatal(err)
			}
			cpuUsageGauge.Set(sample.Usage)
			time.Sleep(agentTime) // Agent run interval
		}

	}()

	//go func() {
	//
	//	for {
	//		samples, err := SampleCPULatency()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		for _, sample := range samples {
	//			cpuLatencyHistogram.WithLabelValues(sample.Command, hostname).Observe(sample.RunQueueLatency)
	//			cpuLatencySpentGauge.WithLabelValues(sample.Command, hostname).Set(sample.TimeSpentOnCPU)
	//		}
	//
	//		time.Sleep(agentTime) // Agent run interval
	//
	//	}
	//
	//}()

	//go func() {
	//
	//	for {
	//		samples, err := SampleMemoryLatency()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		for _, sample := range samples {
	//			memLatencyGauge.WithLabelValues(sample.Command, hostname).Set(sample.SizeKb / 1000)
	//		}
	//
	//		time.Sleep(agentTime) // Agent run interval
	//	}
	//
	//}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("sysperf prometheus exported started and is running on port: 9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
