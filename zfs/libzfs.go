package zfs

// cgo CFLAGS reference
// https://github.com/johnramsden/zectl/pull/34/commits/f1531921899c8114943cd62b519d977d24f819bb

/*
#cgo CFLAGS: -I /usr/include/libzfs -I /usr/include/libspl -DHAVE_IOCTL_IN_SYS_IOCTL_H -D_GNU_SOURCE -D__USE_LARGEFILE64=1 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE
#cgo LDFLAGS: -lzfs -lzpool -lnvpair -lzfs_core -luutil

#include <stdlib.h>
#include <libzfs.h>
*/
import "C"
import (
	"errors"
	"sync"
	"syscall"

	"github.com/puzpuzpuz/xsync/v3"
)

// Persist global mappings between C pointers and their Golang counterparts to
// ensure we only ever have a 1:1 mapping. Locking and memory destruction relies
// upon a strict and persistent object lifetime.
var (
	libzfs      = xsync.NewMapOf[*C.libzfs_handle_t, *LibZFS]()
	allPools    = xsync.NewMapOf[*C.zpool_handle_t, *Pool]()
	allDatasets = xsync.NewMapOf[*C.zfs_handle_t, *Dataset]()
)

func New() (*LibZFS, error) {
	handle, err := C.libzfs_init()
	if handle == nil {
		errno := err.(syscall.Errno)
		return nil, errors.New(C.GoString(C.libzfs_error_init(C.int(errno))))
	}
	l := &LibZFS{handle: handle}
	libzfs.Store(handle, l)
	return l, nil
}

type LibZFS struct {
	// include/libzfs_impl.h
	handle *C.libzfs_handle_t
	lock   sync.Mutex
}

func (l *LibZFS) Close() {
	handle := l.handle
	l.handle = nil
	libzfs.Delete(handle)
	C.libzfs_fini(handle)
}

func (l *LibZFS) getPool(handle *C.zpool_handle_t) *Pool {
	pool, _ := allPools.LoadOrStore(handle, &Pool{handle: handle})
	return pool
}

func (l *LibZFS) getDataset(handle *C.zfs_handle_t) *Dataset {
	dataset, _ := allDatasets.LoadOrStore(handle, &Dataset{handle: handle})
	return dataset
}
