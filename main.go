package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
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
)

func init() {
	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(cpuSaturation)
	prometheus.MustRegister(cpuSaturationBusy)
}

func main() {

	go func() {

		for {
			sample, err := SampleCPUSaturation(3 * time.Second)
			if err != nil {
				log.Fatal(err)
			}

			cpuSaturation.Set(sample.Usage)
			cpuSaturationBusy.Set(sample.Busy)

			time.Sleep(10 * time.Second) // Agent run interval
		}

	}()

	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9001", nil))
}
