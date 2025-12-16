package models

type User struct {
	ID           string `json:"id"`
	Username     string `json:"uname"`
	Email        string `json:"email"`
	Password     string `json:"-"`
	RefreshToken string `json:"-"`
}
