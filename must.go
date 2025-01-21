package miscutils

// Must enforces error-free func calls: a panic 
// is triggered by a non-nil error returned by any 
// func of the form: MyFunc(..) (someValue, error) 
func Must[T any](val T, err error) T {
    if err != nil {
        panic(err)
    }
    return val
}
