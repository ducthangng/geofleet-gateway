package copier

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// MustCopy performs a deep copy of FromA into ToB.
// It panics if serialization fails or if ToB is not a pointer.
func MustCopy(toB, fromA any) {
	// Pre-check: ToB must be a pointer to be settable
	if reflect.ValueOf(toB).Kind() != reflect.Ptr {
		panic("copier.MustCopy: ToB must be a pointer")
	}

	// 1. convert A -> []byte (Serialize)
	bytes, err := json.Marshal(fromA)
	if err != nil {
		panic(fmt.Errorf("copier.MustCopy: failed to marshal FromA: %w", err))
	}

	// 2. convert []byte -> B (Deserialize)
	// (We don't convert B to bytes; we copy the bytes FROM A INTO B)
	err = json.Unmarshal(bytes, toB)
	if err != nil {
		panic(fmt.Errorf("copier.MustCopy: failed to unmarshal into ToB: %w", err))
	}
}
