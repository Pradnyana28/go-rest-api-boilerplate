package mockdb

import (
	"github.com/go-pg/pg/orm"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, rest.User) (*rest.User, error)
	ViewFn           func(orm.DB, int) (*rest.User, error)
	FindByUsernameFn func(orm.DB, string) (*rest.User, error)
	FindByTokenFn    func(orm.DB, string) (*rest.User, error)
	ListFn           func(orm.DB, *rest.ListQuery, *rest.Pagination) ([]rest.User, error)
	DeleteFn         func(orm.DB, *rest.User) error
	UpdateFn         func(orm.DB, *rest.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr rest.User) (*rest.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (*rest.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (*rest.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (*rest.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *rest.ListQuery, p *rest.Pagination) ([]rest.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr *rest.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr *rest.User) error {
	return u.UpdateFn(db, usr)
}
