package zfs

import (
	"strings"
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

// DatasetPropertyValue ZFS dataset property value
type DatasetPropertyValue struct {
	Property DatasetProperty
	Source   PropertySource
	Inherit  string
	Value    string
}
