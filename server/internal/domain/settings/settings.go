package settings

import (
	"time"
)

type Settings struct {
	ID                  uint
	UserId              uint
	DefaultMedicationId uint
	TimeZone            string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func New(id, userId, medicationId uint, timeZone string) *Settings {
	return &Settings{
		ID:                  id,
		UserId:              userId,
		DefaultMedicationId: medicationId,
		TimeZone:            timeZone,
	}
}

type UpdateDto struct {
	DefaultMedicationId *uint
	TimeZone            *string
}
