package zfs

/*
#include <libzfs.h>
#include <zfs_prop.h>
#include <libintl.h>
#include <libuutil.h>
*/
import "C"

import (
	"bytes"
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
	return C.GoString(C.zfs_prop_to_name(C.zfs_prop_t(dp)))
}

func (dp DatasetProperty) Type() PropertyType {
	return PropertyType(C.zfs_prop_get_type(C.zfs_prop_t(dp)))
}

// Dataset - ZFS dataset object
type Dataset struct {
	handle *C.zfs_handle_t
}

func (d *Dataset) Close() {
	handle := d.handle
	d.handle = nil
	allDatasets.Delete(handle)
	C.zfs_close(handle)
}

func (d *Dataset) LibZFS() *LibZFS {
	l, _ := libzfs.Load(C.zfs_get_handle(d.handle))
	return l
}

func (d *Dataset) Name() string {
	return C.GoString(C.zfs_get_name(d.handle))
}

func (d *Dataset) Type() DatasetType {
	return DatasetType(C.zfs_get_type(d.handle))
}

func (d *Dataset) Pool() *Pool {
	return d.LibZFS().getPool(C.zfs_get_pool_handle(d.handle))
}

func (d *Dataset) Get(prop DatasetProperty) (DatasetPropertyValue, error) {
	var ret, src C.int

	switch prop.Type() {
	case PropertyTypeNumber:
		var value uint64
		ret = C.zfs_prop_get_numeric(
			d.handle, C.zfs_prop_t(prop),
			(*C.uint64_t)(unsafe.Pointer(&value)),
			(*C.zprop_source_t)(unsafe.Pointer(&src)),
			nil, 0,
		)
		if ret != 0 {
			return nil, d.LibZFS().Errno()
		}
		return &DatasetPropertyNumber{
			property: prop,
			source:   PropertySource(src),
			value:    value,
		}, nil

	case PropertyTypeIndex:
		var value uint64
		ret = C.zfs_prop_get_numeric(
			d.handle, C.zfs_prop_t(prop),
			(*C.uint64_t)(unsafe.Pointer(&value)),
			(*C.zprop_source_t)(unsafe.Pointer(&src)),
			nil, 0,
		)
		if ret != 0 {
			return nil, d.LibZFS().Errno()
		}
		return &DatasetPropertyIndex{
			property: prop,
			source:   PropertySource(src),
			value:    value,
		}, nil

	case PropertyTypeString:
		var propBuf = make([]byte, C.ZFS_MAXPROPLEN)

		/*
			* Retrieve a property from the given object.  If 'literal' is specified, then
			* numbers are left as exact values.  Otherwise, numbers are converted to a
			* human-readable form.
			*
			* Returns 0 on success, or -1 on error.

			int zfs_prop_get(zfs_handle_t *zhp, zfs_prop_t prop, char *propbuf,
				size_t proplen, zprop_source_t *src, char *statbuf, size_t statlen, boolean_t literal)
		*/
		ret = C.zfs_prop_get(
			d.handle, C.zfs_prop_t(prop),
			(*C.char)(unsafe.Pointer(&propBuf[0])), (C.ulong)(len(propBuf)),
			(*C.zprop_source_t)(unsafe.Pointer(&src)),
			nil, 0,
			C.B_TRUE,
		)

		if ret != 0 {
			return nil, d.LibZFS().Errno()
		}

		return &DatasetPropertyString{
			property: prop,
			source:   PropertySource(src),
			value:    string(propBuf[:bytes.IndexByte(propBuf, 0)]),
		}, nil
	}

	panic("unknown property type")
}

func (d *Dataset) Gets(props ...DatasetProperty) (map[DatasetProperty]DatasetPropertyValue, error) {
	vals := make(map[DatasetProperty]DatasetPropertyValue, len(props))

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
	var handles list[*C.zfs_handle_t]
	defer handles.clear()

	l := d.LibZFS()
	if types&DatasetTypeFilesystem == DatasetTypeFilesystem {
		ret := C.zfs_iter_filesystems(d.handle, listAppend, handles.pointer())
		if int(ret) != 0 {
			_ = handles.iter(datasetClose)
			return nil, l.Errno()
		}
	}
	if types&DatasetTypeSnapshot == DatasetTypeSnapshot {
		tries := 0
		for {
			ret := C.zfs_iter_snapshots(d.handle, C.B_FALSE, listAppend, handles.pointer(), 0, 0)
			if ret != 0 {
				// Try again if we get EINTR; it's most likely a random signal
				if tries < 3 && C.libzfs_errno(l.handle) == C.EZFS_INTR {
					tries++
					continue
				}
				_ = handles.iter(datasetClose)
				return nil, l.Errno()
			}
			break
		}
	}
	if types&DatasetTypeBookmark == DatasetTypeBookmark {
		ret := C.zfs_iter_bookmarks(d.handle, listAppend, handles.pointer())
		if int(ret) != 0 {
			_ = handles.iter(datasetClose)
			return nil, l.Errno()
		}
	}

	// Shortcut
	if handles.len() == 0 {
		return nil, nil
	}

	datasets := d.LibZFS().datasetInitAll(handles.slice(), types)
	if depth == -1 || depth > 0 {
		if depth > 0 {
			depth--
		}
		for _, d := range datasets {
			// recurse
			children, err := d.Children(types, depth)
			if err != nil {
				for _, d := range datasets {
					d.Close()
				}
				return nil, err
			}
			datasets = append(datasets, children...)
		}
	}

	return datasets, nil
}

func (l *LibZFS) datasetInitAll(handles []*C.zfs_handle_t, types DatasetType) (datasets []*Dataset) {
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

		datasets = append(datasets, l.getDataset(handle))
	}
	return datasets
}

// DatasetOpen opens a single dataset
func (l *LibZFS) DatasetOpen(path string) (*Dataset, error) {
	csPath := C.CString(path)
	defer C.free(unsafe.Pointer(csPath))

	handle := C.zfs_open(l.handle, csPath, C.ZFS_TYPE_DATASET)
	if handle == nil {
		return nil, l.Errno()
	}

	return &Dataset{handle: handle}, nil
}

// DatasetOpenAll recursive get handles to all available datasets on system
// (file-systems, volumes or snapshots).
func (l *LibZFS) DatasetOpenAll(types DatasetType, depth int) ([]*Dataset, error) {
	var handles list[*C.zfs_handle_t]
	defer handles.clear()

	l.lock.Lock()
	ret := C.zfs_iter_root(l.handle, listAppend, handles.pointer())
	l.lock.Unlock()
	if int(ret) != 0 {
		return nil, fmt.Errorf("zfs_iter_root returned %d", int(ret))
	}

	roots := l.datasetInitAll(handles.slice(), types)
	var datasets = roots[:]

	if depth == -1 || depth > 0 {
		if depth > 0 {
			depth--
		}
		for _, root := range roots {
			children, err := root.Children(types, depth)
			if err != nil {
				for _, d := range datasets {
					d.Close()
				}
				return nil, err
			}
			datasets = append(datasets, children...)
		}
	}

	return datasets, nil
}

func datasetClose(_ int, handle *C.zfs_handle_t) error {
	C.zfs_close(handle)
	return nil
}
