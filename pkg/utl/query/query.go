package query

import (
	"github.com/labstack/echo"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"
)

// List prepares data for list queries
func List(u *rest.AuthUser) (*rest.ListQuery, error) {
	switch true {
	case u.Role <= rest.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == rest.CompanyAdminRole:
		return &rest.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == rest.LocationAdminRole:
		return &rest.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
