package mock

import (
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(*rest.User) (string, string, error)
}

// GenerateToken mock
func (j *JWT) GenerateToken(u *rest.User) (string, string, error) {
	return j.GenerateTokenFn(u)
}
