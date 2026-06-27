package bp

import "time"

type BloodPressure struct {
	ID        uint
	UserId    uint
	DateTime  time.Time
	Systolic  uint8
	Diastolic uint8
	Pulse     uint8
}

func New(id, userId uint, dateTime time.Time, systolic, diastolic, pulse uint8) *BloodPressure {
	return &BloodPressure{
		ID:        id,
		UserId:    userId,
		DateTime:  dateTime,
		Systolic:  systolic,
		Diastolic: diastolic,
		Pulse:     pulse,
	}
}
