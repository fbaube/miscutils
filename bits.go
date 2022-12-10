package miscutils

// IsBitSet returns the numbered bit. Bit 0 is rightmost.
// It may or may not work for more than 16 bits in flagset.
// .
func IsBitSet(flagset, index byte) bool {
	return flagset&(1<<index) != 0
}
