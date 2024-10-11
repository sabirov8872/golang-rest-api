package database

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sabirov8872/golang-rest-api/internal/types"
	"github.com/stretchr/testify/require"
)

func TestRepository_SignIn(t *testing.T) {
	type args struct {
		username string
	}

	type want struct {
		userdata *types.GetUserByUserDB
		err      error
	}

	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args, *fields) error
	}{
		{
			name: "success case",
			args: args{
				username: "test",
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT id, password
						FROM users
						WHERE username = \$1`).
					WithArgs(args.username).
					WillReturnRows(
						mock.NewRows([]string{"id", "password"}).
							AddRow(1, "testpass"),
					)

				fields.db = db

				return err
			},
			want: want{
				userdata: &types.GetUserByUserDB{
					ID:       1,
					Password: "testpass",
				},
				err: nil,
			},
		},
		{
			name: "fail case",
			args: args{
				username: "bar",
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT id, password
						FROM users
						WHERE username = \$1`).
					WithArgs(args.username).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(tt.args, &ff))
			repo := NewRepository(ff.db)
			got, err := repo.GetUserByUser(tt.args.username)
			require.Equal(t, tt.want.err, err)
			require.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_GetAllUsers(t *testing.T) {
	type want struct {
		userdata []*types.UserDB
		err      error
	}

	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		want    want
		prepare func(*fields) error
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
						ID:        3,
						Firstname: "test",
						Lastname:  "test",
						Username:  "test",
						Password:  "test",
					},
				},
				err: nil,
			},
			prepare: func(fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT
						id,
						firstname,
						lastname,
						username,
						password
						FROM users`).
					WillReturnRows(
						mock.NewRows(
							[]string{"id", "firstname", "lastname", "username", "password"}).
							AddRow(1, "foo", "foo", "foo", "foo").
							AddRow(3, "test", "test", "test", "test"))

				fields.db = db

				return err
			},
		},
		{
			name: "fail case",
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
			prepare: func(fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT
						id,
						firstname,
						lastname,
						username,
						password
						FROM users`).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(&ff))
			repo := NewRepository(ff.db)
			got, err := repo.GetAllUsers()
			require.Equal(t, tt.want.err, err)
			require.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_GetUserByID(t *testing.T) {
	type args struct {
		id string
	}

	type want struct {
		userdata *types.UserDB
		err      error
	}

	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args, *fields) error
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
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT id,
							firstname,
							lastname,
							username,
							password
						FROM users
						WHERE id = \$1`).
					WithArgs(args.id).
					WillReturnRows(
						mock.NewRows([]string{"id", "firstname", "lastname", "username", "password"}).
							AddRow(1, "foo", "foo", "foo", "foo"))

				fields.db = db

				return err
			},
		},
		{
			name: "fail case",
			args: args{
				id: "2",
			},
			want: want{
				userdata: nil,
				err:      sql.ErrNoRows,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`SELECT id,
							    firstname,
							    lastname,
								username,
								password
							FROM users
							WHERE id = \$1`).
					WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(tt.args, &ff))
			repo := NewRepository(ff.db)
			got, err := repo.GetUserByID(tt.args.id)
			require.Equal(t, tt.want.err, err)
			require.Equal(t, tt.want.userdata, got)
		})
	}
}

func TestRepository_CreateUser(t *testing.T) {
	type args struct {
		req types.CreateUserRequest
	}

	type want struct {
		id  int64
		err error
	}

	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args, *fields) error
	}{
		{
			name: "success case",
			args: args{
				req: types.CreateUserRequest{
					Firstname: "foo",
					Lastname:  "foo",
					Username:  "foo",
					Password:  "foo",
				},
			},
			want: want{
				id:  1,
				err: nil,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`INSERT INTO users \(firstname, 
                    		lastname, 
                    		username, 
                    		password\)
						VALUES \(\$1, \$2, \$3, \$4\)
						RETURNING id`).
					WithArgs(args.req.Firstname, args.req.Lastname, args.req.Username, args.req.Password).
					WillReturnRows(
						mock.NewRows([]string{"id"}).
							AddRow(1))

				fields.db = db

				return err
			},
		},
		{
			name: "fail case",
			args: args{
				req: types.CreateUserRequest{
					Firstname: "bar",
					Lastname:  "bar",
					Username:  "bar",
					Password:  "bar",
				},
			},
			want: want{
				id:  0,
				err: sql.ErrNoRows,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`INSERT INTO users \(firstname, lastname, username, password\)
						VALUES \(\$1, \$2, \$3, \$4\)
						RETURNING id`).
					WithArgs(args.req.Firstname, args.req.Lastname, args.req.Username, args.req.Password).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(tt.args, &ff))
			repo := NewRepository(ff.db)
			got, err := repo.CreateUser(tt.args.req)
			require.Equal(t, tt.want.err, err)
			require.Equal(t, tt.want.id, got)
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	type args struct {
		id  string
		req types.UpdateUserRequest
	}
	type want struct {
		err error
	}
	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args, *fields) error
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
			func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`UPDATE users 
						SET firstname = \$1, 
						    lastname = \$2, 
						    username = \$3, 
						    password = \$4 
						WHERE id = \$5`).
					WithArgs(args.req.Firstname, args.req.Lastname, args.req.Username, args.req.Password, args.id).
					WillReturnRows(
						mock.NewRows([]string{}))

				fields.db = db

				return err
			},
		},
		{
			name: "fail case",
			args: args{},
			want: want{
				err: sql.ErrNoRows,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`UPDATE users
                    	SET firstname = \$1,
							lastname = \$2,
							username = \$3,
							password = \$4
						WHERE id = \$5`).
					WithArgs(args.req.Firstname, args.req.Lastname, args.req.Username, args.req.Password, args.id).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(tt.args, &ff))
			repo := NewRepository(ff.db)
			err := repo.UpdateUser(tt.args.id, tt.args.req)
			require.Equal(t, tt.want.err, err)
		})
	}
}

func TestRepository_DeleteUser(t *testing.T) {
	type args struct {
		id string
	}
	type want struct {
		err error
	}
	type fields struct {
		db *sql.DB
	}

	tests := []struct {
		name    string
		args    args
		want    want
		prepare func(args, *fields) error
	}{
		{
			name: "success case",
			args: args{
				id: "1",
			},
			want: want{
				err: nil,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`DELETE FROM users
						WHERE id = \$1`).
					WithArgs(args.id).
					WillReturnRows(
						mock.NewRows([]string{}))

				fields.db = db

				return err
			},
		},
		{
			name: "fail case",
			args: args{},
			want: want{
				err: sql.ErrNoRows,
			},
			prepare: func(args args, fields *fields) error {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
				}

				mock.ExpectQuery(
					`DELETE FROM users
						WHERE id = \$1`).
					WithArgs(args.id).
					WillReturnError(sql.ErrNoRows)

				fields.db = db

				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ff := fields{}
			require.NoError(t, tt.prepare(tt.args, &ff))
			repo := NewRepository(ff.db)
			err := repo.DeleteUser(tt.args.id)
			require.Equal(t, tt.want.err, err)
		})
	}
}
