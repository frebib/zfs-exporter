package collector

import (
	"log"
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/frebib/zfs-exporter/zfs"
)

var (
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
	datasets, err := collector.libzfs.DatasetOpenAll(zfs.DatasetTypeFilesystem, -1)

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

	descs := map[zfs.DatasetProperty]*prometheus.Desc{
		zfs.DatasetPropUsed:          datasetUsedBytes,
		zfs.DatasetPropReferenced:    datasetRefBytes,
		zfs.DatasetPropAvailable:     datasetAvailBytes,
		zfs.DatasetPropWritten:       datasetWrittenBytes,
		zfs.DatasetPropReadonly:      datasetReadOnly,
		zfs.DatasetPropCompressratio: datasetCompressRatio,
	}

	switch typ {
	case zfs.DatasetTypeVolume:
		// Only applicable to volumes
		descs[zfs.DatasetPropVolsize] = datasetVolsizeBytes
	default:
		descs[zfs.DatasetPropQuota] = datasetQuotaBytes
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

		propType:
			switch prop.Type() {
			case zfs.PropertyTypeNumber:
				value, err = strconv.ParseFloat(val.Value, 10)
				if err != nil {
					log.Printf("error parsing property '%s' value '%s': %v", prop, val.Value, err)
					collector.datasetErrors[name]++
					continue
				}

			case zfs.PropertyTypeIndex:
				switch prop {
				case zfs.DatasetPropReadonly:
					if val.Value == "on" {
						value = 1
					} else {
						value = 0
					}
					break propType
				}

				fallthrough

			default:
				collector.datasetErrors[name]++
				log.Printf("unknown property value for '%s'", prop)
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
