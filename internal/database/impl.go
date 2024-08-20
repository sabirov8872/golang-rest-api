package database

import (
	"database/sql"

	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Repository struct {
	DB *sql.DB
}

type IRepository interface {
	GetAllUsers() (resp []*types.UserDB, err error)
	GetUserById(id string) (*types.UserDB, error)
	CreateUser(req types.CreateUser) (*types.UserDB, error)
	UpdateUser(id string, req types.UpdateUser) (*types.UserDB, error)
	DeleteUser(id string) (*types.UserDB, error)
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) GetAllUsers() (resp []*types.UserDB, err error) {
	rows, err := repo.DB.Query(getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user types.UserDB
		err = rows.Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &user)
	}

	return resp, nil
}

func (repo *Repository) GetUserById(id string) (*types.UserDB, error) {
	var user types.UserDB
	err := repo.DB.QueryRow(getUserByIdQuery, id).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *Repository) CreateUser(req types.CreateUser) (*types.UserDB, error) {
	var user types.UserDB
	err := repo.DB.QueryRow(createUserQuery, req.FirstName, req.Username, req.Phone).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *Repository) UpdateUser(id string, req types.UpdateUser) (*types.UserDB, error) {
	var user types.UserDB
	err := repo.DB.QueryRow(updateUserQuery, req.FirstName, req.Username, req.Phone, id).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *Repository) DeleteUser(id string) (*types.UserDB, error) {
	var user types.UserDB
	err := repo.DB.QueryRow(deleteUserQuery, id).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
