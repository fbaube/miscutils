package miscutils

// Error is used to embed an error in a struct,
// when we want to execute (pointer) methods on
// a struct pipeline-style, i.e. chainable, and
// executed left-to-right.
//
// We make the error public so that it is easily
// set, so that we can wrap errors easily using
// the "%w" printf format spec.
//
type Errer struct {
	Err error
}

// HasError is a convenience function.
// Since Err is publicly visible, HasError is
// not really needed, but it seems appropriate
// given that we also have func Error() .
func (p Errer) HasError() bool {
	return p.Err != nil
}

// Error is NPE-proof.
func (p Errer) Error() string {
	if p.Err == nil { // !p.HasError() {
		return ""
	}
	return p.Err.Error()
}

/*
// SetError is unnecessary if Err is publicly visible.
func (p *Errable) SetError(e error) {
	p.Err = e
}
*/
