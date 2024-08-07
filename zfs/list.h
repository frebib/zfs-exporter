#ifndef ZFS_H
#define ZFS_H

#include <stdlib.h>

struct list {
	size_t size;
	size_t next;
	void **data;
};

int list_append(void *data, struct list *l);
#endif