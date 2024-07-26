package zfs

/*
#include <libnvpair.h>

char *MODE_W = "w";
*/
import "C"

import (
	"time"
	"unsafe"
)

type NVList struct {
	handle *C.nvlist_t
}

func NewNVList(hdl *C.nvlist_t) NVList {
	return NVList{handle: hdl}
}

func (nvl *NVList) Pointer() *C.nvlist_t {
	return nvl.handle
}

func (nvl *NVList) Exists(name string) bool {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return C.nvlist_exists(nvl.handle, cname) == C.B_TRUE
}

func (nvl *NVList) Lookup(name string) (NVPair, bool) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var nvp *C.nvpair_t
	// Find the 'name'ed nvpair in the nvlist. If 'name' found, the function
	// returns zero and a pointer to the matching nvpair is returned in '*ret'
	// (given 'ret' is non-NULL). See nvlist_lookup_nvpair_ei_sep() for details
	ret := C.nvlist_lookup_nvpair(nvl.handle, cname, &nvp)

	ok := false
	if ret == 0 {
		ok = true
	}

	return NVPair{handle: nvp}, ok
}

func (nvl NVList) LookupBool(name string) (bool, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.boolean_t
	ret := C.nvlist_lookup_boolean_value(nvl.handle, cname, &val)
	if ret != 0 {
		return false, nvlistLookupError(ret)
	}
	return val == C.B_TRUE, nil
}

