package hashmap

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// This conversion *does not* copy data. Note that casting via "([]byte)(string)" *does* copy data.
// Also note that you *should not* change the byte slice after conversion, because Go strings
// are treated as immutable. This would cause a segmentation violation panic.
func S2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func ParseKey(key any) (r []byte) {
	switch key := key.(type) {
	case []byte:
		return key
	case string:
		return S2b(key)
	case int64:
		r = make([]byte, 8)
		binary.NativeEndian.PutUint64(r, uint64(key))
		return r
	case uint64:
		r = make([]byte, 8)
		binary.NativeEndian.PutUint64(r, key)
		return r
	default:
		panic(fmt.Sprintf("invalid value type: %T", key))
	}
}

func Beq(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
