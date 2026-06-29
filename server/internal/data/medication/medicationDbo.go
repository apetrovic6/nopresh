package medication

import (
	"context"

	"gorm.io/gorm"
	// "nopresh.apetrovic.com/internal/data/user"
	"nopresh.apetrovic.com/internal/domain/medication"
)

type MedicationModel struct {
	DB *gorm.DB
}

type MedicationDbo struct {
	gorm.Model
	Name              string `gorm:"uniqueIndex"`
	RecommendedDosage float32
	DosageMeasurment  medication.Measurment `sql:"type:enum('mg', 'g';default:'mg')"`
	UserId            uint
	// user.UserDbo
}

func (m *MedicationDbo) TableName() string {
	return "medications"
}

func (m *MedicationModel) Insert(med *medication.Medication) (*medication.Medication, error) {
	ctx := context.Background()

	medication := new(med)

	err := gorm.G[MedicationDbo](m.DB).Create(ctx, medication)

	if err != nil {
		return nil, err
	}

	return toDomain(medication), nil
}

func new(m *medication.Medication) *MedicationDbo {
	return &MedicationDbo{
		Name:              m.Name,
		RecommendedDosage: m.RecommendedDosage,
		DosageMeasurment:  m.DosageMeasurment,
		UserId:            m.UserId,
	}
}

func toDomain(m *MedicationDbo) *medication.Medication {
	return &medication.Medication{
		ID:                m.ID,
		Name:              m.Name,
		RecommendedDosage: m.RecommendedDosage,
		UserId:            m.UserId,
		DosageMeasurment:  m.DosageMeasurment,
	}
}
