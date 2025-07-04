package collector

import (
	"errors"
	"log"
	"runtime"
	"strings"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/frebib/zfs-exporter/zfs"
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
		[]string{"pool", "type", "parent", "device", "path", "op"},
		nil,
	)

	vdevBytesDesc = prometheus.NewDesc(
		"zfs_pool_bytes_total",
		"number of bytes handled",
		[]string{"pool", "type", "parent", "device", "path", "op"},
		nil,
	)

	vdevErrorsDesc = prometheus.NewDesc(
		"zfs_pool_errors_total",
		"number of errors seen",
		[]string{"pool", "type", "parent", "device", "path", "errortype"},
		nil,
	)

	vdevStateDesc = prometheus.NewDesc(
		"zfs_pool_vdev_state",
		"vdev state: Unknown, Closed, Offline, Removed, CantOpen, Faulted, Degraded, Healthy.",
		[]string{"pool", "type", "parent", "device", "path", "state"},
		nil,
	)

	vdevAllocDesc = prometheus.NewDesc(
		"zfs_pool_allocated_bytes",
		"number of bytes allocated (usage)",
		[]string{"pool", "type", "parent", "device", "path"},
		nil,
	)

	vdevSizeDesc = prometheus.NewDesc(
		"zfs_pool_size_bytes",
		"size of the vdev in bytes (total capacity).",
		[]string{"pool", "type", "parent", "device", "path"},
		nil,
	)

	vdevFreeDesc = prometheus.NewDesc(
		"zfs_pool_free_bytes",
		"free space on the vdev in bytes.",
		[]string{"pool", "type", "parent", "device", "path"},
		nil,
	)

	vdevFragDesc = prometheus.NewDesc(
		"zfs_pool_fragmentation_percent",
		"device fragmentation percentage",
		[]string{"pool", "type", "parent", "device", "path"},
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

	poolScrubStatus = prometheus.NewDesc(
		"zfs_pool_scrub_status",
		"Scrub status [0: inactive, 1: scanning, 2:finished, 3: cancelled]",
		[]string{"pool"},
		nil,
	)
	poolScrubPauseTime = prometheus.NewDesc(
		"zfs_pool_scrub_paused_timestamp",
		"Unix timestamp of when the scrub was paused. Zero indicates that the scrub is not paused",
		[]string{"pool"},
		nil,
	)
	poolScrubStartTimeDesc = prometheus.NewDesc(
		"zfs_pool_last_scrub_start_timestamp",
		"Unix timestamp of the start of the last scrub",
		[]string{"pool"},
		nil,
	)
	poolScrubEndTimeDesc = prometheus.NewDesc(
		"zfs_pool_last_scrub_end_timestamp",
		"Unix timestamp of the end of the last scrub",
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
	descs <- poolScrubStatus
	descs <- poolScrubStartTimeDesc
	descs <- poolScrubEndTimeDesc
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

	state := pool.State()
	metrics <- prometheus.MustNewConstMetric(
		poolStateDesc,
		prometheus.GaugeValue,
		float64(state),
		name, strings.ToLower(state.String()),
	)

	status := pool.Status()
	metrics <- prometheus.MustNewConstMetric(
		poolStatusDesc,
		prometheus.GaugeValue,
		float64(status),
		name, strings.ToLower(status.String()),
	)

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
		// Pass empty "parent" because pools are top-level. Label will be empty
		// and appear absent in Prometheus.
		err = collector.collectVdev(metrics, vdt, name, "")
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
				poolScrubStatus,
				prometheus.GaugeValue,
				float64(scan.State),
				name,
			)
			metrics <- prometheus.MustNewConstMetric(
				poolScrubPauseTime,
				prometheus.GaugeValue,
				float64(scan.PassScrubPause.Unix()),
				name,
			)
			metrics <- prometheus.MustNewConstMetric(
				poolScrubStartTimeDesc,
				prometheus.GaugeValue,
				float64(scan.StartTime.Unix()),
				name,
			)
			metrics <- prometheus.MustNewConstMetric(
				poolScrubEndTimeDesc,
				prometheus.GaugeValue,
				float64(scan.EndTime.Unix()),
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

func (collector *ZpoolCollector) collectVdev(ch chan<- prometheus.Metric, vdt zfs.VDevTree, pool, parent string) error {
	stat, err := vdt.Stat()
	if err != nil {
		return err
	}

	name := vdt.Name()
	devType := vdt.Type()
	path := vdt.Path()

	// Try to resolve the root disk in /dev for the vdev disk (partition)
	if path != "" && devType == zfs.VDevTypeDisk {
		disk, err := diskFromPartition(path)
		if err != nil {
			log.Printf("error resolving disk path '%s': %s", path, err)
		} else {
			path = disk
		}
	}

	isLog, err := vdt.Config().LookupUint64(zfs.PoolConfigIsLog)
	if err != nil && !errors.Is(err, zfs.ErrNotFound) {
		panic(err)
	} else if isLog > 0 {
		// Falsify the "log" device type for log disks
		devType = zfs.VDevTypeLog
	}

	typ := string(devType)

	ch <- prometheus.MustNewConstMetric(
		vdevStateDesc,
		prometheus.GaugeValue,
		float64(stat.State),
		pool, typ, parent, name, path,
		strings.ToLower(stat.State.String()),
	)

	if devType != zfs.VDevTypeDisk && devType != zfs.VDevTypeFile {
		ch <- prometheus.MustNewConstMetric(
			vdevAllocDesc,
			prometheus.GaugeValue,
			float64(stat.Alloc),
			pool, typ, parent, name, path,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevSizeDesc,
			prometheus.GaugeValue,
			float64(stat.Space),
			pool, typ, parent, name, path,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevFreeDesc,
			prometheus.GaugeValue,
			float64(stat.Space-stat.Alloc),
			pool, typ, parent, name, path,
		)
		ch <- prometheus.MustNewConstMetric(
			vdevFragDesc,
			prometheus.GaugeValue,
			float64(stat.Fragmentation),
			pool, typ, parent, name, path,
		)
	}

	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.ReadErrors),
		pool, typ, parent, name, path, "read",
	)
	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.WriteErrors),
		pool, typ, parent, name, path, "write",
	)
	ch <- prometheus.MustNewConstMetric(
		vdevErrorsDesc,
		prometheus.CounterValue,
		float64(stat.ChecksumErrors),
		pool, typ, parent, name, path, "checksum",
	)

	for op := zfs.ZIOTypeNull + 1; op < zfs.ZIOTypes; op++ {
		ch <- prometheus.MustNewConstMetric(
			vdevOpsDesc, prometheus.CounterValue,
			float64(stat.Ops[op]),
			pool, typ, parent, name, path, zioTypeNames[op],
		)

		ch <- prometheus.MustNewConstMetric(
			vdevBytesDesc, prometheus.CounterValue,
			float64(stat.Bytes[op]),
			pool, typ, parent, name, path, zioTypeNames[op],
		)
	}

	// recurse
	for _, child := range vdt.Children() {
		err := collector.collectVdev(ch, child, pool, name)
		if err != nil {
			return err
		}
	}

	return nil
}
