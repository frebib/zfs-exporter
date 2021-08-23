package zfs

/*
#include <libnvpair.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"time"
	"unsafe"
)

var (
	ErrNotFound    = errors.New("nvlist named value not found")
	ErrUnknownType = errors.New("nvpair type unknown")
)

type NVPair struct {
	handle *C.nvpair_t
}

func NewNVPair(hdl *C.nvpair_t) NVPair {
	return NVPair{handle: hdl}
}

func (nvp NVPair) Name() string {
	return C.GoString(C.nvpair_name(nvp.handle))
}

func (nvp NVPair) Type() NVType {
	return NVType(C.nvpair_type(nvp.handle))
}

func (nvp NVPair) Bool() (bool, bool) {
	var val C.boolean_t
	ret := C.nvpair_value_boolean_value(nvp.handle, &val)
	if ret != 0 {
		return false, false
	}
	return val == C.B_TRUE, true
}

func (nvp NVPair) Byte() (byte, bool) {
	var val C.uchar_t
	ret := C.nvpair_value_byte(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return byte(val), true
}

func (nvp NVPair) Int8() (int8, bool) {
	var val C.int8_t
	ret := C.nvpair_value_int8(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return int8(val), true
}

func (nvp NVPair) Int16() (int16, bool) {
	var val C.int16_t
	ret := C.nvpair_value_int16(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return int16(val), true
}

func (nvp NVPair) Int32() (int32, bool) {
	var val C.int32_t
	ret := C.nvpair_value_int32(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return int32(val), true
}

func (nvp NVPair) Int64() (int64, bool) {
	var val C.int64_t
	ret := C.nvpair_value_int64(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return int64(val), true
}

func (nvp NVPair) Uint8() (uint8, bool) {
	var val C.uint8_t
	ret := C.nvpair_value_uint8(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return uint8(val), true
}

func (nvp NVPair) Uint16() (uint16, bool) {
	var val C.uint16_t
	ret := C.nvpair_value_uint16(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return uint16(val), true
}

func (nvp NVPair) Uint32() (uint32, bool) {
	var val C.uint32_t
	ret := C.nvpair_value_uint32(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return uint32(val), true
}

func (nvp NVPair) Uint64() (uint64, bool) {
	var val C.uint64_t
	ret := C.nvpair_value_uint64(nvp.handle, &val)
	if ret != 0 {
		return 0, false
	}
	return uint64(val), true
}

func (nvp NVPair) String() (string, bool) {
	var val *C.char
	ret := C.nvpair_value_string(nvp.handle, &val)
	if ret != 0 {
		return "", false
	}
	return C.GoString(val), true
}

func (nvp NVPair) Time() (time.Time, bool) {
	var ns C.hrtime_t // is a C long long int, aka int64
	ret := C.nvpair_value_hrtime(nvp.handle, &ns)
	if ret != 0 {
		return time.Time{}, false
	}
	return time.Unix(0, int64(ns)), true
}

func (nvp NVPair) NVList() (NVList, bool) {
	var val *C.nvlist_t
	ret := C.nvpair_value_nvlist(nvp.handle, &val)
	if ret != 0 {
		return NVList{}, false
	}
	return NVList{handle: val}, true
}

func (nvp NVPair) BoolArray() ([]bool, bool) {
	var val *C.boolean_t
	var nelem C.uint_t
	ret := C.nvpair_value_boolean_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]bool, count)
	for i := 0; i < count; i++ {
		elem := *(*C.boolean_t)(ptrIndex(unsafe.Pointer(val), i))
		if elem == C.B_TRUE {
			arr[i] = true
		}
	}
	return arr, true
}

func (nvp NVPair) ByteArray() ([]byte, bool) {
	var val *C.uchar
	var nelem C.uint_t
	ret := C.nvpair_value_byte_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]byte, count)
	for i := 0; i < count; i++ {
		arr[i] = byte(*(*C.uchar)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Int8Array() ([]int8, bool) {
	var val *C.int8_t
	var nelem C.uint_t
	ret := C.nvpair_value_int8_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]int8, count)
	for i := 0; i < count; i++ {
		arr[i] = int8(*(*C.int8_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Int16Array() ([]int16, bool) {
	var val *C.int16_t
	var nelem C.uint_t
	ret := C.nvpair_value_int16_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]int16, count)
	for i := 0; i < count; i++ {
		arr[i] = int16(*(*C.int16_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Int32Array() ([]int32, bool) {
	var val *C.int32_t
	var nelem C.uint_t
	ret := C.nvpair_value_int32_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]int32, count)
	for i := 0; i < count; i++ {
		arr[i] = int32(*(*C.int32_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Int64Array() ([]int64, bool) {
	var val *C.int64_t
	var nelem C.uint_t
	ret := C.nvpair_value_int64_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]int64, count)
	for i := 0; i < count; i++ {
		arr[i] = int64(*(*C.int64_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Uint8Array() ([]uint8, bool) {
	var val *C.uint8_t
	var nelem C.uint_t
	ret := C.nvpair_value_uint8_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]uint8, count)
	for i := 0; i < count; i++ {
		arr[i] = uint8(*(*C.uint8_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Uint16Array() ([]uint16, bool) {
	var val *C.uint16_t
	var nelem C.uint_t
	ret := C.nvpair_value_uint16_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]uint16, count)
	for i := 0; i < count; i++ {
		arr[i] = uint16(*(*C.uint16_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Uint32Array() ([]uint32, bool) {
	var val *C.uint32_t
	var nelem C.uint_t
	ret := C.nvpair_value_uint32_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]uint32, count)
	for i := 0; i < count; i++ {
		arr[i] = uint32(*(*C.uint32_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) Uint64Array() ([]uint64, bool) {
	var val *C.uint64_t
	var nelem C.uint_t
	ret := C.nvpair_value_uint64_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]uint64, count)
	for i := 0; i < count; i++ {
		arr[i] = uint64(*(*C.uint64_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) StringArray() ([]string, bool) {
	var val **C.char
	var nelem C.uint_t
	ret := C.nvpair_value_string_array(nvp.handle, &val, &nelem)
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = C.GoString(*(**C.char)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, true
}

func (nvp NVPair) NVListArray() ([]NVList, bool) {
	var val **C.nvlist_t
	var nelem C.uint_t
	ret := (C.nvpair_value_nvlist_array(nvp.handle, &val, &nelem))
	if ret != 0 {
		return nil, false
	}

	count := int(nelem)
	arr := make([]NVList, count)
	for i := 0; i < count; i++ {
		ptr := *(**C.nvlist_t)(ptrIndex(unsafe.Pointer(val), i))
		arr[i] = NVList{handle: ptr}
	}
	return arr, true
}

func (nvp NVPair) Value() interface{} {
	value, ok := func() (interface{}, bool) {
		switch nvp.Type() {
		case TypeBoolean:
			// https://illumos.org/man/3nvpair/nvpair_value_byte
			// There is no nvpair_value_boolean(); the existence of the name implies
			// the value is true.
			var b = true
			return &b, true

		case TypeBooleanValue:
			return nvp.Bool()
		case TypeByte:
			return nvp.Byte()
		case TypeInt8:
			return nvp.Int8()
		case TypeInt16:
			return nvp.Int16()
		case TypeInt32:
			return nvp.Int32()
		case TypeInt64:
			return nvp.Int64()
		case TypeUint8:
			return nvp.Uint8()
		case TypeUint16:
			return nvp.Uint16()
		case TypeUint32:
			return nvp.Uint32()
		case TypeUint64:
			return nvp.Uint64()
		case TypeString:
			return nvp.String()
		case TypeTime:
			return nvp.Time()
		case TypeNVList:
			return nvp.NVList()
		case TypeBooleanArray:
			return nvp.BoolArray()
		case TypeByteArray:
			return nvp.ByteArray()
		case TypeInt8Array:
			return nvp.Int8Array()
		case TypeInt16Array:
			return nvp.Int16Array()
		case TypeInt32Array:
			return nvp.Int32Array()
		case TypeInt64Array:
			return nvp.Int64Array()
		case TypeUint8Array:
			return nvp.Uint8Array()
		case TypeUint16Array:
			return nvp.Uint16Array()
		case TypeUint32Array:
			return nvp.Uint32Array()
		case TypeUint64Array:
			return nvp.Uint64Array()
		case TypeStringArray:
			return nvp.StringArray()
		case TypeNVListArray:
			return nvp.NVListArray()

		default:
			panic(fmt.Sprintf("nvpair.Value() unrecognised type: %s", nvp.Type()))
		}
	}()

	if !ok {
		panic(fmt.Sprintf("nvpair.Value() programming error for %s", nvp.Type()))
	}

	return value
}

func ptrIndex(ptr unsafe.Pointer, n int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(C.sizeof_size_t*n))
}

func nvlistLookupError(ret C.int) error {
	/*
		static int
		nvlist_lookup_common(nvlist_t *nvl, const char *name, data_type_t type,
		    uint_t *nelem, void *data)
		{
			if (name == NULL || nvl == NULL || nvl->nvl_priv == 0)
				return (EINVAL);

			if (!(nvl->nvl_nvflag & (NV_UNIQUE_NAME | NV_UNIQUE_NAME_TYPE)))
				return (ENOTSUP);

			nvpair_t *nvp = nvt_lookup_name_type(nvl, name, type);
			if (nvp == NULL)
				return (ENOENT);

			return (nvpair_value_common(nvp, type, nelem, data));
		}
	*/
	switch ret {
	case 0:
		return nil

	case C.ENOENT:
		return ErrNotFound

	case C.ENOTSUP:
		return ErrUnknownType

	case C.EINVAL:
		// nvpair_value_common only returns EINVAL when a passed pointer is
		// null, so if we did our job correctly this should never happen
		panic("nvpair value programming error")

	default:
		panic(fmt.Sprintf("unknown return value from nvpair_value_common: %d", int(ret)))
	}
}
