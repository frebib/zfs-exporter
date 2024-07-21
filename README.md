# ZFS Prometheus exporter

Exports dataset, pool and disk/vdev statistics in a Prometheus-compatible format.
It works (as tested) on Linux and OpenBSD, plus probably others too.

`go-libzfs` is based around github.com/bicomsystems/go-libzfs with a 
considerable portion of the libzfs wrapping logic rewritten. Credit goes to that
implementation for the framework and pointers on how to reference the largely
undocumented (though the source itself was very helpful) libzfs.

Provided is an example Grafana dashboard that looks a little something like this
![grafana](https://github.com/frebib/zfs-exporter/blob/master/contrib/grafana.png?raw=true)

## Usage

The exporter takes very few options. It supports (m)TLS/authentication via [exporter-toolkit](https://github.com/prometheus/exporter-toolkit/tree/master/web)

```
$ zfs-exporter --help
  -web.config.file string
    	Path to web-config file
  -web.listen-address string
    	Address on which to expose metrics and web interface. (default ":9254")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
```
