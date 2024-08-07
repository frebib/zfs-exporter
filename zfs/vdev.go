package zfs

/*
#include <libzfs.h>
*/
import "C"

import (
	"errors"
	"time"
	"unsafe"
)

// VDevType type of device in the pool
type VDevType string

// Types of Virtual Devices
const (
	VDevTypeRoot      VDevType = "root"      // Root device in ZFS pool
	VDevTypeMirror             = "mirror"    // Mirror device in ZFS pool
	VDevTypeReplacing          = "replacing" // Replacing
	VDevTypeRaidz              = "raidz"     // RAIDZ device
	VDevTypeDisk               = "disk"      // Device is disk
	VDevTypeFile               = "file"      // Device is file
	VDevTypeMissing            = "missing"   // Missing device
	VDevTypeHole               = "hole"      // Hole
	VDevTypeSpare              = "spare"     // Spare device
	VDevTypeLog                = "log"       // ZIL device
	VDevTypeL2cache            = "l2cache"   // Cache device (disk)
)

// VDevState values are ordered from least to most healthy.
// Less than or equal to VDevStateCantOpen is considered unusable.
type VDevState uint64

const (
	VDevStateUnknown  VDevState = iota // Uninitialized vdev
	VDevStateClosed                    // Not currently open
	VDevStateOffline                   // Not allowed to open
	VDevStateRemoved                   // Explicitly removed from system
	VDevStateCantOpen                  // Tried to open, but failed
	VDevStateFaulted                   // External request to fault device
	VDevStateDegraded                  // Replicated vdev with unhealthy kids
	VDevStateHealthy                   // Presumed good
)

func (s VDevState) String() string {
	switch s {
	case VDevStateUnknown:
		return "uninitialized"
	case VDevStateClosed:
		return "closed"
	case VDevStateOffline:
		return "offline"
	case VDevStateRemoved:
		return "removed"
	case VDevStateCantOpen:
		return "cantopen"
	case VDevStateFaulted:
		return "faulted"
	case VDevStateDegraded:
		return "degraded"
	case VDevStateHealthy:
		return "online"
	default:
		return "unknown"
	}
}

// VDevAux - vdev aux states
type VDevAux uint64

// vdev aux states.  When a vdev is in the VDevStateCantOpen state, the aux field
// of the vdev stats structure uses these constants to distinguish why.
const (
	VDevAuxNone         VDevAux = iota // No error
	VDevAuxOpenFailed                  // Ldi_open_*() or vn_open() failed
	VDevAuxCorruptData                 // Bad label or disk contents
	VDevAuxNoReplicas                  // Insufficient number of replicas
	VDevAuxBadGUIDSum                  // Vdev guid sum doesn't match
	VDevAuxTooSmall                    // Vdev size is too small
	VDevAuxBadLabel                    // The label is OK but invalid
	VDevAuxVersionNewer                // On-disk version is too new
	VDevAuxVersionOlder                // On-disk version is too old
	VDevAuxUnsupFeat                   // Unsupported features
	VDevAuxSpared                      // Hot spare used in another pool
	VDevAuxErrExceeded                 // Too many errors
	VDevAuxIOFailure                   // Experienced I/O failure
	VDevAuxBadLog                      // Cannot read log chain(s)
	VDevAuxExternal                    // External diagnosis
	VDevAuxSplitPool                   // Vdev was split off into another pool
)

// VDevTree ZFS virtual device tree
type VDevTree struct {
	pool *Pool
	nvl  NVList

	// cached after first use
	name string
}

func (vdt VDevTree) Config() NVList {
	return vdt.nvl
}

func (vdt VDevTree) Name() string {
	if vdt.name != "" {
		return vdt.name
	}

	/*
		Pools with disks that are unavailable will return the guid as the name
		but also show the previous name to the right. We want that one.

		https://github.com/openzfs/zfs/blob/zfs-2.0.5/cmd/zpool/zpool_main.c#L2255-L2258

		NAME                     STATE     READ WRITE CKSUM
		test                     DEGRADED     0     0     0
		  raidz1-0               DEGRADED     0     0     0
			/tmp/disk1           ONLINE       0     0     0
			/tmp/disk2           ONLINE       0     0     0
			/tmp/disk3           ONLINE       0     0     0
			9445367449878001541  UNAVAIL      0     0     0  was /tmp/disk4
	*/
	_, err := vdt.nvl.LookupUint64(PoolConfigNotPresent)
	if err == nil {
		vdt.name, err = vdt.nvl.LookupString(PoolConfigPath)
		// Must succeed: https://github.com/openzfs/zfs/blob/zfs-2.0.5/cmd/zpool/zpool_main.c#L2257
		if err != nil {
			panic(err)
		}
		return vdt.name
	}

	libzfs := vdt.pool.LibZFS().Handle()
	// Flag VDEV_NAME_TYPE_ID gives `raidz1-n` where n is the vdev-id. Without
	// it (flag value of zero), we'd simply get `raidz1` which isn't unique.
	ptr := C.zpool_vdev_name(libzfs, vdt.pool.handle, vdt.nvl.handle, C.VDEV_NAME_TYPE_ID)
	if ptr == nil {
		panic("zpool_vdev_name() returned nil")
	}
	defer C.free(unsafe.Pointer(ptr))

	// cache for future calls
	vdt.name = C.GoString(ptr)
	return vdt.name
}

func (vdt VDevTree) Path() string {
	switch vdt.Type() {
	case VDevTypeDisk:
	case VDevTypeFile:
		break
	default:
		return ""
	}

	path, err := vdt.nvl.LookupString(PoolConfigPath)
	if err != nil {
		panic(err)
	}
	return path
}

