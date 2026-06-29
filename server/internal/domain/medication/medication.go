package medication

import (
	medicationv1 "nopresh.apetrovic.com/gen/proto/medication/v1"
)

type Medication struct {
	ID                uint
	Name              string
	RecommendedDosage float32
	UserId            uint
	DosageMeasurment  Measurment
}

type Measurment string

const (
	NA Measurment = "N/A"
	Mg Measurment = "mg"
	G  Measurment = "g"
)

func New(id uint, name string, recommendedDosage float32, dosageMeasument Measurment, userId uint) *Medication {
	return &Medication{
		ID:                id,
		Name:              name,
		RecommendedDosage: recommendedDosage,
		UserId:            userId,
		DosageMeasurment:  dosageMeasument,
	}
}

func ToPbMeasurment(m Measurment) medicationv1.MEDICATIONMEAUSERMENT {
	switch m {
	case Mg:
		return medicationv1.MEDICATIONMEAUSERMENT_MEDICATIONMEAUSERMENT_MG
	case G:
		return medicationv1.MEDICATIONMEAUSERMENT_MEDICATIONMEAUSERMENT_G
	default:
		return medicationv1.MEDICATIONMEAUSERMENT_MEDICATIONMEAUSERMENT_UNSPECIFIED
	}
}

func ToDomainMeasurment(m medicationv1.MEDICATIONMEAUSERMENT) Measurment {
	switch m {
	case medicationv1.MEDICATIONMEAUSERMENT_MEDICATIONMEAUSERMENT_G:
		return Mg
	case medicationv1.MEDICATIONMEAUSERMENT_MEDICATIONMEAUSERMENT_MG:
		return G
	default:
		return NA
	}
}
