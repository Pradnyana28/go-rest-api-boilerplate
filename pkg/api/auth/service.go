package auth

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/api/auth/platform/pgsql"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// New creates new iam service
func New(db *pg.DB, udb UserDB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return &Auth{
		db:   db,
		udb:  udb,
		tg:   j,
		sec:  sec,
		rbac: rbac,
	}
}

// Initialize initializes auth application service
func Initialize(db *pg.DB, j TokenGenerator, sec Securer, rbac RBAC) *Auth {
	return New(db, pgsql.NewUser(), j, sec, rbac)
}

// Service represents auth service interface
type Service interface {
	Authenticate(echo.Context, string, string) (*rest.AuthToken, error)
	Refresh(echo.Context, string) (*rest.RefreshToken, error)
	Me(echo.Context) (*rest.User, error)
}

// Auth represents auth application service
type Auth struct {
	db   *pg.DB
	udb  UserDB
	tg   TokenGenerator
	sec  Securer
	rbac RBAC
}

// UserDB represents user repository interface
type UserDB interface {
	View(orm.DB, int) (*rest.User, error)
	FindByUsername(orm.DB, string) (*rest.User, error)
	FindByToken(orm.DB, string) (*rest.User, error)
	Update(orm.DB, *rest.User) error
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*rest.User) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	HashMatchesPassword(string, string) bool
	Token(string) string
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) *rest.AuthUser
}
