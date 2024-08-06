package zfs

/*
#include <stdlib.h>
#include <libzfs.h>
*/
import "C"
import (
	"strings"
	"unsafe"
)

type PropertySource int

const (
	SourceNone PropertySource = 1 << iota
	SourceDefault
	SourceTemporary
	SourceLocal
	SourceInherited
	SourceReceived

	SourceAll = SourceNone | SourceDefault | SourceTemporary |
		SourceLocal | SourceInherited | SourceReceived
)

func (ps PropertySource) String() string {
	if ps == 0 {
		return "unknown"
	}
	if ps == SourceNone {
		return "-"
	}

	var sources []string
	if ps&SourceDefault == SourceDefault {
		sources = append(sources, "default")
	}
	if ps&SourceTemporary == SourceTemporary {
		sources = append(sources, "temporary")
	}
	if ps&SourceLocal == SourceLocal {
		sources = append(sources, "local")
	}
	if ps&SourceInherited == SourceInherited {
		sources = append(sources, "inherited")
	}
	if ps&SourceReceived == SourceReceived {
		sources = append(sources, "received")
	}

	if len(sources) == 0 {
		return "unknown"
	}

	return strings.Join(sources, " | ")
}

//go:generate stringer -type PropertyType -trimprefix PropertyType
type PropertyType int

const (
	PropertyTypeNumber PropertyType = iota /* numeric value */
	PropertyTypeString                     /* string value */
	PropertyTypeIndex                      /* numeric value indexed by string */
)

type DatasetPropertyValue interface {
	Type() PropertyType
	Property() DatasetProperty
	Source() PropertySource
}

type DatasetPropertyNumber struct {
	property DatasetProperty
	source   PropertySource
	value    uint64
}

func (d DatasetPropertyNumber) Type() PropertyType {
	return PropertyTypeNumber
}

func (d DatasetPropertyNumber) Property() DatasetProperty {
	return d.property
}

func (d DatasetPropertyNumber) Source() PropertySource {
	return d.source
}

func (d DatasetPropertyNumber) Value() uint64 {
	return d.value
}

type DatasetPropertyIndex struct {
	property DatasetProperty
	source   PropertySource
	value    uint64
}

func (d DatasetPropertyIndex) Type() PropertyType {
	return PropertyTypeIndex
}

func (d DatasetPropertyIndex) Property() DatasetProperty {
	return d.property
}

func (d DatasetPropertyIndex) Source() PropertySource {
	return d.source
}

func (d DatasetPropertyIndex) Name() string {
	var cstr *C.char
	ret := C.zfs_prop_index_to_string(
		(C.zfs_prop_t)(d.property),
		(C.uint64_t)(d.value),
		(**C.char)(unsafe.Pointer(&cstr)),
	)
	if ret != 0 {
		panic("womp")
	}
	return C.GoString(cstr)
}

func (d DatasetPropertyIndex) Value() uint64 {
	return d.value
}

type DatasetPropertyString struct {
	property DatasetProperty
	source   PropertySource
	value    string
}

func (d DatasetPropertyString) Type() PropertyType {
	return PropertyTypeString
}

func (d DatasetPropertyString) Property() DatasetProperty {
	return d.property
}

func (d DatasetPropertyString) Source() PropertySource {
	return d.source
}

func (d DatasetPropertyString) Value() string {
	return d.value
}
