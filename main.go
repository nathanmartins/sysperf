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
	cpuSaturation = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_saturation",
		Help: "Current temperature of the CPU.",
	})
	cpuSaturationBusy = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_saturation_busy",
		Help: "Current temperature of the CPU.",
	})
	cpuSaturationTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_saturation_total",
		Help: "Current temperature of the CPU.",
	})
	cpuLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "cpu_latency",
		},
		[]string{"command", "hostname"},
	)
	cpuLatencySpent = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "spent_cpu_latency",
		},
		[]string{"command", "hostname"},
	)
	memLatency = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_latency",
		},
		[]string{"command", "hostname"},
	)
)

func init() {
	prometheus.MustRegister(cpuSaturation)
	prometheus.MustRegister(cpuSaturationBusy)
	prometheus.MustRegister(cpuSaturationTotal)
	prometheus.MustRegister(cpuLatency)
	prometheus.MustRegister(cpuLatencySpent)
	prometheus.MustRegister(memLatency)
}

func main() {

	hostnameBytes, _ := os.ReadFile("/etc/hostname")
	hostname := string(hostnameBytes)
	hostname = strings.TrimSuffix(hostname, "\n")

	go func() {

		for {
			sample, err := SampleCPUSaturation(3 * time.Second)
			if err != nil {
				log.Fatal(err)
			}

			cpuSaturation.Set(sample.Usage)
			cpuSaturationBusy.Set(sample.Busy)

			time.Sleep(5 * time.Second) // Agent run interval
		}

	}()

	go func() {

		for {
			samples, err := SampleCPULatency()
			if err != nil {
				log.Fatal(err)
			}

			for _, sample := range samples {
				cpuLatency.WithLabelValues(sample.Command, hostname).Observe(sample.RunQueueLatency)
				cpuLatencySpent.WithLabelValues(sample.Command, hostname).Set(sample.TimeSpentOnCPU)
			}

			time.Sleep(1 * time.Second) // Agent run interval

		}

	}()

	go func() {

		for {
			samples, err := SampleMemoryLatency()
			if err != nil {
				log.Fatal(err)
			}

			for _, sample := range samples {
				memLatency.WithLabelValues(sample.Command, hostname).Set(sample.SizeKb / 1000)
			}

			time.Sleep(1 * time.Second) // Agent run interval
		}

	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9001", nil))
}
