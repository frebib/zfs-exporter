package collector

import (
	"log"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/frebib/zfs-exporter/zfs"
)

var (
	datasetCreatedAt = prometheus.NewDesc(
		"zfs_dataset_created_timestamp_seconds",
		"Unix timestamp representing the created date/time of the dataset",
		[]string{"name", "pool", "type"},
		nil,
	)
	datasetUsedBytes = prometheus.NewDesc(
		"zfs_dataset_used_bytes",
		"space used by dataset and all its descendents in bytes",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetRefBytes = prometheus.NewDesc(
		"zfs_dataset_referenced_bytes",
		"",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetAvailBytes = prometheus.NewDesc(
		"zfs_dataset_available_bytes",
		"space available in the dataset in bytes",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetWrittenBytes = prometheus.NewDesc(
		"zfs_dataset_written_bytes",
		"",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetQuotaBytes = prometheus.NewDesc(
		"zfs_dataset_quota_bytes",
		"",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetVolsizeBytes = prometheus.NewDesc(
		"zfs_volume_size_bytes",
		"size in bytes of a zfs volume",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetReadOnly = prometheus.NewDesc(
		"zfs_dataset_readonly",
		"",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetCompressRatio = prometheus.NewDesc(
		"zfs_dataset_compress_ratio",
		"",
		[]string{"name", "pool", "type"},
		nil,
	)

	datasetCollectErrors = prometheus.NewDesc(
		"zfs_dataset_collect_errors_total",
		"errors collecting ZFS dataset metrics",
		[]string{"dataset"},
		nil,
	)
)

type DatasetCollector struct {
	libzfs *zfs.LibZFS

	datasetErrors map[string]int
}

// Describe implements prometheus.Collector.
func (collector *DatasetCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- datasetCollectErrors
}

func NewDatasetCollector(libzfs *zfs.LibZFS) *DatasetCollector {
	return &DatasetCollector{
		libzfs:        libzfs,
		datasetErrors: make(map[string]int),
	}
}

// Collect implements prometheus.Collector.
func (collector *DatasetCollector) Collect(ch chan<- prometheus.Metric) {
	datasets, err := collector.libzfs.DatasetOpenAll(zfs.DatasetTypeFilesystem|zfs.DatasetTypeSnapshot, -1)

	if err != nil {
		log.Printf("error opening datasets: %v", err)
		//ch <- prometheus.NewInvalidMetric(nil, err)
		return
	}

	for _, dataset := range datasets {
		collector.collectDataset(ch, dataset)
		dataset.Close()
	}

	runtime.GC()
}

func (collector *DatasetCollector) collectDataset(metrics chan<- prometheus.Metric, dataset *zfs.Dataset) {
	name := dataset.Name()
	pool := dataset.Pool().Name()
	typ := dataset.Type()

	if _, ok := collector.datasetErrors[name]; !ok {
		collector.datasetErrors[name] = 0
	}

	// Common properties for all dataset types (filesystems, volumes, snapshots)
	descs := map[zfs.DatasetProperty]*prometheus.Desc{
		zfs.DatasetPropCreation:   datasetCreatedAt,
		zfs.DatasetPropUsed:       datasetUsedBytes,
		zfs.DatasetPropReferenced: datasetRefBytes,
		zfs.DatasetPropWritten:    datasetWrittenBytes,
	}

	// Common properties to filesystem & volume datasets
	if typ == zfs.DatasetTypeFilesystem || typ == zfs.DatasetTypeVolume {
		descs[zfs.DatasetPropAvailable] = datasetAvailBytes
		descs[zfs.DatasetPropCompressratio] = datasetCompressRatio
		descs[zfs.DatasetPropReadonly] = datasetReadOnly
	}

	// Various random properties specific to each dataset types
	switch typ {
	case zfs.DatasetTypeFilesystem:
		descs[zfs.DatasetPropQuota] = datasetQuotaBytes
	case zfs.DatasetTypeVolume:
		descs[zfs.DatasetPropVolsize] = datasetVolsizeBytes
	}

	props := make([]zfs.DatasetProperty, 0, len(descs))
	for prop := range descs {
		props = append(props, prop)
	}

	vals, err := dataset.Gets(props...)
	if err != nil {
		collector.datasetErrors[name]++
		log.Printf("error reading dataset properties: %v", err)

	} else {
		for prop, val := range vals {
			var value float64

			switch v := val.(type) {
			case *zfs.DatasetPropertyNumber:
				value = float64(v.Value())
			case *zfs.DatasetPropertyIndex:
				value = float64(v.Value())

			default:
				collector.datasetErrors[name]++
				log.Printf("unknown property type for '%s'", prop)
				continue
			}

			metrics <- prometheus.MustNewConstMetric(
				descs[prop], prometheus.GaugeValue,
				value, name, pool, typ.String(),
			)
		}
	}

	metrics <- prometheus.MustNewConstMetric(
		datasetCollectErrors,
		prometheus.CounterValue,
		float64(collector.datasetErrors[name]),
		name,
	)
}
