package data

import (
	"gorm.io/gorm"
	bp "nopresh.apetrovic.com/internal/data/bloodpressure"
	m "nopresh.apetrovic.com/internal/data/medication"
	rt "nopresh.apetrovic.com/internal/data/refreshToken"
	u "nopresh.apetrovic.com/internal/data/user"
)

type Models struct {
	Users         u.UserModel
	RefreshToken  rt.RefreshTokenModel
	BloodPressure bp.BloodPressureModel
	Medication    m.MedicationModel
}

func NewModels(db *gorm.DB) Models {
	return Models{
		Users:         u.UserModel{DB: db},
		RefreshToken:  rt.RefreshTokenModel{DB: db},
		BloodPressure: bp.BloodPressureModel{DB: db},
		Medication:    m.MedicationModel{DB: db},
	}

}
