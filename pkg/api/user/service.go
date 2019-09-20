package user

import (
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/api/user/platform/pgsql"
	rest "github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, rest.User) (*rest.User, error)
	List(echo.Context, *rest.Pagination) ([]rest.User, error)
	View(echo.Context, int) (*rest.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, *Update) (*rest.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.NewUser(), rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, rest.User) (*rest.User, error)
	View(orm.DB, int) (*rest.User, error)
	List(orm.DB, *rest.ListQuery, *rest.Pagination) ([]rest.User, error)
	Update(orm.DB, *rest.User) error
	Delete(orm.DB, *rest.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *rest.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, rest.AccessRole, int) error
	IsLowerRole(echo.Context, rest.AccessRole) error
}
