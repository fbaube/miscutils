package miscutils

import (
	"fmt"
	S "strings"
	"time"
)

// Into starts a timer and (if the string argument is non-empty)
// writes to `os.Stderr`; the timer is ended by passing it to
// `Outa(..)`, which can be placed in a `defer` statement.
// Note that the caller must store the start time,
// but that this means we can have nested timers.
//
// These functions add a dependency on package `time`, but
// it is also possible to rewrite them to use `interface{}`.
func Into(s string) time.Time {
	if s != "" {
		println("[>BEGAN<]", s)
	}
	return time.Now()
}

// Outa stops a clock that was started by `Into(..)`, and also writes
// the elapsed time (in a human-friendly format) to `os.Stderr`.
// The message is a bit more informative if the string starts with `!`
// (which is not printed).
func Outa(s string, t time.Time) {
	if S.HasPrefix(s, "!") {
		println(fmt.Sprintf("[>ENDED<] %s (elapsed %s)", s[1:], time.Since(t)))
	} else {
		println(fmt.Sprintf("==> %s took %s", s, time.Since(t)))
	}
}
