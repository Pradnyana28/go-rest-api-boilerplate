package pgsql

import (
	"github.com/go-pg/pg/orm"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// NewUser returns a new user database instance
func NewUser() *User {
	return &User{}
}

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u *User) View(db orm.DB, id int) (*rest.User, error) {
	user := &rest.User{Base: rest.Base{ID: id}}
	err := db.Select(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Update updates user's info
func (u *User) Update(db orm.DB, user *rest.User) error {
	return db.Update(user)
}
