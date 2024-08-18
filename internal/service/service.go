package service

import (
	"github.com/sabirov8872/golang-rest-api/internal/database"
	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Service struct {
	repo database.IRepository
}

type IService interface {
	GetAllUsers() ([]*types.User, error)
	GetUserById(id string) (*types.User, error)
	CreateUser(firstName, username, phone string) error
	UpdateUser(id, firstName, username, phone string) error
	DeleteUser(id string) error
}

func NewService(repo database.IRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAllUsers() ([]*types.User, error) {
	res, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	resp := make([]*types.User, len(res))
	for i, v := range res {
		resp[i] = &types.User{
			ID:        v.ID,
			FirstName: v.FirstName,
			Username:  v.Username,
			Phone:     v.Phone,
		}
	}

	return resp, nil
}

func (s *Service) GetUserById(id string) (*types.User, error) {
	res, err := s.repo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	resp := &types.User{
		ID:        res.ID,
		FirstName: res.FirstName,
		Username:  res.Username,
		Phone:     res.Phone,
	}

	return resp, nil
}

func (s *Service) CreateUser(firstName, username, phone string) error {
	err := s.repo.CreateUser(firstName, username, phone)
	return err
}

func (s *Service) UpdateUser(id, newFirstName, newUsername, newPhone string) error {
	err := s.repo.UpdateUser(id, newFirstName, newUsername, newPhone)
	return err
}

func (s *Service) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	return err
}
