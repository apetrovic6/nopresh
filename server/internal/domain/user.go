package domain

import (
	"time"
)

type User struct {
	ID             uint
	Name           string
	Email          string
	Hashed_pasword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Version        int
}

func New(email, name string) *User {
	return &User{
		Name:  name,
		Email: email,
	}
}
