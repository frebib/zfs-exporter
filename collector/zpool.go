package collector

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	zfs "github.com/frebib/go-libzfs"
)

var (
	zioTypeNames = []string{
		"null",
		"read",
		"write",
		"free",
		"claim",
		"ioctl",
	}

	vdevOpsDesc = prometheus.NewDesc(
		"zfs_pool_ops_total",
		"number of operations performed.",
		[]string{"pool", "type", "vdev", "vdevid", "op"},
		nil,
	)

	vdevBytesDesc = prometheus.NewDesc(
		"zfs_pool_bytes_total",
		"number of bytes handled",
		[]string{"pool", "type", "vdev", "vdevid", "op"},
		nil,
	)

	vdevErrorsDesc = prometheus.NewDesc(
		"zfs_pool_errors_total",
		"number of errors seen",
		[]string{"pool", "type", "vdev", "vdevid", "errortype"},
		nil,
	)

	vdevStateDesc = prometheus.NewDesc(
		"zfs_pool_vdev_state",
		"vdev state: Unknown, Closed, Offline, Removed, CantOpen, Faulted, Degraded, Healthy.",
		[]string{"pool", "type", "vdev", "vdevid", "state"},
		nil,
	)

	vdevAllocDesc = prometheus.NewDesc(
		"zfs_pool_allocated_bytes",
		"number of bytes allocated (usage)",
		[]string{"pool", "type", "vdev", "vdevid"},
		nil,
	)

	vdevSizeDesc = prometheus.NewDesc(
		"zfs_pool_size_bytes",
		"size of the vdev in bytes (total capacity).",
		[]string{"pool", "type", "vdev", "vdevid"},
		nil,
	)

	vdevFreeDesc = prometheus.NewDesc(
		"zfs_pool_free_bytes",
		"free space on the vdev in bytes.",
		[]string{"pool", "type", "vdev", "vdevid"},
		nil,
	)

	vdevFragDesc = prometheus.NewDesc(
		"zfs_pool_fragmentation_percent",
		"device fragmentation percentage",
		[]string{"pool", "type", "vdev", "vdevid"},
		nil,
	)

	poolStateDesc = prometheus.NewDesc(
		"zfs_pool_state",
		"pool state enum: Active, Exported, Destroyed, Spare, L2cache, uninitialized, unavail, potentiallyactive",
		[]string{"pool", "state"},
		nil,
	)

	poolStatusDesc = prometheus.NewDesc(
		"zfs_pool_status",
		"pool status enum: CorruptCache, MissingDevR, MissingDevNr, CorruptLabelR, CorruptLabelNr, BadGUIDSum, CorruptPool, CorruptData, FailingDev, VersionNewer, HostidMismatch, IoFailureWait, IoFailureContinue, BadLog, Errata, UnsupFeatRead, UnsupFeatWrite, FaultedDevR, FaultedDevNr, VersionOlder, FeatDisabled, Resilvering, OfflineDev, RemovedDev, Ok",
		[]string{"pool", "status"},
		nil,
	)

	poolReadonlyDesc = prometheus.NewDesc(
		"zfs_pool_readonly",
		"Read-only status of the pool [0: read-write, 1: read-only].",
		[]string{"pool"},
		nil,
	)

	poolScrubTimeDesc = prometheus.NewDesc(
		"zfs_pool_last_scrub_timestamp",
		"Unix timestamp of the start of the last scrub",
		[]string{"pool"},
		nil,
	)

	poolCollectErrors = prometheus.NewDesc(
		"zfs_pool_collect_errors_total",
		"errors collecting ZFS metrics",
		[]string{"pool"},
		nil,
	)
)

type ZpoolCollector struct {
	libzfs *zfs.LibZFS

	poolErrors map[string]int
}

// Describe implements prometheus.Collector.
func (collector *ZpoolCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- vdevOpsDesc
	descs <- vdevBytesDesc
	descs <- vdevErrorsDesc
	descs <- vdevStateDesc
	descs <- vdevAllocDesc
	descs <- vdevSizeDesc
	descs <- vdevFreeDesc
	descs <- vdevFragDesc
	descs <- poolStateDesc
	descs <- poolStatusDesc
	descs <- poolReadonlyDesc
	descs <- poolScrubTimeDesc
	descs <- poolCollectErrors
}

func NewZpoolCollector(libzfs *zfs.LibZFS) *ZpoolCollector {
	return &ZpoolCollector{
		libzfs:     libzfs,
		poolErrors: make(map[string]int),
	}
}

// Collect implements prometheus.Collector.
func (collector *ZpoolCollector) Collect(ch chan<- prometheus.Metric) {
	pools, err := collector.libzfs.PoolOpenAll()
	if err != nil {
		log.Printf("error opening pools: %v", err)
		ch <- prometheus.NewInvalidMetric(nil, err)
		return
	}

	for _, pool := range pools {
		collector.collectPool(ch, pool)
		pool.Close()
	}

	runtime.GC()
}

