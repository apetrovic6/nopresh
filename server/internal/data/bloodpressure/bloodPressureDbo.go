package bp

import (
	"context"
	"time"

	"gorm.io/gorm"
	"nopresh.apetrovic.com/internal/data/user"
	"nopresh.apetrovic.com/internal/domain/bp"
	domain "nopresh.apetrovic.com/internal/domain/bp"
)

type BloodPressureModel struct {
	DB *gorm.DB
}

// If you need more than 8bit int, you probably dead
type BloodPressureDbo struct {
	gorm.Model
	DateTime  time.Time
	Systolic  uint8
	Diastolic uint8
	Pulse     uint8
	UserId    uint
	User      user.UserDbo
}

func (b *BloodPressureDbo) TableName() string {
	return "blood_pressure"
}

func (b BloodPressureModel) Insert(bp *domain.BloodPressure) (*domain.BloodPressure, error) {
	ctx := context.Background()

	ugala := new(bp)

	err := gorm.G[BloodPressureDbo](b.DB).Create(ctx, ugala)

	if err != nil {
		return nil, err
	}

	return toDomain(ugala), nil
}

func new(bp *bp.BloodPressure) *BloodPressureDbo {
	return &BloodPressureDbo{
		DateTime:  bp.DateTime,
		Systolic:  bp.Systolic,
		Diastolic: bp.Diastolic,
		Pulse:     bp.Pulse,
		UserId:    bp.UserId,
	}
}

func toDomain(bp *BloodPressureDbo) *domain.BloodPressure {
	return &domain.BloodPressure{
		ID:        bp.ID,
		UserId:    bp.UserId,
		DateTime:  bp.DateTime,
		Systolic:  bp.Systolic,
		Diastolic: bp.Diastolic,
		Pulse:     bp.Pulse,
	}
}
