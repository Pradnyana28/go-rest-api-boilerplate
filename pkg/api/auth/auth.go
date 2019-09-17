package auth

import (
	"net/http"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"

	"github.com/labstack/echo"
)

// Custom errors
var (
	ErrInvalidCredentials = echo.NewHTTPError(http.StatusUnauthorized, "Username or password does not exist")
)

// Authenticate tries to authenticate the user provided by username and password
func (a *Auth) Authenticate(c echo.Context, user, pass string) (*rest.AuthToken, error) {
	u, err := a.udb.FindByUsername(a.db, user)
	if err != nil {
		return nil, err
	}

	if !a.sec.HashMatchesPassword(u.Password, pass) {
		return nil, ErrInvalidCredentials
	}

	if !u.Active {
		return nil, rest.ErrUnauthorized
	}

	token, expire, err := a.tg.GenerateToken(u)
	if err != nil {
		return nil, rest.ErrUnauthorized
	}

	u.UpdateLastLogin(a.sec.Token(token))

	if err := a.udb.Update(a.db, u); err != nil {
		return nil, err
	}

	return &rest.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}, nil
}

// Refresh refreshes jwt token and puts new claims inside
func (a *Auth) Refresh(c echo.Context, token string) (*rest.RefreshToken, error) {
	user, err := a.udb.FindByToken(a.db, token)
	if err != nil {
		return nil, err
	}
	token, expire, err := a.tg.GenerateToken(user)
	if err != nil {
		return nil, err
	}
	return &rest.RefreshToken{Token: token, Expires: expire}, nil
}

// Me returns info about currently logged user
func (a *Auth) Me(c echo.Context) (*rest.User, error) {
	au := a.rbac.User(c)
	return a.udb.View(a.db, au.ID)
}
