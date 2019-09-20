package auth_test

import (
	"testing"
	"time"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/api/auth"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/mock"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/mock/mockdb"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"

	"github.com/go-pg/pg/orm"
	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	type args struct {
		user string
		pass string
	}
	cases := []struct {
		name     string
		args     args
		wantData *rest.AuthToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
		sec      *mock.Secure
	}{
		{
			name:    "Fail on finding user",
			args:    args{user: "juzernejm"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return nil, rest.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on wrong password",
			args:    args{user: "juzernejm", pass: "notHashedPassword"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return &rest.User{
						Username: user,
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return false
				},
			},
		},
		{
			name:    "Inactive user",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return &rest.User{
						Username: user,
						Password: "pass",
						Active:   false,
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return &rest.User{
						Username: user,
						Password: "pass",
						Active:   true,
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *rest.User) (string, string, error) {
					return "", "", rest.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on updating last login",
			args:    args{user: "juzernejm", pass: "pass"},
			wantErr: true,
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return &rest.User{
						Username: user,
						Password: "pass",
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u *rest.User) error {
					return rest.ErrGeneric
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *rest.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
		},
		{
			name: "Success",
			args: args{user: "juzernejm", pass: "pass"},
			udb: &mockdb.User{
				FindByUsernameFn: func(db orm.DB, user string) (*rest.User, error) {
					return &rest.User{
						Username: user,
						Password: "password",
						Active:   true,
					}, nil
				},
				UpdateFn: func(db orm.DB, u *rest.User) error {
					return nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *rest.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				TokenFn: func(string) string {
					return "refreshtoken"
				},
			},
			wantData: &rest.AuthToken{
				Token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires:      mock.TestTime(2000).Format(time.RFC3339),
				RefreshToken: "refreshtoken",
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, tt.sec, nil)
			token, err := s.Authenticate(nil, tt.args.user, tt.args.pass)
			if tt.wantData != nil {
				tt.wantData.RefreshToken = token.RefreshToken
				assert.Equal(t, tt.wantData, token)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
func TestRefresh(t *testing.T) {
	type args struct {
		c     echo.Context
		token string
	}
	cases := []struct {
		name     string
		args     args
		wantData *rest.RefreshToken
		wantErr  bool
		udb      *mockdb.User
		jwt      *mock.JWT
	}{
		{
			name:    "Fail on finding token",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*rest.User, error) {
					return nil, rest.ErrGeneric
				},
			},
		},
		{
			name:    "Fail on token generation",
			args:    args{token: "refreshtoken"},
			wantErr: true,
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*rest.User, error) {
					return &rest.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *rest.User) (string, string, error) {
					return "", "", rest.ErrGeneric
				},
			},
		},
		{
			name: "Success",
			args: args{token: "refreshtoken"},
			udb: &mockdb.User{
				FindByTokenFn: func(db orm.DB, token string) (*rest.User, error) {
					return &rest.User{
						Username: "username",
						Password: "password",
						Active:   true,
						Token:    token,
					}, nil
				},
			},
			jwt: &mock.JWT{
				GenerateTokenFn: func(u *rest.User) (string, string, error) {
					return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", mock.TestTime(2000).Format(time.RFC3339), nil
				},
			},
			wantData: &rest.RefreshToken{
				Token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
				Expires: mock.TestTime(2000).Format(time.RFC3339),
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, tt.jwt, nil, nil)
			token, err := s.Refresh(tt.args.c, tt.args.token)
			assert.Equal(t, tt.wantData, token)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestMe(t *testing.T) {
	cases := []struct {
		name     string
		wantData *rest.User
		udb      *mockdb.User
		rbac     *mock.RBAC
		wantErr  bool
	}{
		{
			name: "Success",
			rbac: &mock.RBAC{
				UserFn: func(echo.Context) *rest.AuthUser {
					return &rest.AuthUser{ID: 9}
				},
			},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (*rest.User, error) {
					return &rest.User{
						Base: rest.Base{
							ID:        id,
							CreatedAt: mock.TestTime(1999),
							UpdatedAt: mock.TestTime(2000),
						},
						FirstName: "John",
						LastName:  "Doe",
						Role: &rest.Role{
							AccessLevel: rest.UserRole,
						},
					}, nil
				},
			},
			wantData: &rest.User{
				Base: rest.Base{
					ID:        9,
					CreatedAt: mock.TestTime(1999),
					UpdatedAt: mock.TestTime(2000),
				},
				FirstName: "John",
				LastName:  "Doe",
				Role: &rest.Role{
					AccessLevel: rest.UserRole,
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.New(nil, tt.udb, nil, nil, tt.rbac)
			user, err := s.Me(nil)
			assert.Equal(t, tt.wantData, user)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestInitialize(t *testing.T) {
	a := auth.Initialize(nil, nil, nil, nil)
	if a == nil {
		t.Error("auth service not initialized")
	}
}
