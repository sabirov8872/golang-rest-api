package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mockdatabase "github.com/sabirov8872/golang-rest-api/internal/database/mock"
	"github.com/sabirov8872/golang-rest-api/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestService_GetUserByUser(t *testing.T) {
	type args struct {
		username string
	}

	type want struct {
		userData *types.GetUserByUserDB
		err      error
	}

	tests := []struct {
		name         string
		args         args
		want         want
		mockBehavior func(*mockdatabase.MockIRepository, args)
	}{
		{
			name: "success case",
			args: args{
				username: "test",
			},
			want: want{
				userData: &types.GetUserByUserDB{
					ID:       1,
					Password: "test",
				},
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					GetUserByUser(a.username).
					Return(&types.GetUserByUserDB{
						ID:       1,
						Password: "test",
					}, nil)
			},
		},
		{
			name: "fail case",
			want: want{
				userData: nil,
				err:      errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					GetUserByUser(a.username).
					Return(nil, errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo, tt.args)

			serv := NewService(repo)
			a, err := serv.GetUserByUser(tt.args.username)
			assert.Equal(t, err, tt.want.err)
			assert.Equal(t, tt.want.userData, a)
		})
	}
}

func TestService_GetAllUsers(t *testing.T) {
	type want struct {
		userData *types.ListUserResponse
		err      error
	}

	tests := []struct {
		name         string
		want         want
		mockBehavior func(*mockdatabase.MockIRepository)
	}{
		{
			name: "success case",
			want: want{
				userData: &types.ListUserResponse{
					Items: []*types.User{
						{
							ID:        1,
							Firstname: "test",
							Lastname:  "test",
							Username:  "test",
							Password:  "test",
						},
					},
				},
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository) {
				repo.EXPECT().
					GetAllUsers().Return([]*types.UserDB{
					{
						ID:        1,
						Firstname: "test",
						Lastname:  "test",
						Username:  "test",
						Password:  "test",
					},
				}, nil)
			},
		},
		{
			name: "fail case",
			want: want{
				userData: nil,
				err:      errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository) {
				repo.EXPECT().
					GetAllUsers().
					Return(nil, errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo)

			serv := NewService(repo)
			res, err := serv.GetAllUsers()
			assert.Equal(t, res, tt.want.userData)
			assert.Equal(t, err, tt.want.err)

		})
	}
}

func TestService_GetUserById(t *testing.T) {
	type args struct {
		id string
	}

	type want struct {
		userData *types.User
		err      error
	}

	tests := []struct {
		name         string
		args         args
		want         want
		mockBehavior func(*mockdatabase.MockIRepository, args)
	}{
		{
			name: "success case",
			want: want{
				userData: &types.User{
					ID:        1,
					Firstname: "test",
					Lastname:  "test",
					Username:  "test",
					Password:  "test",
				},
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					GetUserByID(a.id).
					Return(&types.UserDB{
						ID:        1,
						Firstname: "test",
						Lastname:  "test",
						Username:  "test",
						Password:  "test",
					}, nil)
			},
		},
		{
			name: "fail case",
			want: want{
				userData: nil,
				err:      errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					GetUserByID(a.id).
					Return(nil, errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo, tt.args)

			serv := NewService(repo)
			a, err := serv.GetUserById(tt.args.id)
			assert.Equal(t, err, tt.want.err)
			assert.Equal(t, tt.want.userData, a)
		})
	}
}

func TestService_CreateUser(t *testing.T) {
	type args struct {
		req types.CreateUserRequest
	}

	type want struct {
		userData *types.CreateUserResponse
		err      error
	}

	tests := []struct {
		name         string
		args         args
		want         want
		mockBehavior func(*mockdatabase.MockIRepository, args)
	}{
		{
			name: "success case",
			args: args{
				req: types.CreateUserRequest{
					Firstname: "test",
					Lastname:  "test",
					Username:  "test",
					Password:  "test",
				},
			},
			want: want{
				userData: &types.CreateUserResponse{
					ID: 1,
				},
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					CreateUser(a.req).
					Return(int64(1), nil)
			},
		},
		{
			name: "fail case",
			want: want{
				userData: nil,
				err:      errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					CreateUser(a.req).
					Return(int64(0), errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo, tt.args)

			serv := NewService(repo)
			a, err := serv.CreateUser(tt.args.req)
			assert.Equal(t, err, tt.want.err)
			assert.Equal(t, tt.want.userData, a)
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	type args struct {
		id  string
		req types.UpdateUserRequest
	}

	type want struct {
		err error
	}

	tests := []struct {
		name         string
		args         args
		want         want
		mockBehavior func(*mockdatabase.MockIRepository, args)
	}{
		{
			name: "success case",
			args: args{
				id: "1",
				req: types.UpdateUserRequest{
					Firstname: "test",
					Lastname:  "test",
					Username:  "test",
					Password:  "test",
				},
			},
			want: want{
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					UpdateUser(a.id, a.req).
					Return(nil)
			},
		},
		{
			name: "fail case",
			want: want{
				err: errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					UpdateUser(a.id, a.req).
					Return(errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo, tt.args)

			serv := NewService(repo)
			err := serv.UpdateUser(tt.args.id, tt.args.req)
			assert.Equal(t, err, tt.want.err)
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	type args struct {
		id string
	}

	type want struct {
		err error
	}

	tests := []struct {
		name         string
		args         args
		want         want
		mockBehavior func(*mockdatabase.MockIRepository, args)
	}{
		{
			name: "success case",
			args: args{
				id: "1",
			},
			want: want{
				err: nil,
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					DeleteUser(a.id).
					Return(nil)
			},
		},
		{
			name: "fail case",
			want: want{
				err: errors.New("error"),
			},
			mockBehavior: func(repo *mockdatabase.MockIRepository, a args) {
				repo.EXPECT().
					DeleteUser(a.id).
					Return(errors.New("error"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockdatabase.NewMockIRepository(ctrl)
			tt.mockBehavior(repo, tt.args)

			serv := NewService(repo)
			err := serv.DeleteUser(tt.args.id)
			assert.Equal(t, err, tt.want.err)
		})
	}
}
