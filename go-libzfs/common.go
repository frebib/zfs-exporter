package zfs

/*
#cgo CFLAGS: -I /usr/include/libzfs -I /usr/include/libspl -DHAVE_IOCTL_IN_SYS_IOCTL_H -D_GNU_SOURCE
#cgo LDFLAGS: -lzfs -lzpool -lnvpair -lzfs_core -luutil

#include <stdlib.h>
#include <libzfs.h>
*/
import "C"
