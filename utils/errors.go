package utils

import (
	"fmt"
	"runtime"
)

var InternalError = NewConfHubError("Internal error %s")

type ConfHubError struct {
	format string
	params []interface{}
}

func (this ConfHubError) Error() string {
	return fmt.Sprintf(this.format, this.params...)
}

func NewConfHubError(format string) *ConfHubError {
	return &ConfHubError{format:format}
}

func (this *ConfHubError) Details(details ...interface{}) *ConfHubError {
	this.params = details
	return this
}

// MyCaller returns the caller of the function that called it :)
func MyCaller() string {

	// we get the callers as uintptrs - but we just need 1
	fpcs := make([]uintptr, 1)

	// skip 3 levels to get to the caller of whoever called Caller()
	n := runtime.Callers(4, fpcs)
	if n == 0 {
		return "n/a" // proper error her would be better
	}

	// get the info of the actual function that's in the pointer
	fun := runtime.FuncForPC(fpcs[0]-1)
	if fun == nil {
		return "n/a"
	}

	// return its name
	return fun.Name()
}