func (nvl NVList) LookupByte(name string) (byte, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.uchar_t
	ret := C.nvlist_lookup_byte(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return byte(val), nil
}

func (nvl NVList) LookupInt8(name string) (int8, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.int8_t
	ret := C.nvlist_lookup_int8(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return int8(val), nil
}

func (nvl NVList) LookupInt16(name string) (int16, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.int16_t
	ret := C.nvlist_lookup_int16(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return int16(val), nil
}

func (nvl NVList) LookupInt32(name string) (int32, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.int32_t
	ret := C.nvlist_lookup_int32(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return int32(val), nil
}

func (nvl NVList) LookupInt64(name string) (int64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.int64_t
	ret := C.nvlist_lookup_int64(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return int64(val), nil
}

func (nvl NVList) LookupUint8(name string) (uint8, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.uint8_t
	ret := C.nvlist_lookup_uint8(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return uint8(val), nil
}

func (nvl NVList) LookupUint16(name string) (uint16, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.uint16_t
	ret := C.nvlist_lookup_uint16(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return uint16(val), nil
}

func (nvl NVList) LookupUint32(name string) (uint32, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.uint32_t
	ret := C.nvlist_lookup_uint32(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return uint32(val), nil
}

func (nvl NVList) LookupUint64(name string) (uint64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val C.uint64_t
	ret := C.nvlist_lookup_uint64(nvl.handle, cname, &val)
	if ret != 0 {
		return 0, nvlistLookupError(ret)
	}
	return uint64(val), nil
}

func (nvl NVList) LookupString(name string) (string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.char
	ret := C.nvlist_lookup_string(nvl.handle, cname, &val)
	if ret != 0 {
		return "", nvlistLookupError(ret)
	}
	return C.GoString(val), nil
}

func (nvl NVList) LookupTime(name string) (time.Time, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var ns C.hrtime_t // is a C long long int, aka int64
	ret := C.nvlist_lookup_hrtime(nvl.handle, cname, &ns)
	if ret != 0 {
		return time.Time{}, nvlistLookupError(ret)
	}
	return time.Unix(0, int64(ns)), nil
}

func (nvl NVList) LookupNVList(name string) (NVList, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.nvlist_t
	ret := C.nvlist_lookup_nvlist(nvl.handle, cname, &val)
	if ret != 0 {
		return NVList{}, nvlistLookupError(ret)
	}
	return NVList{handle: val}, nil
}

func (nvl NVList) LookupBoolArray(name string) ([]bool, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.boolean_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_boolean_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]bool, count)
	for i := 0; i < count; i++ {
		elem := *(*C.boolean_t)(ptrIndex(unsafe.Pointer(val), i))
		if elem == C.B_TRUE {
			arr[i] = true
		}
	}
	return arr, nil
}

func (nvl NVList) LookupByteArray(name string) ([]byte, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.uchar
	var nelem C.uint_t
	ret := C.nvlist_lookup_byte_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]byte, count)
	for i := 0; i < count; i++ {
		arr[i] = byte(*(*C.uchar)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupInt8Array(name string) ([]int8, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.int8_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_int8_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]int8, count)
	for i := 0; i < count; i++ {
		arr[i] = int8(*(*C.int8_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupInt16Array(name string) ([]int16, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.int16_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_int16_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]int16, count)
	for i := 0; i < count; i++ {
		arr[i] = int16(*(*C.int16_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupInt32Array(name string) ([]int32, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.int32_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_int32_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]int32, count)
	for i := 0; i < count; i++ {
		arr[i] = int32(*(*C.int32_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupInt64Array(name string) ([]int64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.int64_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_int64_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]int64, count)
	for i := 0; i < count; i++ {
		arr[i] = int64(*(*C.int64_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupUint8Array(name string) ([]uint8, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.uint8_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_uint8_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]uint8, count)
	for i := 0; i < count; i++ {
		arr[i] = uint8(*(*C.uint8_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupUint16Array(name string) ([]uint16, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.uint16_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_uint16_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]uint16, count)
	for i := 0; i < count; i++ {
		arr[i] = uint16(*(*C.uint16_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupUint32Array(name string) ([]uint32, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.uint32_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_uint32_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]uint32, count)
	for i := 0; i < count; i++ {
		arr[i] = uint32(*(*C.uint32_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupUint64Array(name string) ([]uint64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val *C.uint64_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_uint64_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]uint64, count)
	for i := 0; i < count; i++ {
		arr[i] = uint64(*(*C.uint64_t)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupStringArray(name string) ([]string, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val **C.char
	var nelem C.uint_t
	ret := C.nvlist_lookup_string_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]string, count)
	for i := 0; i < count; i++ {
		arr[i] = C.GoString(*(**C.char)(ptrIndex(unsafe.Pointer(val), i)))
	}
	return arr, nil
}

func (nvl NVList) LookupNVListArray(name string) ([]NVList, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	var val **C.nvlist_t
	var nelem C.uint_t
	ret := C.nvlist_lookup_nvlist_array(nvl.handle, cname, &val, &nelem)
	if ret != 0 {
		return nil, nvlistLookupError(ret)
	}

	count := int(nelem)
	arr := make([]NVList, count)
	for i := 0; i < count; i++ {
		ptr := *(**C.nvlist_t)(ptrIndex(unsafe.Pointer(val), i))
		arr[i] = NVList{handle: ptr}
	}
	return arr, nil
}

func (nvl NVList) Map() map[string]interface{} {
	if C.nvlist_empty(nvl.handle) != 0 {
		return nil
	}

	elems := make(map[string]interface{})
	handle := C.nvlist_next_nvpair(nvl.handle, nil)
	for handle != nil {
		nvp := NVPair{handle: handle}

		var data interface{}
		switch nvp.Type() {
		case TypeNVList:
			list, _ := nvp.NVList()
			m := list.Map()
			data = &m
		case TypeNVListArray:
			list, _ := nvp.NVListArray()
			maps := make([]map[string]interface{}, len(list))
			for i, nvl := range list {
				maps[i] = nvl.Map()
			}
			data = &maps
		default:
			data = nvp.Value()
		}

		name := nvp.Name()
		elems[name] = data

		handle = C.nvlist_next_nvpair(nvl.handle, nvp.handle)
	}

	return elems
}

func (nvl NVList) JSON(fd uintptr) int {
	file := C.fdopen(C.int(fd), C.MODE_W)
	defer C.fclose(file)
	return int(C.nvlist_print_json(file, nvl.handle))
}
