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
	cpuSaturationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_saturation",
		Help: "Current saturation of CPU resource",
	})
	memUsageGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_usage",
		},
		[]string{"command", "hostname"},
	)
	agentTime = 3 * time.Second
)

func init() {
	// CPU
	prometheus.MustRegister(cpuUsageGauge)
	prometheus.MustRegister(cpuSaturationGauge)

	// Memory
	prometheus.MustRegister(memUsageGauge)
}

func main() {

	hostnameBytes, _ := os.ReadFile("/etc/hostname")
	hostname := string(hostnameBytes)
	hostname = strings.TrimSuffix(hostname, "\n")

	// Should these be separate goroutines?
	go func() {

		for {
			sample, err := SampleCPUUsage(agentTime)
			if err != nil {
				log.Fatal(err)
			}
			cpuUsageGauge.Set(sample.Usage)
			cpuSaturationGauge.Set(sample.Busy)
			time.Sleep(agentTime) // Agent run interval
		}

	}()

	go func() {

		for {
			samples, err := SampleMemoryUsage()
			if err != nil {
				log.Fatal(err)
			}

			for _, sample := range samples {
				memUsageGauge.WithLabelValues(strings.TrimSuffix(sample.Command, "\n"), hostname).Set(sample.SizeKb / 1000)
			}

			time.Sleep(agentTime) // Agent run interval
		}

	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("sysperf prometheus exported started and is running on port: 9001")
	log.Fatal(http.ListenAndServe(":9001", nil))
}
