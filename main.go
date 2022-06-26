package main

import (
	"fmt"
	"github.com/nathanmartins/sysperf/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

func init() {
	var memCollector collectors.MeminfoCollector

	// This will register fine.
	if err := prometheus.Register(&memCollector); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("memcollector registered.")
	}
}

func main() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
