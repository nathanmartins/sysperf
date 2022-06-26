package main

import (
	"github.com/nathanmartins/sysperf/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

var MemInfoMetric = collectors.MemInfoCollector{}

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	prometheus.MustRegister(MemInfoMetric)
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msg("running on port :8080")
	log.Err(http.ListenAndServe(":8080", nil))
}
