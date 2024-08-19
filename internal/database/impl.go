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
	GetUserById(id string) (resp *types.UserDB, err error)
	CreateUser(req types.CreateUser) (err error)
	UpdateUser(id string, req types.UpdateUser) (err error)
	DeleteUser(id string) (err error)
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

func (repo *Repository) GetUserById(id string) (resp *types.UserDB, err error) {
	var user types.UserDB
	err = repo.DB.QueryRow(getUserByIdQuery, id).Scan(&user.ID, &user.FirstName, &user.Username, &user.Phone)
	if err != nil {
		return nil, err
	}
	resp = &user

	return resp, nil
}

func (repo *Repository) CreateUser(req types.CreateUser) (err error) {
	_, err = repo.DB.Exec(createUserQuery, req.FirstName, req.Username, req.Phone)
	return err
}

func (repo *Repository) UpdateUser(id string, req types.UpdateUser) (err error) {
	_, err = repo.DB.Exec(updateUserQuery, req.FirstName, req.Username, req.Phone, id)
	return err
}

func (repo *Repository) DeleteUser(id string) (err error) {
	_, err = repo.DB.Exec(deleteUserQuery, id)
	return err
}