func (collector *ZpoolCollector) collectPool(metrics chan<- prometheus.Metric, pool *zfs.Pool) {
	name := pool.Name()
	if _, ok := collector.poolErrors[name]; !ok {
		collector.poolErrors[name] = 0
	}

	state, err := pool.State()
	if err != nil {
		log.Printf("error getting state of pool '%s': %v\n", name, err)
		collector.poolErrors[name]++
	} else {
		metrics <- prometheus.MustNewConstMetric(
			poolStateDesc,
			prometheus.GaugeValue,
			float64(state),
			name, strings.ToLower(state.String()),
		)
	}

	status, err := pool.Status()
	if err != nil {
		log.Printf("error getting status of pool '%s': %v", name, err)
		collector.poolErrors[name]++
	} else {
		metrics <- prometheus.MustNewConstMetric(
			poolStatusDesc,
			prometheus.GaugeValue,
			float64(status),
			name, strings.ToLower(status.String()),
		)
	}

	roProp, err := pool.Get(zfs.PoolPropReadonly)
	if err != nil {
		log.Printf("error getting property '%s' of pool '%s': %v",
			zfs.PoolPropReadonly, name, err,
		)
		collector.poolErrors[name]++

	} else {
		readonly := 0.0

		if roProp.Value != "on" && roProp.Value != "off" {
			log.Printf("readonly value is unexpected: %s", roProp.Value)
			collector.poolErrors[name]++

		} else {
			if roProp.Value == "on" {
				readonly = 1.0
			}

			metrics <- prometheus.MustNewConstMetric(
				poolReadonlyDesc,
				prometheus.GaugeValue,
				readonly,
				name,
			)
		}
	}

	var vdt zfs.VDevTree
	vdt, err = pool.VDevTree()
	if err != nil {
		log.Printf("unable to read vdevtree for pool '%s': %v", name, err)
		collector.poolErrors[name]++
	} else {
		err = collector.collectVdev(metrics, vdt, name)
		if err != nil {
			log.Printf("unable to read vdevtree stats for pool '%s': %v", name, err)
			collector.poolErrors[name]++
		}
	}

	scan, err := vdt.ScanStat()
	if err != nil {
		if !errors.Is(err, zfs.ErrNotFound) {
			log.Printf("unable to read scan statistics for pool '%s': %v", name, err)
			collector.poolErrors[name]++
		}
	} else {
		if scan.Func == zfs.ScanScrub {
			metrics <- prometheus.MustNewConstMetric(
				poolScrubTimeDesc,
				prometheus.GaugeValue,
				float64(scan.StartTime.Unix()),
				name,
			)
		}
	}

	metrics <- prometheus.MustNewConstMetric(
		poolCollectErrors,
		prometheus.CounterValue,
		float64(collector.poolErrors[name]),
		name,
	)
}

func (collector *ZpoolCollector) collectVdev(ch chan<- prometheus.Metric, vdt zfs.VDevTree, name string) error {
	stat, err := vdt.Stat()
	if err != nil {
		return err
	}

	vName := vdt.Name()
	vType := string(vdt.Type())
	id := fmt.Sprintf("%d", vdt.ID())

	ch <- prometheus.MustNewConstMetric(
		vdevStateDesc,
		prometheus.GaugeValue,
		float64(stat.State),
		name, vType, vName, id,
		strings.ToLower(stat.State.String()),
	)

	if vType != zfs.VDevTypeDisk && vType != zfs.VDevTypeFile {
		ch <- prometheus.MustNewConstMetric(
			vdevAllocDesc,
			prometheus.GaugeValue,
			float64(stat.Alloc),
			name, vType, vName, id,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevSizeDesc,
			prometheus.GaugeValue,
			float64(stat.Space),
			name, vType, vName, id,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevFreeDesc,
			prometheus.GaugeValue,
			float64(stat.Space-stat.Alloc),
			name, vType, vName, id,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevFragDesc,
			prometheus.GaugeValue,
			float64(stat.Fragmentation),
			name, vType, vName, id,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.ReadErrors),
		name, vType, vName, id, "read",
	)
	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.WriteErrors),
		name, vType, vName, id, "write",
	)
	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.ChecksumErrors),
		name, vType, vName, id, "checksum",
	)

	for op := zfs.ZIOTypeNull + 1; op < zfs.ZIOTypes; op++ {
		ch <- prometheus.MustNewConstMetric(
			vdevOpsDesc, prometheus.CounterValue,
			float64(stat.Ops[op]),
			name, vType, vName, id, zioTypeNames[op],
		)

		ch <- prometheus.MustNewConstMetric(
			vdevBytesDesc, prometheus.CounterValue,
			float64(stat.Bytes[op]),
			name, vType, vName, id, zioTypeNames[op],
		)
	}

	// recurse
	for _, child := range vdt.Children() {
		err := collector.collectVdev(ch, child, name)
		if err != nil {
			return err
		}
	}

	return nil
}
