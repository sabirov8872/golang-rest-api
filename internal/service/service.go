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
	CreateUser(req types.CreateUser) error
	UpdateUser(id string, req types.UpdateUser) error
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

func (s *Service) CreateUser(req types.CreateUser) error {
	err := s.repo.CreateUser(req)
	return err
}

func (s *Service) UpdateUser(id string, req types.UpdateUser) error {
	err := s.repo.UpdateUser(id, req)
	return err
}

func (s *Service) DeleteUser(id string) error {
	err := s.repo.DeleteUser(id)
	return err
}
