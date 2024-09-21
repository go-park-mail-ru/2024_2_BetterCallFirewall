package repository

import (
	"github.com/2024_2_BetterCallFirewall/internal/models"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (u *User) ToServiceUser() models.User {
	return models.User{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  u.Password,
		Email:     u.Email,
	}
}

func FromServiceUser(user models.User, id int) User {
	return User{
		ID:        id,
		Email:     user.Email,
		Password:  user.Password,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}
