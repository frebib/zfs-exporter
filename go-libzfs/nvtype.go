package zfs

//go:generate stringer -type NVType
type NVType uint

const (
	TypeUnknown NVType = iota
	TypeBoolean
	TypeByte
	TypeInt16
	TypeUint16
	TypeInt32
	TypeUint32
	TypeInt64
	TypeUint64
	TypeString
	TypeByteArray
	TypeInt16Array
	TypeUint16Array
	TypeInt32Array
	TypeUint32Array
	TypeInt64Array
	TypeUint64Array
	TypeStringArray
	TypeTime
	TypeNVList
	TypeNVListArray
	TypeBooleanValue
	TypeInt8
	TypeUint8
	TypeBooleanArray
	TypeInt8Array
	TypeUint8Array
)
