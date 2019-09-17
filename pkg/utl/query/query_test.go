package query_test

import (
	"testing"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"

	"github.com/labstack/echo"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/query"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	type args struct {
		user *rest.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *rest.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: &rest.AuthUser{
				Role: rest.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: &rest.AuthUser{
				Role:      rest.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &rest.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: &rest.AuthUser{
				Role:       rest.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &rest.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: &rest.AuthUser{
				Role: rest.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
