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
	CreateUser(firstName, username, phone string) (err error)
	UpdateUser(id, newFirstName, newUsername, newPhone string) (err error)
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

func (repo *Repository) CreateUser(firstName, username, phone string) (err error) {
	_, err = repo.DB.Exec(createUserQuery, firstName, username, phone)
	return err
}

func (repo *Repository) UpdateUser(id, newFirstName, newUsername, newPhone string) (err error) {
	_, err = repo.DB.Exec(updateUserQuery, newFirstName, newUsername, newPhone, id)
	return err
}

func (repo *Repository) DeleteUser(id string) (err error) {
	_, err = repo.DB.Exec(deleteUserQuery, id)
	return err
}
