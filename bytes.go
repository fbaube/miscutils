package miscutils

// DupeByteSlice allocates & returns a duplicate of the input byte slice.
func DupeByteSlice(in []byte) (out []byte) {
	out = make([]byte, len(in))
	copy(out, in)
	return out
}
