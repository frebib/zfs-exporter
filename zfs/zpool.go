package zfs

/*
#include "list.h"
#include <stdlib.h>
#include <libzfs.h>
#include <zfs_prop.h>
*/
import "C"

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"unsafe"
)

type PoolProperty int

// Pool properties. Enumerates available ZFS pool properties. Use it to access
// pool properties either to read or set specific property.
const (
	PoolPropCont PoolProperty = iota - 2
	PoolPropInval
	PoolPropName
	PoolPropSize
	PoolPropCapacity
	PoolPropAltroot
	PoolPropHealth
	PoolPropGUID
	PoolPropVersion
	PoolPropBootfs
	PoolPropDelegation
	PoolPropAutoreplace
	PoolPropCachefile
	PoolPropFailuremode
	PoolPropListsnaps
	PoolPropAutoexpand
	PoolPropDedupditto
	PoolPropDedupratio
	PoolPropFree
	PoolPropAllocated
	PoolPropReadonly
	PoolPropAshift
	PoolPropComment
	PoolPropExpandsz
	PoolPropFreeing
	PoolPropFragmentaion
	PoolPropLeaked
	PoolPropMaxBlockSize
	PoolPropTName
	PoolPropMaxNodeSize
	PoolPropMultiHost
	PoolPropCheckpoint
	PoolPropLoadGUID
	PoolPropAutotrim
	PoolNumProps
)

func (pp PoolProperty) String() string {
	ptr := C.zpool_prop_to_name(C.zpool_prop_t(pp))
	return C.GoString(ptr)
}

func (pp PoolProperty) Type() PropertyType {
	return PropertyType(C.zpool_prop_get_type(C.zpool_prop_t(pp)))
}

// Enable or disable pool feature with this constants
const (
	FENABLED  = "enabled"
	FDISABLED = "disabled"
)

/*
 * ZIO types.  Needed to interpret vdev statistics below.
 */
const (
	ZIOTypeNull = iota
	ZIOTypeRead
	ZIOTypeWrite
	ZIOTypeFree
	ZIOTypeClaim
	ZIOTypeIOCtl
	ZIOTypes
)

type ScanState int

// Scan states
const (
	DSSNone     ScanState = iota // No scan
	DSSScanning                  // Scanning
	DSSFinished                  // Scan finished
	DSSCanceled                  // Scan canceled

	DSSNumStates // Total number of scan states
)

func (ss ScanState) String() string {
	switch ss {
	case DSSNone:
		return "none"
	case DSSScanning:
		return "scanning"
	case DSSFinished:
		return "finished"
	case DSSCanceled:
		return "cancelled"
	default:
		return "unknown"
	}
}

type ScanFunc int

// Scan functions
const (
	ScanNone     ScanFunc = iota // No scan function
	ScanScrub                    // Pools is checked against errors
	ScanResilver                 // Pool is resilvering
	ScanNumFuncs                 // Number of scan functions
)

func (sf ScanFunc) String() string {
	switch sf {
	case ScanNone:
		return "none"
	case ScanScrub:
		return "scrub"
	case ScanResilver:
		return "resilver"
	default:
		return "unknown"
	}
}

// PoolInitializeAction type representing pool initialize action
type PoolInitializeAction int

// Initialize actions
const (
	PoolInitializeStart   PoolInitializeAction = iota // start initialization
	PoolInitializeCancel                              // cancel initialization
	PoolInitializeSuspend                             // suspend initialization
)

func (s PoolInitializeAction) String() string {
	switch s {
	case PoolInitializeStart:
		return "start"
	case PoolInitializeCancel:
		return "cancel"
	case PoolInitializeSuspend:
		return "suspend"
	default:
		return "unknown"
	}
}

// PoolStatus type representing status of the pool
//
//go:generate stringer -type PoolStatus -trimprefix PoolStatus
type PoolStatus int