func (vdt VDevTree) Type() VDevType {
	typ, err := vdt.nvl.LookupString(PoolConfigType)
	if err != nil {
		panic(err)
	}
	return VDevType(typ)
}

func (vdt VDevTree) ID() uint64 {
	id, err := vdt.nvl.LookupUint64(PoolConfigID)
	if err != nil {
		panic(err)
	}
	return id
}

func (vdt VDevTree) GUID() uint64 {
	guid, err := vdt.nvl.LookupUint64(PoolConfigGUID)
	if err != nil {
		panic(err)
	}
	return guid
}

func (vdt VDevTree) Children() []VDevTree {
	nvls, err := vdt.nvl.LookupNVListArray(PoolConfigChildren)
	if errors.Is(err, ErrNotFound) {
		return []VDevTree{}
	} else if err != nil {
		panic(err)
	}

	children := make([]VDevTree, len(nvls))
	for i, nvl := range nvls {
		children[i].pool = vdt.pool
		children[i].nvl = nvl
	}

	return children
}

func (vdt VDevTree) ScanStat() (PoolScanStat, error) {
	var ss *C.pool_scan_stat_t
	var count C.uint_t

	scanStats := C.CString(PoolConfigScanStats)
	defer C.free(unsafe.Pointer(scanStats))

	// Here we "cheat" by unloading the uint64_t array into a scan_stat_t struct
	// as the fields are all uint64_t and in the correct order for us already
	ret := C.nvlist_lookup_uint64_array(vdt.nvl.Pointer(), scanStats,
		(**C.uint64_t)(unsafe.Pointer(&ss)), &count)
	if ret != 0 {
		return PoolScanStat{}, nvlistLookupError(ret)
	}

	stat := PoolScanStat{
		Func:           ScanFunc(ss.pss_func),
		State:          ScanState(ss.pss_state),
		StartTime:      time.Unix(int64(ss.pss_start_time), 0).UTC(),
		EndTime:        time.Unix(int64(ss.pss_end_time), 0).UTC(),
		ToExamine:      uint64(ss.pss_to_examine),
		Examined:       uint64(ss.pss_examined),
		Skipped:        uint64(ss.pss_skipped),
		Processed:      uint64(ss.pss_processed),
		Errors:         uint64(ss.pss_errors),
		PassExam:       uint64(ss.pss_pass_exam),
		PassStart:      time.Unix(int64(ss.pss_pass_start), 0).UTC(),
		PassScrubPause: time.Unix(int64(ss.pss_pass_scrub_pause), 0).UTC(),
	}

	return stat, nil
}

func (vdt VDevTree) Stat() (VDevStat, error) {
	var vs *C.struct_vdev_stat
	var count C.uint_t

	vdevStats := C.CString(PoolConfigVdevStats)
	defer C.free(unsafe.Pointer(vdevStats))

	// Here we "cheat" by unloading the uint64_t array into a vdev_stat_t struct
	// as the fields are all uint64_t and in the correct order for us already
	ret := C.nvlist_lookup_uint64_array(vdt.nvl.Pointer(), vdevStats,
		(**C.uint64_t)(unsafe.Pointer(&vs)), &count)
	if ret != 0 {
		return VDevStat{}, nvlistLookupError(ret)
	}

	stat := VDevStat{
		Timestamp:      time.Unix(int64(vs.vs_timestamp), 0).UTC(),
		State:          VDevState(vs.vs_state),
		Aux:            VDevAux(vs.vs_aux),
		Alloc:          uint64(vs.vs_alloc),
		Space:          uint64(vs.vs_space),
		DSpace:         uint64(vs.vs_dspace),
		RSize:          uint64(vs.vs_rsize),
		ESize:          uint64(vs.vs_esize),
		ReadErrors:     uint64(vs.vs_read_errors),
		WriteErrors:    uint64(vs.vs_write_errors),
		ChecksumErrors: uint64(vs.vs_checksum_errors),
		SelfHealed:     uint64(vs.vs_self_healed),
		ScanRemoving:   uint64(vs.vs_scan_removing),
		ScanProcessed:  uint64(vs.vs_scan_processed),
		Fragmentation:  uint64(vs.vs_fragmentation),
	}

	for z := 0; z < ZIOTypes; z++ {
		stat.Ops[z] = uint64(vs.vs_ops[z])
		stat.Bytes[z] = uint64(vs.vs_bytes[z])
	}

	return stat, nil
}

// VDevStat are vdev statistics. All fields should be 64-bit because this is
// passed between kernel and userland as an nvlist uint64 array.
type VDevStat struct {
	Timestamp      time.Time        // Time since vdev load
	State          VDevState        // vdev state
	Aux            VDevAux          // See vdev_aux_t
	Alloc          uint64           // Space allocated
	Space          uint64           // Total capacity
	DSpace         uint64           // Deflated capacity
	RSize          uint64           // Replaceable dev size
	ESize          uint64           // Expandable dev size
	Ops            [ZIOTypes]uint64 // Operation count
	Bytes          [ZIOTypes]uint64 // Bytes read/written
	ReadErrors     uint64           // Read errors
	WriteErrors    uint64           // Write errors
	ChecksumErrors uint64           // Checksum errors
	SelfHealed     uint64           // Self-healed bytes
	ScanRemoving   uint64           // Removing?
	ScanProcessed  uint64           // Scan processed bytes
	Fragmentation  uint64           // Device fragmentation
}
