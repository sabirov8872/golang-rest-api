package database

import (
	"database/sql"

	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Repository struct {
	DB *sql.DB
}

type IRepository interface {
	CheckUser(string) (int64, error)
	GetAllUsers() (resp []*types.UserDB, err error)
	GetUserByID(id string) (*types.UserDB, error)
	CreateUser(req types.CreateUserRequest) (int64, error)
	UpdateUser(id string, req types.UpdateUserRequest) error
	DeleteUser(id string) error
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) CheckUser(username string) (int64, error) {
	var id int64
	err := repo.DB.QueryRow(checkUserQuery, username).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
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

func (repo *Repository) GetUserByID(id string) (*types.UserDB, error) {
	var user types.UserDB
	err := repo.DB.QueryRow(getUserByIdQuery, id).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *Repository) CreateUser(req types.CreateUserRequest) (int64, error) {
	var id int64
	err := repo.DB.QueryRow(createUserQuery, req.FirstName, req.Username, req.Phone).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *Repository) UpdateUser(id string, req types.UpdateUserRequest) error {
	_, err := repo.DB.Query(updateUserQuery, req.FirstName, req.Username, req.Phone, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DeleteUser(id string) error {
	_, err := repo.DB.Query(deleteUserQuery, id)
	if err != nil {
		return err
	}

	return nil
}
