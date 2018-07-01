package miscutils

import (
	"fmt"
	S "strings"
	"time"
)

// Note that these funcs add a dependency on package "time".
// It is also possible to write them using "interface{}".
// Note also that the caller must store the start time,
// but that this means we can have nested timers.

// Into starts the clock, and writes to os.Stderr.
func Into(s string) time.Time {
	if s != "" {
		println("[>BEGAN<]", s)
	}
	return time.Now()
}

// Outa stops the clock, and also writes to os.Stderr.
func Outa(s string, t time.Time) {
	if S.HasPrefix(s, "!") {
		println(fmt.Sprintf("[>ENDED<] %s (elapsed %s)", s[1:], time.Since(t)))
	} else {
		println(fmt.Sprintf("==> %s took %s", s, time.Since(t)))
	}
}
