package main

import (
	"github.com/nathanmartins/sysperf/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var memCollector collectors.MeminfoCollector
	if err := prometheus.Register(&memCollector); err != nil {
		log.Err(err)
	} else {
		log.Info().Msg("memcollector registered.")
	}
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	log.Info().Msg("Beginning to serve on port :8080")
	log.Err(http.ListenAndServe(":8080", nil))
}