// Pool status
// https://github.com/openzfs/zfs/blob/e5e76bd6432de9592c4b4319fa826ad39971abd7/include/libzfs.h#L339-L405
const (
	/*
	 * The following correspond to faults as defined in the (fault.fs.zfs.*)
	 * event namespace.  Each is associated with a corresponding message ID.
	 */
	PoolStatusCorruptCache      PoolStatus = iota /* corrupt /kernel/drv/zpool.cache */
	PoolStatusMissingDevR                         /* missing device with replicas */
	PoolStatusMissingDevNr                        /* missing device with no replicas */
	PoolStatusCorruptLabelR                       /* bad device label with replicas */
	PoolStatusCorruptLabelNr                      /* bad device label with no replicas */
	PoolStatusBadGUIDSum                          /* sum of device guids didn't match */
	PoolStatusCorruptPool                         /* pool metadata is corrupted */
	PoolStatusCorruptData                         /* data errors in user (meta)data */
	PoolStatusFailingDev                          /* device experiencing errors */
	PoolStatusVersionNewer                        /* newer on-disk version */
	PoolStatusHostidMismatch                      /* last accessed by another system */
	PoolStatusHosidActive                         /* currently active on another system */
	PoolStatusHostidRequired                      /* multihost=on and hostid=0 */
	PoolStatusIoFailureWait                       /* failed I/O, failmode 'wait' */
	PoolStatusIoFailureContinue                   /* failed I/O, failmode 'continue' */
	PoolStatusIOFailureMMP                        /* ailed MMP, failmode not 'panic' */
	PoolStatusBadLog                              /* cannot read log chain(s) */
	PoolStatusErrata                              /* informational errata available */

	/*
	 * If the pool has unsupported features but can still be opened in
	 * read-only mode, its status is ZPOOL_STATUS_UNSUP_FEAT_WRITE. If the
	 * pool has unsupported features but cannot be opened at all, its
	 * status is ZPOOL_STATUS_UNSUP_FEAT_READ.
	 */
	PoolStatusUnsupFeatRead  /* unsupported features for read */
	PoolStatusUnsupFeatWrite /* unsupported features for write */

	/*
	 * These faults have no corresponding message ID.  At the time we are
	 * checking the status, the original reason for the FMA fault (I/O or
	 * checksum errors) has been lost.
	 */
	PoolStatusFaultedDevR  /* faulted device with replicas */
	PoolStatusFaultedDevNr /* faulted device with no replicas */

	/*
	 * The following are not faults per se, but still an error possibly
	 * requiring administrative attention.  There is no corresponding
	 * message ID.
	 */
	PoolStatusVersionOlder    /* older legacy on-disk version */
	PoolStatusFeatDisabled    /* supported features are disabled */
	PoolStatusResilvering     /* device being resilvered */
	PoolStatusOfflineDev      /* device offline */
	PoolStatusRemovedDev      /* removed device */
	PoolStatusRebuilding      /* device being rebuilt */
	PoolStatusRebuildScrub    /* recommend scrubbing the pool */
	PoolStatusNonNativeAshift /* (e.g. 512e dev with ashift of 9) */

	/*
	 * Finally, the following indicates a healthy pool.
	 */
	PoolStatusOk
)

// PoolState type representing pool state
type PoolState uint64

// Possible ZFS pool states
const (
	PoolStateActive            PoolState = iota /* In active use		*/
	PoolStateExported                           /* Explicitly exported		*/
	PoolStateDestroyed                          /* Explicitly destroyed		*/
	PoolStateSpare                              /* Reserved for hot spare use	*/
	PoolStateL2cache                            /* Level 2 ARC device		*/
	PoolStateUninitialized                      /* Internal spa_t state		*/
	PoolStateUnavailable                        /* Internal libzfs state	*/
	PoolStatePotentiallyActive                  /* Internal libzfs state	*/
)

func (ps PoolState) String() string {
	str := C.GoString(C.zpool_pool_state_to_name(C.pool_state_t(ps)))
	return strings.ToLower(str)
}

// PoolScanStat - Pool scan statistics
type PoolScanStat struct {
	// Values stored on disk
	Func      ScanFunc  // Current scan function e.g. none, scrub ...
	State     ScanState // Current scan state e.g. scanning, finished ...
	StartTime time.Time // Scan start time
	EndTime   time.Time // Scan end time
	ToExamine uint64    // Total bytes to scan
	Examined  uint64    // Total bytes scanned
	Skipped   uint64    // Total bytes skipped
	Processed uint64    // Total bytes processed
	Errors    uint64    // Scan errors
	// Values not stored on disk
	PassExam       uint64    // Examined bytes per scan pass
	PassStart      time.Time // Start time of scan pass
	PassScrubPause time.Time
}

// ExportedPool is type representing ZFS pool available for import
type ExportedPool struct {
	VDevs   VDevTree
	Name    string
	Comment string
	GUID    uint64
	State   PoolState
	Status  PoolStatus
}

// PoolPropertyValue ZFS pool property value
type PoolPropertyValue struct {
	Property PoolProperty
	Source   PropertySource
	Value    string
}

