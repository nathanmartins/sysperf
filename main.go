package main

import (
	"fmt"
	"github.com/google/uuid"
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
	cpuUsageGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage",
			Help: "Current usage of CPU resource",
		},
		[]string{"hostname", "session_id"},
	)
	cpuSaturationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_saturation",
			Help: "Current saturation of CPU resource",
		},
		[]string{"hostname", "session_id"},
	)
	memUsageGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "mem_usage",
		},
		[]string{"command", "hostname", "session_id"},
	)

	memSwapSaturationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_saturation_swap_usage",
			Help: "Amount of swap used by the node",
		},
		[]string{"hostname", "session_id"},
	)

	memOOMKillingSaturationGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "oom_killings",
		Help: "Amount of processes killed due to memory overload",
	})
	agentTime = 3 * time.Second
)

func init() {
	// CPU
	prometheus.MustRegister(cpuUsageGauge)
	prometheus.MustRegister(cpuSaturationGauge)

	// Memory
	prometheus.MustRegister(memUsageGauge)
	prometheus.MustRegister(memSwapSaturationGauge)
	prometheus.MustRegister(memOOMKillingSaturationGauge)
}

func main() {

	hostnameBytes, _ := os.ReadFile("/etc/hostname")
	hostname := string(hostnameBytes)
	hostname = strings.TrimSuffix(hostname, "\n")
	sessionId := uuid.New()

	// Should these be separate goroutines?
	// cpu usage and saturation
	go func() {

		for {
			sample, err := SampleCPUUsage(agentTime)
			if err != nil {
				log.Fatal(err)
			}
			cpuUsageGauge.WithLabelValues(hostname, fmt.Sprintf("%s", sessionId)).Set(sample.Usage)
			cpuSaturationGauge.WithLabelValues(hostname, fmt.Sprintf("%s", sessionId)).Set(sample.Busy)
			time.Sleep(agentTime) // Agent run interval
		}

	}()

	// memory usage
	//go func() {
	//
	//	for {
	//		samples, err := SampleMemoryUsage()
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//
	//		for _, sample := range samples {
	//			memUsageGauge.WithLabelValues(strings.TrimSuffix(sample.Command, "\n"), hostname).Set(sample.SizeKb / 1000)
	//		}
	//
	//		time.Sleep(agentTime) // Agent run interval
	//	}
	//
	//}()

	// memory saturation
	go func() {
		for {
			sample, err := SampleMemorySaturation()
			if err != nil {
				log.Fatal(err)
			}
			memSwapSaturationGauge.WithLabelValues(hostname, fmt.Sprintf("%s", sessionId)).Set(sample.SwapUsage)
			memOOMKillingSaturationGauge.Set(sample.OOMKillings)
			time.Sleep(agentTime) // Agent run interval
		}
	}()

	port := 9001
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("sysperf export is running on port: %d and your session ID is %s\n", port, sessionId)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
