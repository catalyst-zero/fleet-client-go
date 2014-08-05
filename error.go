package client

import (
	"github.com/juju/errgo"
)

const (
	ERROR_TYPE_NOT_FOUND = 10000 + iota
)

type FleetClientError struct {
	StatusCode int
	StatusText string
}

func (this FleetClientError) Error() string {
	return this.StatusText
}

func NewFleetClientError(code int, text string) FleetClientError {
	return FleetClientError{
		StatusCode: code,
		StatusText: text,
	}
}

func IsNotFoundError(err error) bool {
	err = errgo.Cause(err)

	if appRepoErr, ok := err.(FleetClientError); ok {
		if appRepoErr.StatusCode == ERROR_TYPE_NOT_FOUND {
			return true
		}
	}

	return false
}