// Pool object represents handler to single ZFS pool
// Map of all ZFS pool properties, changing any of this will not affect ZFS
// pool, for that use SetProperty( name, value string) method of the pool
// object. This map is initial loaded when ever you open or create pool to
// give easy access to listing all available properties. It can be refreshed
// with up to date values with call to (*Pool) ReloadProperties
type Pool struct {
	handle *C.zpool_handle_t
	name   string
}

func (p *Pool) Close() {
	handle := p.handle
	p.handle = nil
	allPools.Delete(handle)
	C.zpool_close(handle)
}

func (p *Pool) LibZFS() *LibZFS {
	l, _ := libzfs.Load(C.zpool_get_handle(p.handle))
	return l
}

func (p *Pool) Name() string {
	return C.GoString(C.zpool_get_name(p.handle))
}

// State get ZFS pool state
// Return the state of the pool (ACTIVE or UNAVAILABLE)
func (p *Pool) State() PoolState {
	return PoolState(C.zpool_get_state(p.handle))
}

// Status get pool status. Let you check if pool healthy.
func (p *Pool) Status() PoolStatus {
	// TODO: maintain and return zpool errata
	var errata C.zpool_errata_t
	return PoolStatus(C.zpool_get_status(p.handle, nil, &errata))
}

func (p *Pool) Config() (NVList, error) {
	config := C.zpool_get_config(p.handle, nil)
	if config == nil {
		return NVList{}, p.LibZFS().Errno()
	}

	return NewNVList(config), nil
}

// VDevTree - Fetch pool's current vdev tree configuration, state and stats
func (p *Pool) VDevTree() (VDevTree, error) {
	config, err := p.Config()
	if err != nil {
		return VDevTree{}, fmt.Errorf("failed to get zpool config: %w", err)
	}
	nvl, err := config.LookupNVList(PoolConfigVdevTree)
	if err != nil {
		return VDevTree{}, fmt.Errorf("failed to fetch vdev tree: %w", err)
	}

	return VDevTree{
		pool: p,
		nvl:  nvl,
	}, nil
}

func (p *Pool) Get(prop PoolProperty) (*PoolPropertyValue, error) {
	var source C.int
	var propBuf = make([]byte, 4096)

	/*
		* Retrieve a property from the given object.  If 'literal' is specified, then
		* numbers are left as exact values.  Otherwise, numbers are converted to a
		* human-readable form.
		*
		* Returns 0 on success, or -1 on error.

		int zfs_prop_get(zfs_handle_t *zhp, zfs_prop_t prop, char *propbuf,
			size_t proplen, zprop_source_t *src, char *statbuf, size_t statlen, boolean_t literal)
	*/
	ret := C.zpool_get_prop(
		p.handle, C.zpool_prop_t(prop),
		(*C.char)(unsafe.Pointer(&propBuf[0])), 4096,
		(*C.zprop_source_t)(unsafe.Pointer(&source)),
		C.B_TRUE,
	)

	if ret != 0 {
		return nil, p.LibZFS().Errno()
	}

	return &PoolPropertyValue{
		Property: prop,
		Source:   PropertySource(source),
		Value:    string(propBuf[:bytes.IndexByte(propBuf, 0)]),
	}, nil
}

// PoolOpen opens a single dataset
func (l *LibZFS) PoolOpen(path string) (*Pool, error) {
	csPath := C.CString(path)
	defer C.free(unsafe.Pointer(csPath))

	handle := C.zpool_open(l.handle, csPath)
	if handle == nil {
		return nil, l.Errno()
	}

	return l.getPool(handle), nil
}

type Pools []*Pool

func (ps Pools) Close() {
	for _, p := range ps {
		p.Close()
	}
}

// PoolOpenAll open all active ZFS pools on current system.
// Returns array of Pool handlers, each have to be closed after not needed
// anymore. Call Pool.Close() method.
func (l *LibZFS) PoolOpenAll() (Pools, error) {
	var handles list[*C.zpool_handle_t]
	defer handles.clear()

	l.lock.Lock()
	err := C.zpool_iter(l.handle, listAppend, handles.pointer())
	l.lock.Unlock()
	if int(err) != 0 {
		return nil, l.Errno()
	}

	var pools = make([]*Pool, handles.len())
	_ = handles.iter(func(i int, handle *C.zpool_handle_t) error {
		pools[i] = l.getPool(handle)
		return nil
	})
	return pools, nil
}
