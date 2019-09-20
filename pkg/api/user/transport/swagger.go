package transport

import (
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*rest.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []rest.User `json:"users"`
		Page  int          `json:"page"`
	}
}
