package zfs

/*
#include <libzfs.h>
#include <zfs_prop.h>
#include <libintl.h>
#include <libuutil.h>

extern int datasetSlice(zfs_handle_t *h, void *ptr);
*/
import "C"

import (
	"bytes"
	"errors"
	"fmt"
	"unsafe"
)

// DatasetType defines enum of dataset types
type DatasetType int32

const (
	DatasetTypeUnknown DatasetType = 0
	// DatasetTypeFilesystem - file system dataset
	DatasetTypeFilesystem = 1 << (iota - 1)
	// DatasetTypeSnapshot - snapshot of dataset
	DatasetTypeSnapshot
	// DatasetTypeVolume - volume (virtual block device) dataset
	DatasetTypeVolume
	// DatasetTypePool - pool dataset
	DatasetTypePool
	// DatasetTypeBookmark - bookmark dataset
	DatasetTypeBookmark
)

func (t DatasetType) String() string {
	switch t {
	case DatasetTypeFilesystem:
		return "filesystem"
	case DatasetTypeSnapshot:
		return "snapshot"
	case DatasetTypeVolume:
		return "volume"
	case DatasetTypePool:
		return "pool"
	case DatasetTypeBookmark:
		return "bookmark"
	default:
		return "unknown"
	}
}

type DatasetProperty int

/*
 * Dataset properties are identified by these constants and must be added to
 * the end of this list to ensure that external consumers are not affected
 * by the change. If you make any changes to this list, be sure to update
 * the property table in module/zcommon/zfs_prop.c.
 */
const (
	DatasetPropCont DatasetProperty = iota - 2
	DatasetPropInval

	DatasetPropType
	DatasetPropCreation
	DatasetPropUsed
	DatasetPropAvailable
	DatasetPropReferenced
	DatasetPropCompressratio
	DatasetPropMounted
	DatasetPropOrigin
	DatasetPropQuota
	DatasetPropReservation
	DatasetPropVolsize
	DatasetPropVolblocksize
	DatasetPropRecordsize
	DatasetPropMountpoint
	DatasetPropSharenfs
	DatasetPropChecksum
	DatasetPropCompression
	DatasetPropAtime
	DatasetPropDevices
	DatasetPropExec
	DatasetPropSetuid
	DatasetPropReadonly
	DatasetPropZoned
	DatasetPropSnapdir
	DatasetPropAclmode
	DatasetPropAclinherit
	DatasetPropCreateTXG
	DatasetPropName /* not exposed to the user */
	DatasetPropCanMount
	DatasetPropiSCSIOptions /* not exposed to the user */
	DatasetPropXattr
	DatasetPropNumClones /* not exposed to the user */
	DatasetPropCopies
	DatasetPropVersion
	DatasetPropUtf8Only
	DatasetPropNormalize
	DatasetPropCase
	DatasetPropVScan
	DatasetPropNbmand
	DatasetPropShareSMB
	DatasetPropRefQuota
	DatasetPropRefReservation
	DatasetPropGUID
	DatasetPropPrimaryCache
	DatasetPropSecondaryCache
	DatasetPropUsedSnap
	DatasetPropUsedDS
	DatasetPropUsedChild
	DatasetPropUsedRefReservation
	DatasetPropUserAccounting /* not exposed to the user */
	DatasetPropStmfShareInfo  /* not exposed to the user */
	DatasetPropDeferDestroy
	DatasetPropUserRefs
	DatasetPropLogBias
	DatasetPropUnique /* not exposed to the user */
	DatasetPropObjSetID
	DatasetPropDedup
	DatasetPropMlsLabel
	DatasetPropSync
	DatasetPropDnodeSize
	DatasetPropRefRatio
	DatasetPropWritten
	DatasetPropClones
	DatasetPropLogicalUsed
	DatasetPropLogicalReferenced
	DatasetPropInconsistent /* not exposed to the user */
	DatasetPropVolMode
	DatasetPropFilesystemLimit
	DatasetPropSnapshotLimit
	DatasetPropFilesystemCount
	DatasetPropSnapshotCount
	DatasetPropSnapDev
	DatasetPropAclType
	DatasetPropSelinuxContext
	DatasetPropSelinuxFSContext
	DatasetPropSelinuxDefContext
	DatasetPropSelinuxRootContext
	DatasetPropRelAtime
	DatasetPropRedundantMetadata
	DatasetPropOverlay
	DatasetPropPrevSnap
	DatasetPropReceiveResumeToken
	DatasetPropEncryption
	DatasetPropKeyLocation
	DatasetPropKeyFormat
	DatasetPropPBKDF2Salt
	DatasetPropPBKDF2Iters
	DatasetPropEncryptionRoot
	DatasetPropKeyGUID
	DatasetPropKeyStatus
	DatasetPropRemapTXG /* obsolete - no longer used */
	DatasetPropSpecialSmallBlocks
	DatasetPropIvsetGUID /* not exposed to the user */
	DatasetPropRedacted
	DatasetPropRedactSnaps

	DatasetNumProps
)

