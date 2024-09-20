package service

import (
	"github.com/sabirov8872/golang-rest-api/internal/database"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Service struct {
	repo database.IRepository
}

type IService interface {
	SignIn(username string) (*types.SignInDB, error)
	GetAllUsers() (*types.ListUserResponse, error)
	GetUserById(id string) (*types.User, error)
	CreateUser(req types.CreateUserRequest) (*types.CreateUserResponse, error)
	UpdateUser(id string, req types.UpdateUserRequest) error
	DeleteUser(id string) error
}

func NewService(repo database.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SignIn(username string) (*types.SignInDB, error) {
	return s.repo.SignIn(username)
}

func (s *Service) GetAllUsers() (*types.ListUserResponse, error) {
	res, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	resp := make([]*types.User, len(res))
	for i, v := range res {
		resp[i] = &types.User{
			ID:        v.ID,
			Firstname: v.Firstname,
			Lastname:  v.Lastname,
			Username:  v.Username,
			Password:  v.Password,
		}
	}

	return &types.ListUserResponse{
		Items: resp,
	}, nil
}

func (s *Service) GetUserById(id string) (*types.User, error) {
	res, err := s.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return &types.User{
		ID:        res.ID,
		Firstname: res.Firstname,
		Lastname:  res.Lastname,
		Username:  res.Username,
		Password:  res.Password,
	}, nil
}

func (s *Service) CreateUser(req types.CreateUserRequest) (*types.CreateUserResponse, error) {
	id, err := s.repo.CreateUser(req)
	if err != nil {
		return nil, err
	}

	return &types.CreateUserResponse{
		ID: id,
	}, nil
}

func (s *Service) UpdateUser(id string, req types.UpdateUserRequest) error {
	return s.repo.UpdateUser(id, req)
}

func (s *Service) DeleteUser(id string) error {
	return s.repo.DeleteUser(id)
}
