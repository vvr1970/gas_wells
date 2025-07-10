package entity

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"` // "admin" или "user"
}
