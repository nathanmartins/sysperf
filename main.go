package main

import (
	"github.com/nathanmartins/sysperf/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

var MemInfoMetrics = collectors.MemInfoCollector{}
var CPUMetrics = collectors.CPUCollector{}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	prometheus.MustRegister(MemInfoMetrics)
	prometheus.MustRegister(CPUMetrics)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msg("running on port :8080")
	log.Err(http.ListenAndServe(":8080", nil))
}
