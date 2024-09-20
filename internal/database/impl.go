package database

import (
	"database/sql"

	"github.com/sabirov8872/golang-rest-api/internal/types"
)

type Repository struct {
	DB *sql.DB
}

type IRepository interface {
	SignIn(username string) (*types.SignInDB, error)
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

func (repo *Repository) SignIn(username string) (*types.SignInDB, error) {
	var s types.SignInDB
	err := repo.DB.QueryRow(signInQuery, username).Scan(&s.ID, &s.Password)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

func (repo *Repository) GetAllUsers() (resp []*types.UserDB, err error) {
	rows, err := repo.DB.Query(getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u types.UserDB
		err = rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username, &u.Password)
		if err != nil {
			return nil, err
		}

		resp = append(resp, &u)
	}

	return resp, nil
}

func (repo *Repository) GetUserByID(id string) (*types.UserDB, error) {
	var u types.UserDB
	err := repo.DB.QueryRow(getUserByIdQuery, id).Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Username, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *Repository) CreateUser(req types.CreateUserRequest) (int64, error) {
	var id int64
	err := repo.DB.QueryRow(createUserQuery, req.Firstname, req.Lastname, req.Username, req.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *Repository) UpdateUser(id string, req types.UpdateUserRequest) error {
	_, err := repo.DB.Query(updateUserQuery, req.Firstname, req.Lastname, req.Username, req.Password, id)
	return err
}

func (repo *Repository) DeleteUser(id string) error {
	_, err := repo.DB.Query(deleteUserQuery, id)
	return err
}
