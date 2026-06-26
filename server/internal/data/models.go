package data

import "gorm.io/gorm"

type Models struct {
	Users        UserModel
	RefreshToken RefreshTokenModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:        UserModel{DB: db},
		RefreshToken: RefreshTokenModel{DB: db},
	}

}
