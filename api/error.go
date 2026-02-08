// simple-bank\api\error.go

package api

import (
	"errors"

	"github.com/lib/pq"
)

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}
	return false
}

var (
	ErrMissingAuthorizationHeader   = errors.New("authorization header is missing")
	ErrInvalidAuthorizationHeader   = errors.New("authorization header format is invalid")
	ErrUnsupportedAuthorizationType = errors.New("unsupported authorization type")
)
