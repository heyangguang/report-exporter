package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"report-exporter/collector"
)

func main() {
	newCheckFile := collector.NewCheckFileCollector()
	prometheus.MustRegister(newCheckFile)
	log.Info("binding to server on Port : 9100")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:9100", nil))
}
