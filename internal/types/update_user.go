package types

type UpdateUser struct {
	FirstName string `json:"firstName"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}
