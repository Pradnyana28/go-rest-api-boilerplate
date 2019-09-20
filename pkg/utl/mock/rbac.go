package mock

import (
	rest "github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
	"github.com/labstack/echo"
)

// RBAC Mock
type RBAC struct {
	UserFn            func(echo.Context) *rest.AuthUser
	EnforceRoleFn     func(echo.Context, rest.AccessRole) error
	EnforceUserFn     func(echo.Context, int) error
	EnforceCompanyFn  func(echo.Context, int) error
	EnforceLocationFn func(echo.Context, int) error
	AccountCreateFn   func(echo.Context, rest.AccessRole, int) error
	IsLowerRoleFn     func(echo.Context, rest.AccessRole) error
}

// User mock
func (a *RBAC) User(c echo.Context) *rest.AuthUser {
	return a.UserFn(c)
}

// EnforceRole mock
func (a *RBAC) EnforceRole(c echo.Context, role rest.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a *RBAC) EnforceUser(c echo.Context, id int) error {
	return a.EnforceUserFn(c, id)
}

// EnforceCompany mock
func (a *RBAC) EnforceCompany(c echo.Context, id int) error {
	return a.EnforceCompanyFn(c, id)
}

// EnforceLocation mock
func (a *RBAC) EnforceLocation(c echo.Context, id int) error {
	return a.EnforceLocationFn(c, id)
}

// AccountCreate mock
func (a *RBAC) AccountCreate(c echo.Context, roleID rest.AccessRole, locationID int) error {
	return a.AccountCreateFn(c, roleID, locationID)
}

// IsLowerRole mock
func (a *RBAC) IsLowerRole(c echo.Context, role rest.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