func (dp DatasetProperty) String() string {
	ptr := C.zfs_prop_to_name(C.zfs_prop_t(dp))
	return C.GoString(ptr)
}

func (dp DatasetProperty) Type() PropertyType {
	return PropertyType(C.zfs_prop_get_type(C.zfs_prop_t(dp)))
}

// Dataset - ZFS dataset object
type Dataset struct {
	handle *C.zfs_handle_t
	name   string
	typ    DatasetType
}

func (d Dataset) Close() {
	C.zfs_close(d.handle)
}

func (d Dataset) LibZFS() *LibZFS {
	return &LibZFS{
		handle: C.zfs_get_handle(d.handle),
	}
}

func (d *Dataset) Name() string {
	if d.name == "" {
		d.name = C.GoString(C.zfs_get_name(d.handle))
	}

	return d.name
}

func (d *Dataset) Type() DatasetType {
	// Cache on first use
	if d.typ == DatasetTypeUnknown {
		d.typ = DatasetType(C.zfs_get_type(d.handle))
	}

	return d.typ
}

func (d *Dataset) Pool() *Pool {
	return &Pool{handle: C.zfs_get_pool_handle(d.handle)}
}

func (d Dataset) Get(prop DatasetProperty) (*DatasetPropertyValue, error) {
	var source C.int
	var statBuf = make([]byte, 1024)
	var propBuf = make([]byte, 4096) //create my buffer

	/*
		* Retrieve a property from the given object.  If 'literal' is specified, then
		* numbers are left as exact values.  Otherwise, numbers are converted to a
		* human-readable form.
		*
		* Returns 0 on success, or -1 on error.

		int zfs_prop_get(zfs_handle_t *zhp, zfs_prop_t prop, char *propbuf,
			size_t proplen, zprop_source_t *src, char *statbuf, size_t statlen, boolean_t literal)
	*/
	ret := C.zfs_prop_get(
		d.handle, C.zfs_prop_t(prop), (*C.char)(unsafe.Pointer(&propBuf[0])),
		4096, (*C.zprop_source_t)(unsafe.Pointer(&source)),
		(*C.char)(unsafe.Pointer(&statBuf[0])), 1024, booleanT(true),
	)

	if ret != 0 {
		return nil, d.LibZFS().Errno()
	}

	return &DatasetPropertyValue{
		Property: prop,
		Source:   PropertySource(source),
		Inherit:  string(statBuf[:bytes.IndexByte(statBuf, 0)]),
		Value:    string(propBuf[:bytes.IndexByte(propBuf, 0)]),
	}, nil
}
func (d Dataset) Gets(props ...DatasetProperty) (map[DatasetProperty]*DatasetPropertyValue, error) {
	vals := make(map[DatasetProperty]*DatasetPropertyValue, len(props))

	for _, prop := range props {
		val, err := d.Get(prop)
		if err != nil {
			return nil, err
		}

		if _, ok := vals[prop]; ok {
			return nil, fmt.Errorf("duplicate property requested: '%s'", prop)
		}
		vals[prop] = val
	}

	return vals, nil
}

