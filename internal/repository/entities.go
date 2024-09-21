package repository

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

func (u *User) ToServiceUser() models.User {
	return models.User{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
	}
}

func FromServiceUser(user models.User, id int) User {
	return User{
		ID:       id,
		Email:    user.Email,
		Password: user.Password,
		Username: user.Username,
	}
}
