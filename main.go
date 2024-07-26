package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	kitlog "github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"

	"github.com/frebib/zfs-exporter/collector"
	"github.com/frebib/zfs-exporter/zfs"
)

var (
	listenAddress = flag.String("web.listen-address", ":9254", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	webConfigFile = flag.String("web.config.file", "", "Path to web-config file")
)

func main() {
	flag.Parse()

	libzfs, err := zfs.New()
	if err != nil {
		panic(err)
	}
	defer libzfs.Close()

	opts := promhttp.HandlerOpts{
		ErrorLog:          log.Default(),
		ErrorHandling:     promhttp.PanicOnError,
		EnableOpenMetrics: true,
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collectors.NewGoCollector())
	registry.MustRegister(collector.NewZpoolCollector(libzfs))
	registry.MustRegister(collector.NewDatasetCollector(libzfs))

	router := http.NewServeMux()
	router.Handle(*metricsPath, promhttp.HandlerFor(registry, opts))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html><head><title>ZFS Exporter</title></head>`+
			`<body><h1>ZFS Exporter</h1>`+
			`<p><a href="%s">Metrics</a></p></body></html>`,
			*metricsPath,
		)
	})

	server := http.Server{Handler: router}
	logger := kitlog.NewLogfmtLogger(log.Writer())
	config := web.FlagConfig{
		WebListenAddresses: &[]string{*listenAddress},
		WebConfigFile:      webConfigFile,
	}

	err = web.ListenAndServe(&server, &config, logger)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("%s", err)
	}
}
