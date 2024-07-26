package collector

import (
	"log"
	"strconv"
	"sync"

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
)

type DatasetCollector struct {
	libzfs *zfs.libZFS

	datasetErrors *prometheus.CounterVec
}

// Describe implements prometheus.Collector.
func (collector *DatasetCollector) Describe(descs chan<- *prometheus.Desc) {
	collector.datasetErrors.Describe(descs)
}

func NewDatasetCollector(libzfs *zfs.libZFS) *DatasetCollector {
	return &DatasetCollector{
		libzfs: libzfs,
		datasetErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: "zfs",
				Subsystem: "dataset",
				Name:      "collect_errors_total",
				Help:      "errors collecting ZFS dataset metrics",
			},
			[]string{"dataset"},
		),
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

	var wg sync.WaitGroup
	wg.Add(len(datasets))
	for _, dataset := range datasets {
		go func(d *zfs.Dataset) {
			collector.collectDataset(ch, d)
			d.Close()
			wg.Done()
		}(dataset)
	}
	wg.Wait()

	collector.datasetErrors.Collect(ch)
}

func (collector *DatasetCollector) collectDataset(metrics chan<- prometheus.Metric, dataset *zfs.Dataset) {
	name := dataset.Name()
	pool := dataset.Pool().Name()
	typ := dataset.Type()

	// Initialise error counter as 0
	errs := collector.datasetErrors.With(prometheus.Labels{"dataset": name})

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
		errs.Inc()
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
					errs.Inc()
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
				errs.Inc()
				log.Printf("unknown property value for '%s'", prop)
				continue
			}

			metrics <- prometheus.MustNewConstMetric(
				descs[prop], prometheus.GaugeValue,
				value, name, pool, typ.String(),
			)
		}
	}
}
