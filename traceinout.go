package miscutils

import (
	"fmt"
	"time"
)

// Note that these funcs add a dependency on package "time".
// It is also possible to write them using "interface{}".
// Note also that the caller must store the start time,
// but that this means we can have nested timers.

// Into starts the clock, and writes to os.Stderr.
func Into(s string) time.Time {
	println("[>BEGAN<]", s)
	return time.Now()
}

// Outa stops the clock, and also writes to os.Stderr.
func Outa(s string, t time.Time) {
	println(fmt.Sprintf("[>ENDED<] %s (elapsed %s) \n", s, time.Since(t)))
}
