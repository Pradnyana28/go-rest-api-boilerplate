package rbac

import (
	rest "github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
	"github.com/labstack/echo"
)

// New creates new RBAC service
func New() *Service {
	return &Service{}
}

// Service is RBAC application service
type Service struct{}

func checkBool(b bool) error {
	if b {
		return nil
	}
	return echo.ErrForbidden
}

// User returns user data stored in jwt token
func (s *Service) User(c echo.Context) *rest.AuthUser {
	id := c.Get("id").(int)
	locationID := c.Get("location_id").(int)
	user := c.Get("username").(string)
	email := c.Get("email").(string)
	role := c.Get("role").(rest.AccessRole)
	return &rest.AuthUser{
		ID:         id,
		Username:   user,
		LocationID: locationID,
		Email:      email,
		Role:       role,
	}
}

// EnforceRole authorizes request by AccessRole
func (s *Service) EnforceRole(c echo.Context, r rest.AccessRole) error {
	return checkBool(!(c.Get("role").(rest.AccessRole) > r))
}

// EnforceUser checks whether the request to change user data is done by the same user
func (s *Service) EnforceUser(c echo.Context, ID int) error {
	// TODO: Implement querying db and checking the requested user's company_id/location_id
	// to allow company/location admins to view the user
	if s.isAdmin(c) {
		return nil
	}
	return checkBool(c.Get("id").(int) == ID)
}

// EnforceLocation checks whether the request to change location data
// is done by the user belonging to the requested location
func (s *Service) EnforceLocation(c echo.Context, ID int) error {
	if err := s.EnforceRole(c, rest.LocationAdminRole); err != nil {
		return err
	}
	return checkBool((c.Get("location_id").(int) == ID))
}

func (s *Service) isAdmin(c echo.Context) bool {
	return !(c.Get("role").(rest.AccessRole) > rest.AdminRole)
}

// AccountCreate performs auth check when creating a new account
func (s *Service) AccountCreate(c echo.Context, roleID rest.AccessRole, locationID int) error {
	if err := s.EnforceLocation(c, locationID); err != nil {
		return err
	}
	return s.IsLowerRole(c, roleID)
}

// IsLowerRole checks whether the requesting user has higher role than the user it wants to change
// Used for account creation/deletion
func (s *Service) IsLowerRole(c echo.Context, r rest.AccessRole) error {
	return checkBool(c.Get("role").(rest.AccessRole) < r)
}