func (d *Dataset) Children(types DatasetType, depth int) ([]*Dataset, error) {
	var handles []*C.zfs_handle_t

	// cgo has silly type restrictions.
	// This is the only way I could make it compile
	callback := (*[0]byte)(C.datasetSlice)

	if types&DatasetTypeFilesystem == DatasetTypeFilesystem {
		ret := C.zfs_iter_filesystems(d.handle, callback, unsafe.Pointer(&handles))
		if int(ret) != 0 {
			return nil, d.LibZFS().Errno()
		}
	}
	if types&DatasetTypeSnapshot == DatasetTypeSnapshot {
		ret := C.zfs_iter_snapshots(d.handle, C.B_TRUE, callback,
			unsafe.Pointer(&handles), C.ulong(0), C.ulong(0),
		)
		if int(ret) != 0 {
			return nil, d.LibZFS().Errno()
		}
	}
	if types&DatasetTypeBookmark == DatasetTypeBookmark {
		ret := C.zfs_iter_bookmarks(d.handle, callback, unsafe.Pointer(&handles))
		if int(ret) != 0 {
			return nil, d.LibZFS().Errno()
		}
	}

	datasets, err := datasetInitAll(handles, types)
	if err != nil {
		return nil, err
	}

	if depth == -1 || depth > 0 {
		if depth > 0 {
			depth--
		}
		for _, d := range datasets {
			// recurse
			children, err := d.Children(types, depth)
			if err != nil {
				return nil, err
			}
			datasets = append(datasets, children...)
		}
	}

	return datasets, nil
}

//export datasetSlice
// datasetSlice appends the passed zfs_handle_t to a slice of []*C.zfs_handle_t
// passed in via ptr of type unsafe.Pointer(*[]*C.zfs_handle_t). This function
// is intended to be used as a callback to the zfs_iter_* suite of libzfs
// functions, matching signature: int (*zfs_iter_f)(zfs_handle_t*, void*)
func datasetSlice(handle *C.zfs_handle_t, ptr unsafe.Pointer) C.int {
	list := (*[]*C.zfs_handle_t)(ptr)
	*list = append(*list, handle)
	return 0
}

func datasetInit(hdl *C.zfs_handle_t) (*Dataset, error) {
	if hdl == nil {
		return nil, errors.New("zfs handle is nil")
	}

	return &Dataset{handle: hdl}, nil
}

func datasetInitAll(handles []*C.zfs_handle_t, types DatasetType) ([]*Dataset, error) {
	var datasets []*Dataset

	for _, handle := range handles {
		if handle == nil {
			panic("nil zfs_handle_t")
		}

		// Skip unwanted types
		typ := DatasetType(C.zfs_get_type(handle))
		if types&typ == 0 {
			// Clean up
			C.zfs_close(handle)
			continue
		}

		dataset, err := datasetInit(handle)
		if err != nil {
			return nil, err
		}
		datasets = append(datasets, dataset)
	}

	return datasets, nil
}

// DatasetOpen opens a single dataset
func (l *LibZFS) DatasetOpen(path string) (*Dataset, error) {
	csPath := C.CString(path)
	defer C.free(unsafe.Pointer(csPath))

	handle := C.zfs_open(l.Handle(), csPath, C.ZFS_TYPE_DATASET)
	if handle == nil {
		return nil, l.Errno()
	}

	return datasetInit(handle)
}

// DatasetOpenAll recursive get handles to all available datasets on system
// (file-systems, volumes or snapshots).
func (l *LibZFS) DatasetOpenAll(types DatasetType, depth int) ([]*Dataset, error) {
	var handles []*C.zfs_handle_t

	l.namespaceMtx.Lock()
	ret := C.zfs_iter_root(l.Handle(), (*[0]byte)(C.datasetSlice), unsafe.Pointer(&handles))
	l.namespaceMtx.Unlock()
	if int(ret) != 0 {
		return nil, l.Errno()
	}

	roots, err := datasetInitAll(handles, types)
	if err != nil {
		return nil, err
	}

	datasets := append([]*Dataset{}, roots...)

	if depth == -1 || depth > 0 {
		if depth > 0 {
			depth--
		}
		for _, root := range roots {
			children, err := root.Children(types, depth)
			if err != nil {
				return nil, err
			}
			datasets = append(datasets, children...)
		}
	}

	return datasets, nil
}
