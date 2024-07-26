package zfs

/*
#include "list.h"
#include <stdlib.h>

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

var listAppend = (*[0]byte)(C.list_append)

type list[T any] C.struct_list

func (l *list[T]) pointer() unsafe.Pointer {
	return unsafe.Pointer(l)
}

func (l *list[T]) len() int {
	return int(l.next)
}

// slice extracts (and casts) data from the C list struct into a Go slice
// then frees the list and returns the data.
func (l *list[T]) slice() []T {
	length := l.len()
	ts := make([]T, length)
	for i := 0; i < length; i++ {
		ts[i] = *(*T)(unsafe.Add(unsafe.Pointer(l.data), i*C.sizeof_size_t))
	}
	return ts
}

func (l *list[T]) clear() {
	if l.data != nil {
		C.free(unsafe.Pointer(l.data))
		l.data = nil
	}
	l.size = 0
	l.next = 0
}
