package zfs

/*
#include <libzfs.h>
*/
import "C"

// ZFS errors
const (
	ESuccess              = 0               /* no error -- success */
	ENomem                = 2000 + iota - 1 /* out of memory */
	EBadprop                                /* invalid property value */
	EPropreadonly                           /* cannot set readonly property */
	EProptype                               /* property does not apply to dataset type */
	EPropnoninherit                         /* property is not inheritable */
	EPropspace                              /* bad quota or reservation */
	EBadtype                                /* dataset is not of appropriate type */
	EBusy                                   /* pool or dataset is busy */
	EExists                                 /* pool or dataset already exists */
	ENoent                                  /* no such pool or dataset */
	EBadstream                              /* bad backup stream */
	EDsreadonly                             /* dataset is readonly */
	EVoltoobig                              /* volume is too large for 32-bit system */
	EInvalidname                            /* invalid dataset name */
	EBadrestore                             /* unable to restore to destination */
	EBadbackup                              /* backup failed */
	EBadtarget                              /* bad attach/detach/replace target */
	ENodevice                               /* no such device in pool */
	EBaddev                                 /* invalid device to add */
	ENoreplicas                             /* no valid replicas */
	EResilvering                            /* currently resilvering */
	EBadversion                             /* unsupported version */
	EPoolunavail                            /* pool is currently unavailable */
	EDevoverflow                            /* too many devices in one vdev */
	EBadpath                                /* must be an absolute path */
	ECrosstarget                            /* rename or clone across pool or dataset */
	EZoned                                  /* used improperly in local zone */
	EMountfailed                            /* failed to mount dataset */
	EUmountfailed                           /* failed to unmount dataset */
	EUnsharenfsfailed                       /* unshare(1M) failed */
	ESharenfsfailed                         /* share(1M) failed */
	EPerm                                   /* permission denied */
	ENospc                                  /* out of space */
	EFault                                  /* bad address */
	EIo                                     /* I/O error */
	EIntr                                   /* signal received */
	EIsspare                                /* device is a hot spare */
	EInvalconfig                            /* invalid vdev configuration */
	ERecursive                              /* recursive dependency */
	ENohistory                              /* no history object */
	EPoolprops                              /* couldn't retrieve pool props */
	EPoolNotsup                             /* ops not supported for this type of pool */
	EPoolInvalarg                           /* invalid argument for this pool operation */
	ENametoolong                            /* dataset name is too long */
	EOpenfailed                             /* open of device failed */
	ENocap                                  /* couldn't get capacity */
	ELabelfailed                            /* write of label failed */
	EBadwho                                 /* invalid permission who */
	EBadperm                                /* invalid permission */
	EBadpermset                             /* invalid permission set name */
	ENodelegation                           /* delegated administration is disabled */
	EUnsharesmbfailed                       /* failed to unshare over smb */
	ESharesmbfailed                         /* failed to share over smb */
	EBadcache                               /* bad cache file */
	EIsl2CACHE                              /* device is for the level 2 ARC */
	EVdevnotsup                             /* unsupported vdev type */
	ENotsup                                 /* ops not supported on this dataset */
	EActiveSpare                            /* pool has active shared spare devices */
	EUnplayedLogs                           /* log device has unplayed logs */
	EReftagRele                             /* snapshot release: tag not found */
	EReftagHold                             /* snapshot hold: tag already exists */
	ETagtoolong                             /* snapshot hold/rele: tag too long */
	EPipefailed                             /* pipe create failed */
	EThreadcreatefailed                     /* thread create failed */
	EPostsplitOnline                        /* onlining a disk after splitting it */
	EScrubbing                              /* currently scrubbing */
	ENoScrub                                /* no active scrub */
	EDiff                                   /* general failure of zfs diff */
	EDiffdata                               /* bad zfs diff data */
	EPoolreadonly                           /* pool is in read-only mode */
	EScrubpaused                            /* scrub currently paused */
	EActivepool                             /* pool is imported on a different system */
	ECryptofailed                           /* failed to setup encryption */
	ENopending                              /* cannot cancel, no operation is pending */
	ECheckpointExists                       /* checkpoint exists */
	EDiscardingCheckpoint                   /* currently discarding a checkpoint */
	ENoCheckpoint                           /* pool has no checkpoint */
	EDevrmInProgress                        /* a device is currently being removed */
	EVdevTooBig                             /* a device is too big to be used */
	EIocNotsupported                        /* operation not supported by zfs module */
	EToomany                                /* argument list too long */
	EInitializing                           /* currently initializing */
	ENoInitialize                           /* no active initialize */
	EWrongParent                            /* invalid parent dataset (e.g ZVOL) */
	ETrimming                               /* currently trimming */
	ENoTrim                                 /* no active trim */
	ETrimNotsup                             /* device does not support trim */
	ENoResilverDefer                        /* pool doesn't support resilver_defer */
	EExportInProgress                       /* currently exporting the pool */
	EUnknown
)

type Error struct {
	errno   int
	message string
}

func (e Error) Errno() int {
	return e.errno
}
func (e Error) Error() string {
	return e.message
}

func (l *LibZFS) Errno() error {
	errno := C.libzfs_errno(l.handle)
	message := C.libzfs_error_description(l.handle)
	return &Error{
		errno:   int(errno),
		message: C.GoString(message),
	}
}
