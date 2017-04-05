package utils

import "fmt"

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