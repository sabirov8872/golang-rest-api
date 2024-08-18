package types

type UserDB struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	Username  string `db:"username"`
	Phone     string `db:"phone"`
}
