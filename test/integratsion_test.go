package test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sabirov8872/golang-rest-api/internal/types"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestPostgresContainer(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)
	defer postgresContainer.Terminate(ctx)

	Host, err := postgresContainer.Host(ctx)
	require.NoError(t, err)
	Port, err := postgresContainer.MappedPort(ctx, "5432")
	require.NoError(t, err)

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"testuser", "testpass", Host, Port.Port(), "testdb")

	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	_, err = db.Exec(
		`create table if not exists users (
			id serial primary key,
			firstname varchar(255) not null,
			lastname varchar(255) not null,
			username varchar(255) not null unique,
			password varchar(255) not null
		);`)
	require.NoError(t, err)

	var id int64
	err = db.QueryRow(
		`INSERT INTO users (firstname, lastname, username, password)
			VALUES ($1, $2, $3, $4) RETURNING id`,
		"John", "Doe", "john", "aaa").Scan(&id)

	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	rows, err := db.Query(`SELECT id, firstname, lastname, username, password FROM users`)
	require.NoError(t, err)
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		err = rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Password)
		require.NoError(t, err)
		users = append(users, user)
	}

	expected := []types.User{
		{
			ID:        1,
			Firstname: "John",
			Lastname:  "Doe",
			Username:  "john",
			Password:  "aaa",
		},
	}

	require.Equal(t, expected, users)
}
