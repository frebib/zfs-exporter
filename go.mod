module github.com/frebib/zfs-exporter

go 1.16

require (
	github.com/frebib/go-libzfs v0.0.0
	github.com/prometheus/client_golang v1.11.0
)

replace github.com/frebib/go-libzfs v0.0.0 => ./go-libzfs
