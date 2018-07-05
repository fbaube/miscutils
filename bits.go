package miscutils

// IsBitSet returns the numbered bit. Bit 0 is rightmost.
func IsBitSet(flag, index byte) bool {
	return flag&(1<<index) != 0
}
