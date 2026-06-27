package refreshToken

import "time"

type RefreshToken struct {
	ID           uint
	UserEamil    string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
	Uuid         string
}
