package wns

import "fmt"

// ErrType is a type to indicate different WNS feed error types
type ErrType int

const (
	// ErrTypeNoNew is a type for an error when no need files are available
	// to be fetched from source
	ErrTypeNoNew ErrType = iota
	// ErrTypeTooFrequent is a type for an error when requests are made too
	// frequently(usually several times per 10seconds)
	ErrTypeTooFrequent
	// ErrTypeUnknown is a type for random errors
	ErrTypeUnknown
)

// APIError is an error wrapper to distinguish between known and unknown
// feed errors, so additional logic could be done on errors that are
// actually only informational messages concealed in error format
type APIError struct {
	Type ErrType
	Err  string
}

// Error returns an error message from type and error value
func (e *APIError) Error() string {
	var t string
	switch e.Type {
	case ErrTypeNoNew:
		t = "No new files available"
	case ErrTypeTooFrequent:
		t = "Too frequent requests"
	case ErrTypeUnknown:
		t = "Unknown error"
	}
	return fmt.Sprintf("%s (%s)", t, e.Err)
}
