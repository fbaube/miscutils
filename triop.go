package miscutils

// Tri implements the ternary operator.
func Tri[T any](cond bool, a, b T) T {
  if cond {
    return a
  }
  return b
}

