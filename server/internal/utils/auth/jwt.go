package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserClaims struct {
	ID    uint
	Name  string
	Email string
	jwt.RegisteredClaims
}

// func NewUserClaims(id int, email string, duration time.Duration) (*UserClaims, error) {
// 	tokenID, err := uuid.NewRandom()

// 	if err != nil {
// 		return nil, fmt.Errorf("error generating token ID: %w", err)
// 	}

// 	return &UserClaims{
// 		ID:    id,
// 		Email: email,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Subject:   email,
// 			Audience:  jwt.ClaimStrings{},
// 			ExpiresAt: &jwt.NumericDate(time.Now().Add(duration)),
// 			IssuedAt:  &jwt.NumericDate(time.Now()),
// 			ID:        tokenID.String(),
// 		},
// 	}, nil
// }

type JWT struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJWT(secretKey string, tokenDuration time.Duration) *JWT {
	return &JWT{secretKey, tokenDuration}
}

func (j *JWT) CreateToken(id uint, name, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		ID:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(j.secretKey))
}

func (j *JWT) CreateRefreshToken(id uint, name, email string, expiresAt time.Time) (string, string, error) {
	uuid, err := uuid.NewV6()

	if err != nil {
		return "", "", errors.New("couldn't generate uuid for refresh token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		ID:    id,
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))

	if err != nil {
		return "", "", errors.New("couldn't generate signed token string")
	}

	return tokenString, uuid.String(), nil
}

func (j *JWT) VerifyToken(tokenStr string) (*UserClaims, error) {
	var claims UserClaims

	_, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("couldn't verify token: %w", err)
	}

	return &claims, nil
}
