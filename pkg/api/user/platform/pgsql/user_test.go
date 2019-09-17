package pgsql_test

import (
	"testing"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/model"

	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/api/user/platform/pgsql"
	"github.com/Pradnyana28/go-rest-api-boilerplate/pkg/utl/mock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      rest.User
		wantData *rest.User
	}{
		{
			name:    "User already exists",
			wantErr: true,
			req: rest.User{
				Email:    "johndoe@mail.com",
				Username: "johndoe",
			},
		},
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: rest.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: rest.Base{
					ID: 1,
				},
			},
		},
		{
			name: "Success",
			req: rest.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: rest.Base{
					ID: 2,
				},
			},
			wantData: &rest.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: rest.Base{
					ID: 2,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &rest.Role{}, &rest.User{})

	if err := mock.InsertMultiple(db, &rest.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].req); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.Create(db, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if resp == nil {
					t.Error("Expected data, but received nil.")
					return
				}
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.UpdatedAt = resp.UpdatedAt
				assert.Equal(t, tt.wantData, resp)
			}
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData *rest.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: &rest.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: rest.Base{
					ID: 2,
				},
				Role: &rest.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &rest.Role{}, &rest.User{})

	if err := mock.InsertMultiple(db, &rest.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				if user == nil {
					t.Errorf("response was nil due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					tt.wantData.UpdatedAt = user.UpdatedAt
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *rest.User
		wantData *rest.User
	}{
		{
			name: "Success",
			usr: &rest.User{
				Base: rest.Base{
					ID: 2,
				},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			wantData: &rest.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Z",
				LastName:   "Freak",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Address:    "Address",
				Phone:      "123456",
				Mobile:     "345678",
				Base: rest.Base{
					ID: 2,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &rest.Role{}, &rest.User{})

	if err := mock.InsertMultiple(db, &rest.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				user := &rest.User{
					Base: rest.Base{
						ID: tt.usr.ID,
					},
				}
				if err := db.Select(user); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.LastLogin = user.LastLogin
				tt.wantData.DeletedAt = user.DeletedAt
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *rest.ListQuery
		pg       *rest.Pagination
		wantData []rest.User
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: &rest.Pagination{
				Limit: -100,
			},
		},
		{
			name: "Success",
			pg: &rest.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &rest.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			wantData: []rest.User{
				{
					Email:      "tomjones@mail.com",
					FirstName:  "Tom",
					LastName:   "Jones",
					Username:   "tomjones",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "newPass",
					Base: rest.Base{
						ID: 2,
					},
					Role: &rest.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
				},
				{
					Email:      "johndoe@mail.com",
					FirstName:  "John",
					LastName:   "Doe",
					Username:   "johndoe",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Base: rest.Base{
						ID: 1,
					},
					Role: &rest.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
					Token: "loginrefresh",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &rest.Role{}, &rest.User{})

	if err := mock.InsertMultiple(db, &rest.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := udb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, users)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      *rest.User
		wantData *rest.User
	}{
		{
			name: "Success",
			usr: &rest.User{
				Base: rest.Base{
					ID:        2,
					DeletedAt: mock.TestTime(2018),
				},
			},
			wantData: &rest.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: rest.Base{
					ID: 2,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &rest.Role{}, &rest.User{})

	if err := mock.InsertMultiple(db, &rest.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.NewUser()

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// Check if the deleted_at was set
		})
	}
}
