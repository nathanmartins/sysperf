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
	cpuSaturationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_saturation",
		Help: "Current saturation of our CPU",
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
	memLatencyGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_latency",
		},
		[]string{"command", "hostname"},
	)
	agentTime = 5 * time.Second
)

func init() {
	prometheus.MustRegister(cpuSaturationGauge)
	prometheus.MustRegister(cpuLatencySpentGauge)
	prometheus.MustRegister(memLatencyGauge)
	prometheus.MustRegister(cpuLatencyHistogram)
}

func main() {

	hostnameBytes, _ := os.ReadFile("/etc/hostname")
	hostname := string(hostnameBytes)
	hostname = strings.TrimSuffix(hostname, "\n")

	go func() {

		for {
			sample, err := SampleCPUSaturation(agentTime)
			if err != nil {
				log.Fatal(err)
			}
			cpuSaturationGauge.Set(sample.Usage)
			time.Sleep(agentTime) // Agent run interval
		}

	}()

	go func() {

		for {
			samples, err := SampleCPULatency()
			if err != nil {
				log.Fatal(err)
			}

			for _, sample := range samples {
				cpuLatencyHistogram.WithLabelValues(sample.Command, hostname).Observe(sample.RunQueueLatency)
				cpuLatencySpentGauge.WithLabelValues(sample.Command, hostname).Set(sample.TimeSpentOnCPU)
			}

			time.Sleep(agentTime) // Agent run interval

		}

	}()

	go func() {

		for {
			samples, err := SampleMemoryLatency()
			if err != nil {
				log.Fatal(err)
			}

			for _, sample := range samples {
				memLatencyGauge.WithLabelValues(sample.Command, hostname).Set(sample.SizeKb / 1000)
			}

			time.Sleep(agentTime) // Agent run interval
		}

	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9001", nil))
}
