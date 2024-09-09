package types

type UserDB struct {
	ID        int64  `db:"id"`
	Firstname string `db:"firstname"`
	Lastname  string `db:"lastname"`
	Username  string `db:"username"`
	Password  string `db:"password"`
}
