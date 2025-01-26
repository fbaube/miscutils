package miscutils

// DupeByteSlice allocates and returns 
// a duplicate of the input byte slice.
// As with many things in Go, there is
// more than one way to skin a cat. 
func DupeByteSlice(in []byte) (out []byte) {
     	// This message is fastest, see
	// https://github.com/golang/go/issues/55905 
	out = make([]byte, len(in))
	copy(out, in)
	return out
	// Other ways:
	// 1) out = append([]byte(nil), in...)
	// 2) out = append(make([]byte, 0, len(in)), in...)
	// 3) out = slices.Clone(in) 
}
