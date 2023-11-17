package miscutils

func Tri[T any](cond bool, a, b T) T {
  if cond {
    return a
  }
  return b
}

