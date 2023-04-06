package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"postgres_logical_replication_exporter/collector"
	"postgres_logical_replication_exporter/pg"
	"time"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	primaryuri string
	standbyuri string
	addr       string
	loglevel   string
)

func init() {
	flag.StringVar(&primaryuri, "primary-uri", "", "Connection URI of the primary instance host.")
	flag.StringVar(&standbyuri, "standby-uri", "", "Connection URI of the standby instance host.")
	flag.StringVar(&addr, "listen-address", ":9394", "The address to listen on for HTTP requests.")
	flag.StringVar(&loglevel, "log-level", "info", "Level of the logs.")
}

func main() {
	flag.Parse()

	logger, err := newLogger(loglevel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialize logger: %v\n", err)
		os.Exit(1)
	}

	primary := pg.MustConnect(primaryuri, 5)
	standby := pg.MustConnect(standbyuri, 5)

	mux := http.NewServeMux()
	reg := prometheus.NewRegistry()

	reg.MustRegister(collector.NewCollector(primary, standby, logger))

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>PG Logical Replication Exporter</title></head>
			<body>
			<h1>PG Logical Replication Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	}))

	mux.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
			ErrorLog:          newPromLogger(logger),
		},
	))

	mux.Handle("/healthcheck", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = level.Debug(logger).Log("msg", "/healthcheck OK")
	}))

	mux.Handle("/readyz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = level.Debug(logger).Log("msg", "/readyz OK")
	}))

	server := &http.Server{
		Addr:              addr,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           mux,
	}

	_ = level.Info(logger).Log("msg", "Starting http server", "address", addr)
	_ = level.Error(logger).Log("msg", "Error starting HTTP server", "err", server.ListenAndServe())
}
