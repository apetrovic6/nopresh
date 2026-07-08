package medication

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	// "nopresh.apetrovic.com/internal/data/user"
	"nopresh.apetrovic.com/internal/domain/medication"
)

type MedicationModel struct {
	DB *gorm.DB
}

type MedicationDbo struct {
	gorm.Model
	Name              string
	RecommendedDosage float32
	DosageMeasurement medication.Measurement `sql:"type:enum('mg', 'g';default:'mg')"`
	UserId            uint
}

func (m *MedicationDbo) TableName() string {
	return "medications"
}

func (m *MedicationModel) Insert(ctx context.Context, med *medication.Medication) (*medication.Medication, error) {
	medicationDbo := toDbo(med)

	err := gorm.G[MedicationDbo](m.DB).Create(ctx, medicationDbo)

	if err != nil {
		return nil, err
	}

	return ToDomain(medicationDbo), nil
}

func (m *MedicationModel) GetAll(ctx context.Context, userId uint) ([]*medication.Medication, error) {
	medications, err := gorm.G[MedicationDbo](m.DB).Where("user_id = ?", userId).Order("name asc").Find(ctx)

	if err != nil {
		return nil, err
	}

	medicationsDomain := lo.Map(medications, func(item MedicationDbo, index int) *medication.Medication {
		return ToDomain(&item)
	})

	return medicationsDomain, nil
}

func (m *MedicationModel) GetById(ctx context.Context, id, userId uint) (*medication.Medication, error) {
	medicationDbo, err := gorm.G[MedicationDbo](m.DB).Where("id = ? AND user_id = ?", id, userId).First(ctx)

	if err != nil {
		return nil, err
	}

	return ToDomain(&medicationDbo), err
}

func (m *MedicationModel) Update(ctx context.Context, medicationId uint, input medication.UpdateInput, userId uint) (*medication.Medication, error) {
	updates := map[string]any{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.RecommendedDosage != nil {
		updates["recommended_dosage"] = *input.RecommendedDosage
	}
	if input.DosageMeasurement != nil {
		updates["dosage_measurement"] = *input.DosageMeasurement
	}

	rowsAffected, err := gorm.G[MedicationDbo](m.DB).
		Where("id = ? AND user_id = ?", medicationId, userId).
		Set(clause.Assignments(updates)).
		Update(ctx)

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("medication not found")
	}

	medicationDbo, err := gorm.G[MedicationDbo](m.DB).Where("id = ? AND user_id = ?", medicationId, userId).First(ctx)

	if err != nil {
		return nil, err
	}

	return ToDomain(&medicationDbo), nil
}

func (m *MedicationModel) DeleteById(ctx context.Context, id, userId uint) (*medication.Medication, error) {
	medication, err := gorm.G[MedicationDbo](m.DB).Where("id = ? AND user_id = ?", id, userId).First(ctx)
	rowsAffected, err := gorm.G[MedicationDbo](m.DB).Where("id = ? AND user_id = ?", id, userId).Delete(ctx)

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}

	return ToDomain(&medication), nil
}

func toDbo(m *medication.Medication) *MedicationDbo {
	return &MedicationDbo{
		Name:              m.Name,
		RecommendedDosage: m.RecommendedDosage,
		DosageMeasurement: m.DosageMeasurement,
		UserId:            m.UserId,
	}
}

func ToDomain(m *MedicationDbo) *medication.Medication {
	return &medication.Medication{
		ID:                m.ID,
		Name:              m.Name,
		RecommendedDosage: m.RecommendedDosage,
		UserId:            m.UserId,
		DosageMeasurement: m.DosageMeasurement,
	}
}
