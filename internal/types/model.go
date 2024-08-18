package types

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}
