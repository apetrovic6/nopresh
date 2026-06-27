package data

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	user "nopresh.apetrovic.com/internal/domain"
)

type UserModel struct {
	DB *gorm.DB
}

type UserDbo struct {
	gorm.Model
	Name            string `gorm:"size:64"`
	Email           string `gorm:"size:64;uniqueIndex"`
	Hashed_password string `gorm:"size:255"`
	Version         int
}

func (u UserDbo) TableName() string {
	return "users"
}

func (u UserModel) Insert(user *user.User, password []byte) (*user.User, error) {
	ctx := context.Background()

	userDbo := toDbo(user)

	userDbo.Hashed_password = string(password)

	err := gorm.G[UserDbo](u.DB).Create(ctx, userDbo)

	if err != nil {
		if strings.Contains(err.Error(), "idx_users_email") {
			return nil, errors.New("Couldn't register the user.")
		}

		return nil, err
	}

	return toUser(userDbo), nil
}

func (u UserModel) GetByEmail(email string) (*user.User, error) {
	ctx := context.Background()

	userDbo, err := gorm.G[UserDbo](u.DB).Where("email = ?", email).First(ctx)

	if err != nil {
		return nil, err
	}

	return toUser(&userDbo), nil
}

func toDbo(u *user.User) *UserDbo {
	return &UserDbo{
		Name:    u.Name,
		Email:   u.Email,
		Version: u.Version,
	}
}

func toUser(u *UserDbo) *user.User {
	return &user.User{
		ID:             u.ID,
		Name:           u.Name,
		Hashed_pasword: u.Hashed_password,
		Email:          u.Email,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Version:        u.Version,
	}
}
