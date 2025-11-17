package domain

type User struct {
	Id         int    `json:"user_id"`
	Username   string `json:"user_name"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	EnableTOTP bool   `json:"enableTOTP"`
}
