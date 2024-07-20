package zfs

// cgo CFLAGS reference
// https://github.com/johnramsden/zectl/pull/34/commits/f1531921899c8114943cd62b519d977d24f819bb

/*
#cgo CFLAGS: -I /usr/include/libzfs -I /usr/include/libspl -DHAVE_IOCTL_IN_SYS_IOCTL_H -D_GNU_SOURCE -D__USE_LARGEFILE64=1 -D_LARGEFILE_SOURCE -D_LARGEFILE64_SOURCE
#cgo LDFLAGS: -lzfs -lzpool -lnvpair -lzfs_core -luutil

#include "list.h"
#include <stdlib.h>
#include <libzfs.h>

int list_append(void *data, struct list *l) {
	if (l->next >= l->size) {
		if (l->size == 0) {
			l->size = 4;
			l->next = 0;
		} else {
			l->size *= 2;
		}
		l->data = realloc(l->data, l->size * sizeof(void *));
	}
	l->data[l->next++] = data;
	return 0;
}
*/
import "C"
import (
	"unsafe"
)

// listToSlice extracts (and casts) data from the C list struct into a Go slice
// then frees the list and returns the data.
func listToSlice[T any](l C.struct_list) []T {
	length := (int)(l.next)
	handles := make([]T, length)
	for i := 0; i < length; i++ {
		handles[i] = *(*T)(unsafe.Add(unsafe.Pointer(l.data), i*C.sizeof_size_t))
	}

	// Clean up
	if l.data != nil {
		C.free(unsafe.Pointer(l.data))
	}
	l.size = 0
	l.next = 0

	return handles
}
