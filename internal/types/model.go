package types

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}

type ListUserResponse struct {
	Items []*User `json:"items"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}

type CreateUserResponse struct {
	ID int64 `json:"Id"`
}

type UpdateUserRequest struct {
	FirstName string `json:"firstName"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
}

type GetTokenResponse struct {
	Username string `json:"username"`
}

type TokenResponse struct {
	UserID int    `json:"userId"`
	Token  string `json:"token"`
}
