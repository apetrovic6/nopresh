package data

import (
	"context"
	"time"

	"gorm.io/gorm"
	r "nopresh.apetrovic.com/internal/domain/refreshToken"
)

type RefreshTokenModel struct {
	DB *gorm.DB
}

type RefreshTokenDbo struct {
	gorm.Model
	UserEmail    string
	RefreshToken string `gorm:"size:256"`
	IsRevoked    bool
	ExpiresAt    time.Time
	Uuid         string
}

func (rt RefreshTokenDbo) TableName() string {
	return "refresh_tokens"
}

func (rt RefreshTokenModel) Insert(token, uuid, userEmail string, expiresAt time.Time) (*r.RefreshToken, error) {
	ctx := context.Background()

	refreshToken := New(token, uuid, userEmail, expiresAt)

	err := gorm.G[RefreshTokenDbo](rt.DB).Create(ctx, refreshToken)

	if err != nil {

		return nil, err
	}

	return toRefreshToken(*refreshToken), nil
}

func (rt RefreshTokenModel) RevokeToken(tokenId string, email string) error {
	ctx := context.Background()
	_, err := gorm.G[RefreshTokenDbo](rt.DB).Where("uuid = ? AND user_email = ?", tokenId, email).Update(ctx, "is_revoked", true)

	if err != nil {
		return err
	}

	return nil
}

func New(refreshToken, uuid, userEmail string, expiresAt time.Time) *RefreshTokenDbo {
	return &RefreshTokenDbo{
		RefreshToken: refreshToken,
		UserEmail:    userEmail,
		IsRevoked:    false,
		ExpiresAt:    expiresAt,
		Uuid:         uuid,
	}
}

func toRefreshToken(dbo RefreshTokenDbo) *r.RefreshToken {
	return &r.RefreshToken{
		ID:           dbo.ID,
		UserEamil:    dbo.UserEmail,
		RefreshToken: dbo.RefreshToken,
		IsRevoked:    dbo.IsRevoked,
		CreatedAt:    dbo.CreatedAt,
		ExpiresAt:    dbo.ExpiresAt,
		Uuid:         dbo.Uuid,
	}
}
