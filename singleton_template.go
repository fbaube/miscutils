package miscutils

import (
	"sync"
)

// singleton is TBS - feel free.
type singleton struct {
}

var theInstance *singleton
var once sync.Once

// It returns the single instance. Come to think of it tho, this
// looks like lazy instantiation; maybe it belongs in a `func init()` ?
func It() *singleton {
	once.Do(func() {
		theInstance = &singleton{}
	})
	return theInstance
}
