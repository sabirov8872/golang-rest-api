package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	mockservice "github.com/sabirov8872/golang-rest-api/internal/service/mock"
	"github.com/sabirov8872/golang-rest-api/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetUserByUser(t *testing.T) {
	tests := []struct {
		name               string
		inputBody          string
		inputUser          types.GetUserByUserRequest
		mockBehavior       func(*mockservice.MockIService, string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputBody: `{
					"username": "foo",
					"password": "foo"
				}`,
			inputUser: types.GetUserByUserRequest{
				Username: "foo",
				Password: "foo",
			},
			mockBehavior: func(repo *mockservice.MockIService, username string) {
				repo.EXPECT().
					GetUserByUser(username).
					Return(&types.GetUserByUserDB{
						ID:       1,
						Password: "$2a$04$hFScWULAoYvPAkfn4KVjWe/JKM4vYPlOZs.zbmK5jge6ry6dhdtqq",
					}, nil)
			},
			expectedStatusCode: 200,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"username": "bar","password": "bar"}`,
			inputUser: types.GetUserByUserRequest{
				Username: "bar",
				Password: "bar",
			},
			mockBehavior: func(repo *mockservice.MockIService, username string) {
				repo.EXPECT().
					GetUserByUser(username).
					Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo, tt.inputUser.Username)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/auth", hand.GetUserByUser).Methods("GET")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/auth",
				bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_GetAllUsers(t *testing.T) {
	tests := []struct {
		name                string
		mockBehavior        func(*mockservice.MockIService)
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(repo *mockservice.MockIService) {
				repo.EXPECT().
					GetAllUsers().
					Return(&types.ListUserResponse{
						Items: []*types.User{
							{
								ID:        1,
								Firstname: "foo",
								Lastname:  "foo",
								Username:  "foo",
								Password:  "foo",
							},
						},
					}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"items":[{"id":1,"firstname":"foo","lastname":"foo","username":"foo","password":"foo"}]}
`,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(repo *mockservice.MockIService) {
				repo.EXPECT().
					GetAllUsers().
					Return(nil, errors.New("user not found"))
			},
			expectedStatusCode: 500,
			expectedRequestBody: `{"message":"internal server error"}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/user", hand.GetAllUsers).Methods("GET")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/user", nil)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_GetUserById(t *testing.T) {
	tests := []struct {
		name                string
		id                  string
		mockBehavior        func(*mockservice.MockIService, string)
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			id:   "1",
			mockBehavior: func(repo *mockservice.MockIService, id string) {
				repo.EXPECT().
					GetUserById(id).
					Return(&types.User{
						ID:        1,
						Firstname: "foo",
						Lastname:  "foo",
						Username:  "foo",
						Password:  "foo",
					}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"id":1,"firstname":"foo","lastname":"foo","username":"foo","password":"foo"}
`,
		},
		{
			name: "Empty Fields",
			id:   "1",
			mockBehavior: func(repo *mockservice.MockIService, id string) {
				repo.EXPECT().
					GetUserById(id).
					Return(nil, errors.New("user not found"))
			},
			expectedStatusCode:  204,
			expectedRequestBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo, tt.id)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/user/{id}", hand.GetUserById).Methods("GET")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/user/1", nil)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name                string
		inputBody           string
		user                types.CreateUserRequest
		mockBehavior        func(*mockservice.MockIService, types.CreateUserRequest)
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"firstname":"foo","lastname":"foo","username":"foo","password":"foo"}`,
			user: types.CreateUserRequest{
				Firstname: "foo",
				Lastname:  "foo",
				Username:  "foo",
				Password:  "foo",
			},
			mockBehavior: func(repo *mockservice.MockIService, user types.CreateUserRequest) {
				repo.EXPECT().
					CreateUser(user).
					Return(&types.CreateUserResponse{
						ID: 1,
					}, nil)
			},
			expectedStatusCode: 200,
			expectedRequestBody: `{"Id":1}
`,
		},
		{
			name: "Empty Fields",
			mockBehavior: func(repo *mockservice.MockIService, user types.CreateUserRequest) {
				repo.EXPECT().
					CreateUser(user).
					Return(nil, errors.New("internal server error"))
			},
			expectedStatusCode: 500,
			expectedRequestBody: `{"message":"internal server error"}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo, tt.user)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/user", hand.CreateUser).Methods("POST")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedRequestBody, w.Body.String())
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name               string
		inputBody          string
		id                 string
		user               types.UpdateUserRequest
		mockBehavior       func(*mockservice.MockIService, string, types.UpdateUserRequest)
		expectedStatusCode int
	}{
		{
			name:      "OK",
			inputBody: `{"firstname":"foo","lastname":"foo","username":"foo","password":"foo"}`,
			id:        "1",
			user: types.UpdateUserRequest{
				Firstname: "foo",
				Lastname:  "foo",
				Username:  "foo",
				Password:  "foo",
			},
			mockBehavior: func(repo *mockservice.MockIService, id string, user types.UpdateUserRequest) {
				repo.EXPECT().
					UpdateUser(id, user).
					Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Empty Fields",
			id:   "1",
			mockBehavior: func(repo *mockservice.MockIService, id string, user types.UpdateUserRequest) {
				repo.EXPECT().
					UpdateUser(id, user).
					Return(errors.New("internal server error"))
			},
			expectedStatusCode: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo, tt.id, tt.user)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/user/{id}", hand.UpdateUser).Methods(http.MethodPut)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPut, "/user/1", bytes.NewBufferString(tt.inputBody))

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name               string
		id                 string
		mockBehavior       func(*mockservice.MockIService, string)
		expectedStatusCode int
	}{
		{
			name: "OK",
			id:   "1",
			mockBehavior: func(repo *mockservice.MockIService, id string) {
				repo.EXPECT().
					DeleteUser(id).
					Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			name: "Empty Fields",
			id:   "1",
			mockBehavior: func(repo *mockservice.MockIService, id string) {
				repo.EXPECT().
					DeleteUser(id).
					Return(errors.New("internal server error"))
			},
			expectedStatusCode: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := mockservice.NewMockIService(ctrl)
			tt.mockBehavior(repo, tt.id)
			hand := NewHandler(repo, "secretKey")

			router := mux.NewRouter()
			router.HandleFunc("/user/{id}", hand.DeleteUser).Methods(http.MethodDelete)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/user/1", nil)

			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedStatusCode, w.Code)
		})
	}
}
