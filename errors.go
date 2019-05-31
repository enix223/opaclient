package opaclient

import (
	"errors"
)

var (
	// ErrNotFound Not found
	ErrNotFound = errors.New("Not found")
	// ErrServerError Server error
	ErrServerError = errors.New("Server error")
	// ErrBadRequest Bad request
	ErrBadRequest = errors.New("Bad request")
	// ErrWriteConflict Write conflict
	ErrWriteConflict = errors.New("Write conflict")
	// ErrStreamingNotImplemented streaming not implemented
	ErrStreamingNotImplemented = errors.New("streaming not implemented")
)
