package models

type User struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balans    float64 `json:"balans"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type CreateUser struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balans    int `json:"balans"`
}

type UpdateUser struct {
	Id        string  `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Balans    float64 `json:"balans"`
}
type UserPrimaryKey struct {
	Id string `json:"id"`
}

type UserGetListRequest struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Search string `json:"search"`
}
type UserGetListResponse struct {
	Count int     `json:"count"`
	Users  []*User `json:"users"`
}
