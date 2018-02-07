package miscutils

import (
	"sync"
)

// singleton is TBS.
type singleton struct {
}

var theInstance *singleton
var once sync.Once

// It returns the single instance. This looks like lazy
// instantiation; shouldn't it be done in func init() ?
func It() *singleton {
	once.Do(func() {
		theInstance = &singleton{}
	})
	return theInstance
}
