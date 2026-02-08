// simple-bank\api\error.go

package api

import "github.com/lib/pq"

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name() == "unique_violation"
	}
	return false
}
