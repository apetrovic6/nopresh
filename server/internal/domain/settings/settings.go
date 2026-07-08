package settings

import (
	"time"
)

type Settings struct {
	ID                  uint
	UserId              uint
	DefaultMedicationId uint
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func New(id, userId, medicationId uint) *Settings {
	return &Settings{
		ID:                  id,
		UserId:              userId,
		DefaultMedicationId: medicationId,
	}
}

type UpdateDto struct {
	DefaultMedicationId *uint
}
