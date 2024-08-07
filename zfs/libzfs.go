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
)

// include/libzfs_impl.h
type LibZFS struct {
	handle       *C.libzfs_handle_t
	namespaceMtx sync.Mutex
}

func (l *LibZFS) Handle() *C.libzfs_handle_t {
	return l.handle
}

func (l *LibZFS) Close() {
	C.libzfs_fini(l.handle)
}

func New() (*LibZFS, error) {
	handle, err := C.libzfs_init()
	if handle == nil {
		errno := err.(syscall.Errno)
		return nil, errors.New(C.GoString(C.libzfs_error_init(C.int(errno))))
	}
	return &LibZFS{handle: handle}, nil
}
