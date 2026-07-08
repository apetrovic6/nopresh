package medication

import (
	medicationv1 "nopresh.apetrovic.com/gen/proto/medication/v1"
)

type Medication struct {
	ID                uint
	Name              string
	RecommendedDosage float32
	UserId            uint
	DosageMeasurement Measurement
}

// UpdateInput describes a partial update of a Medication. A nil field means
// "leave this column unchanged"; a non-nil field is applied. The transport layer
// populates it from a protobuf field mask; the data layer consumes it.
type UpdateInput struct {
	Name              *string
	RecommendedDosage *float32
	DosageMeasurement *Measurement
}

type Measurement string

const (
	NA Measurement = "N/A"
	Mg Measurement = "mg"
	G  Measurement = "g"
)

func New(id uint, name string, recommendedDosage float32, dosageMeasument Measurement, userId uint) *Medication {
	return &Medication{
		ID:                id,
		Name:              name,
		RecommendedDosage: recommendedDosage,
		UserId:            userId,
		DosageMeasurement: dosageMeasument,
	}
}

func ToDomainMeasurement(m medicationv1.MEDICATIONMEAUSEREMENT) Measurement {
	switch m {
	case medicationv1.MEDICATIONMEAUSEREMENT_MEDICATIONMEAUSEREMENT_MG:
		return Mg
	case medicationv1.MEDICATIONMEAUSEREMENT_MEDICATIONMEAUSEREMENT_G:
		return G
	default:
		return NA
	}
}
