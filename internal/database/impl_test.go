package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sabirov8872/golang-rest-api/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetUserByUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		username string
	}

	type want struct {
		userdata *types.GetUserByUserDB
		err      error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args)
	}{
		{
			name: "success case",
			args: args{
				username: "test",
			},
			want: want{
				userdata: &types.GetUserByUserDB{
					ID:       1,
					Password: "test",
				},
				err: nil,
			},
			prepare: func(a args) {
				rows := sqlmock.NewRows([]string{"id", "username"}).
					AddRow(1, "test")

				mock.ExpectQuery(
					`SELECT id, password
						FROM users
						WHERE username = \$1`).
					WithArgs(a.username).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
			prepare: func(a args) {
				mock.ExpectQuery(
					`SELECT id, password
						FROM users
						WHERE username = \$1`).
					WithArgs(a.username).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			got, err := repo.GetUserByUser(tt.args.username)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_GetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type want struct {
		userdata []*types.UserDB
		err      error
	}

	tests := []struct {
		name    string
		want    want
		prepare func()
	}{
		{
			name: "success case",
			want: want{
				userdata: []*types.UserDB{
					{
						ID:        1,
						Firstname: "foo",
						Lastname:  "foo",
						Username:  "foo",
						Password:  "foo",
					},
					{
						ID:        2,
						Firstname: "test",
						Lastname:  "test",
						Username:  "test",
						Password:  "test",
					},
				},
				err: nil,
			},
			prepare: func() {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "username", "password"}).
					AddRow(1, "foo", "foo", "foo", "foo").
					AddRow(2, "test", "test", "test", "test")

				mock.ExpectQuery(
					`SELECT
						id,
						firstname,
						lastname,
						username,
						password
						FROM users`).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
			prepare: func() {
				mock.ExpectQuery(
					`SELECT
						id,
						firstname,
						lastname,
						username,
						password
						FROM users`).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()

			got, err := repo.GetAllUsers()
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		id string
	}

	type want struct {
		userdata *types.UserDB
		err      error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args)
	}{
		{
			name: "success case",
			args: args{
				id: "1",
			},
			want: want{
				userdata: &types.UserDB{
					ID:        1,
					Firstname: "foo",
					Lastname:  "foo",
					Username:  "foo",
					Password:  "foo",
				},
				err: nil,
			},
			prepare: func(a args) {
				rows := sqlmock.NewRows([]string{"id", "firstname", "lastname", "username", "password"}).
					AddRow(1, "foo", "foo", "foo", "foo")

				mock.ExpectQuery(
					`SELECT id,
							firstname,
							lastname,
							username,
							password
						FROM users
						WHERE id = \$1`).
					WithArgs(a.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
			prepare: func(a args) {
				mock.ExpectQuery(
					`SELECT id,
							    firstname,
							    lastname,
								username,
								password
							FROM users
							WHERE id = \$1`).
					WithArgs(a.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			got, err := repo.GetUserByID(tt.args.id)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		req types.CreateUserRequest
	}

	type want struct {
		id  int64
		err error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args)
	}{
		{
			name: "success case",
			args: args{
				req: types.CreateUserRequest{
					Firstname: "foo",
					Lastname:  "foo",
					Username:  "foo",
					Password:  "$2a$04$JI0ndi7QLiAjIsobT4KThOJjlaLjWTm2kpOw1.hFKmDGY0dmrPu/a",
				},
			},
			want: want{
				id:  1,
				err: nil,
			},
			prepare: func(a args) {
				rows := mock.NewRows([]string{"id"}).
					AddRow(1)

				mock.ExpectQuery(
					`INSERT INTO users \(firstname, 
                    		lastname, 
                    		username, 
                    		password\)
						VALUES \(\$1, \$2, \$3, \$4\)
						RETURNING id`).
					WithArgs(a.req.Firstname, a.req.Lastname, a.req.Username, a.req.Password).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				id:  0,
				err: sql.ErrNoRows,
			},
			prepare: func(a args) {
				mock.ExpectQuery(
					`INSERT INTO users \(firstname, 
                    		lastname, 
                    		username, 
                    		password\)
						VALUES \(\$1, \$2, \$3, \$4\)
						RETURNING id`).
					WithArgs(a.req.Firstname, a.req.Lastname, a.req.Username, a.req.Password).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			got, err := repo.CreateUser(tt.args.req)
			assert.Equal(t, tt.want.err, err)
			assert.Equal(t, tt.want.id, got)
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		id  string
		req types.UpdateUserRequest
	}

	type want struct {
		err error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args)
	}{
		{
			"success case",
			args{
				id: "1",
				req: types.UpdateUserRequest{
					Firstname: "foo",
					Lastname:  "foo",
					Username:  "foo",
					Password:  "foo",
				},
			},
			want{
				err: nil,
			},
			func(a args) {
				rows := mock.NewRows([]string{})

				mock.ExpectQuery(
					`UPDATE users 
						SET firstname = \$1, 
						    lastname = \$2, 
						    username = \$3, 
						    password = \$4 
						WHERE id = \$5`).
					WithArgs(a.req.Firstname, a.req.Lastname, a.req.Username, a.req.Password, a.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				err: sql.ErrNoRows,
			},
			prepare: func(a args) {
				mock.ExpectQuery(
					`UPDATE users
                    	SET firstname = \$1,
							lastname = \$2,
							username = \$3,
							password = \$4
						WHERE id = \$5`).
					WithArgs(a.req.Firstname, a.req.Lastname, a.req.Username, a.req.Password, a.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			err = repo.UpdateUser(tt.args.id, tt.args.req)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewRepository(db)

	type args struct {
		id string
	}
	type want struct {
		err error
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args)
	}{
		{
			name: "success case",
			args: args{
				id: "1",
			},
			want: want{
				err: nil,
			},
			prepare: func(a args) {
				rows := mock.NewRows([]string{})

				mock.ExpectQuery(
					`DELETE FROM users
						WHERE id = \$1`).
					WithArgs(a.id).
					WillReturnRows(rows)
			},
		},
		{
			name: "fail case",
			want: want{
				err: sql.ErrNoRows,
			},
			prepare: func(a args) {
				mock.ExpectQuery(
					`DELETE FROM users
						WHERE id = \$1`).
					WithArgs(a.id).
					WillReturnError(sql.ErrNoRows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(tt.args)

			err = repo.DeleteUser(tt.args.id)
			assert.Equal(t, tt.want.err, err)
		})
	}
}
