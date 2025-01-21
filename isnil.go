package miscutils

import "reflect"

// IsNil checks whether an interface value is nil,
// and to do this correctly, it checks BOTH values. 
func IsNil(x interface{}) bool {
  if x == nil {
    return true
  }
  return reflect.ValueOf(x).IsNil()
}

