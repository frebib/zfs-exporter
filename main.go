package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	zfs "github.com/frebib/go-libzfs"

	"github.com/frebib/zfs-exporter/collector"
)

var (
	listenAddress = flag.String("web.listen-address", ":9254", "Address on which to expose metrics and web interface.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
)

func main() {
	flag.Parse()

	libzfs, err := zfs.New()
	if err != nil {
		panic(err)
	}
	defer libzfs.Close()
	defer libzfs.Close()

	opts := promhttp.HandlerOpts{
		ErrorLog:          log.Default(),
		ErrorHandling:     promhttp.PanicOnError,
		EnableOpenMetrics: true,
	}

	registry := prometheus.NewRegistry()
	registry.MustRegister(collector.NewZpoolCollector(libzfs))
	registry.MustRegister(collector.NewDatasetCollector(libzfs))

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, opts))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<html><head><title>ZFS Exporter</title></head>`+
			`<body><h1>ZFS Exporter</h1>`+
			`<p><a href="%s">Metrics</a></p></body></html>`,
			*metricsPath,
		)
	})

	log.Printf("Listening on %s", *listenAddress)
	err = http.ListenAndServe(*listenAddress, nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("%s", err)
	}
}
