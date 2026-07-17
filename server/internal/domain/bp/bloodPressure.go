package bp

import (
	"time"

	"nopresh.apetrovic.com/internal/domain/medication"
)

type BloodPressure struct {
	ID              uint
	UserId          uint
	DateTimeUtc     time.Time
	Systolic        uint16
	Diastolic       uint16
	Pulse           uint16
	Dosage          float32
	MedicationId    uint
	Medication      medication.Medication
	MedicationTaken bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UpdateDto struct {
	DateTimeUtc     *time.Time
	Systolic        *uint16
	Diastolic       *uint16
	Pulse           *uint16
	Dosage          *float32
	MedicationId    *uint
	Medication      *medication.Medication
	MedicationTaken *bool
}

func New(id, userId uint, dateTime time.Time, systolic, diastolic, pulse uint16, dosage float32, medicationId uint, medicationTaken bool) *BloodPressure {
	return &BloodPressure{
		ID:              id,
		UserId:          userId,
		DateTimeUtc:     dateTime,
		Systolic:        systolic,
		Diastolic:       diastolic,
		Pulse:           pulse,
		Dosage:          dosage,
		MedicationId:    medicationId,
		MedicationTaken: medicationTaken,
	}
}
