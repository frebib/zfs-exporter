package zfs

/*
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
