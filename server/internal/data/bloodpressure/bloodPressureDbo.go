package bp

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nopresh.apetrovic.com/internal/data/medication"
	"nopresh.apetrovic.com/internal/data/user"
	"nopresh.apetrovic.com/internal/domain/bp"
	domain "nopresh.apetrovic.com/internal/domain/bp"
	"nopresh.apetrovic.com/internal/utils/database"
)

type BloodPressureModel struct {
	DB *gorm.DB
}

// If you need more than 8bit int, you probably dead
type BloodPressureDbo struct {
	gorm.Model
	DateTimeUtc     time.Time
	Systolic        uint16
	Diastolic       uint16
	Pulse           uint16
	UserId          uint
	User            user.UserDbo
	Dosage          float32
	Comment         string
	MedicationId    uint
	Medication      medication.MedicationDbo
	MedicationTaken bool `gorm:"default:false"`
}

func (b *BloodPressureDbo) BeforeCreate(tx *gorm.DB) error {
	if b.DateTimeUtc.IsZero() {
		b.DateTimeUtc = time.Now().UTC()
	}
	return nil
}

func (b *BloodPressureDbo) TableName() string {
	return "blood_pressure"
}

func WithDateRange(keyset string, startDate time.Time, endDate time.Time) func(db *gorm.Statement) {
	return func(stmt *gorm.Statement) {
		if startDate.IsZero() || endDate.IsZero() {
			return
		}

		exprs := stmt.BuildCondition(keyset, startDate, endDate)

		stmt.AddClause(clause.Where{Exprs: exprs})
	}
}

// WithCursorPagination returns a gorm scope that appends a keyset predicate for
// the given cursor. A nil cursor (the first page) is a no-op, so it's always
// safe to pass through Scopes.
func WithCursorPagination(cursor *database.Cursor, keyset string) func(db *gorm.Statement) {
	return func(stmt *gorm.Statement) {
		if cursor == nil {
			return
		}

		exprs := stmt.BuildCondition(keyset, cursor.Date, cursor.Date, cursor.LastId)
		stmt.AddClause(clause.Where{Exprs: exprs})
	}
}

func (b BloodPressureModel) Get(ctx context.Context, userId, pageSize uint, cursor string, sortOrder database.SortOrder, startDate, endDate time.Time) ([]*domain.BloodPressure, *database.Cursor, error) {
	var decodedCursor *database.Cursor

	if cursor != "" {
		decoded, err := database.DecodeCursor(cursor)
		if err != nil {
			return nil, nil, err
		}

		decodedCursor = decoded
	}

	// Outer parens are required: the scope ANDs this onto the user_id filter, and
	// without them AND/OR precedence would break the user scoping.
	keyset := "(blood_pressure.date_time_utc < ? OR (blood_pressure.date_time_utc = ? AND blood_pressure.id < ?))"

	dateRangeFilter := "(blood_pressure.date_time_utc BETWEEN ? AND ?)"

	// Fetch one extra row to learn whether another page exists.
	entries, err := gorm.G[BloodPressureDbo](b.DB).
		Joins(clause.LeftJoin.Association("Medication"), nil).
		Where("blood_pressure.user_id = ?", userId).
		Scopes(WithCursorPagination(decodedCursor, keyset)).
		Scopes(WithDateRange(dateRangeFilter, startDate, endDate)).
		Order("blood_pressure.date_time_utc DESC, blood_pressure.id DESC").
		Limit(int(pageSize) + 1).
		Find(ctx)

	if err != nil {
		return nil, nil, err
	}

	var nextCursor *database.Cursor

	if len(entries) > int(pageSize) {
		entries = entries[:pageSize]    // drop the probe row
		last := entries[len(entries)-1] // cursor = last row we actually return
		nextCursor = database.CreateCursor(last.ID, last.DateTimeUtc, true)
	}

	domainEntries := lo.Map(entries, func(item BloodPressureDbo, index int) *domain.BloodPressure {
		return toDomain(&item)
	})

	return domainEntries, nextCursor, nil
}

func (b BloodPressureModel) Insert(ctx *context.Context, bp *domain.BloodPressure) (*domain.BloodPressure, error) {
	bpEntry := new(bp)

	err := gorm.G[BloodPressureDbo](b.DB).Create(*ctx, bpEntry)

	if err != nil {
		return nil, err
	}

	return toDomain(bpEntry), nil
}

func (b BloodPressureModel) Update(ctx *context.Context, bpId uint, bp *domain.UpdateDto, userId uint) error {
	updates := map[string]any{}

	if bp.Systolic != nil {
		updates["systolic"] = *bp.Systolic
	}

	if bp.Diastolic != nil {
		updates["diastolic"] = *bp.Diastolic
	}

	if bp.Pulse != nil {
		updates["pulse"] = *bp.Pulse
	}

	if bp.DateTimeUtc != nil {
		updates["date_time_utc"] = *bp.DateTimeUtc
	}

	if bp.MedicationId != nil {
		updates["medication_id"] = *bp.MedicationId
	}

	if bp.Dosage != nil {
		updates["dosage"] = *bp.Dosage
	}

	if bp.Comment != nil {
		updates["comment"] = *bp.Comment
	}

	if bp.MedicationTaken != nil {
		updates["medication_taken"] = *bp.MedicationTaken
	}

	rowsAffected, err := gorm.G[BloodPressureDbo](b.DB).
		Where("id = ? AND user_id = ?", bpId, userId).
		Set(clause.Assignments(updates)).
		Update(*ctx)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("entry not found")
	}

	return nil
}

func (b BloodPressureModel) Delete(ctx context.Context, bpId uint, userId uint) error {

	rowsAffected, err := gorm.G[BloodPressureDbo](b.DB).Where("id = ? AND user_id = ?", bpId, userId).Delete(ctx)

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("couldn't delete the blood pressure entry")
	}

	return nil
}

func new(bp *bp.BloodPressure) *BloodPressureDbo {
	return &BloodPressureDbo{
		DateTimeUtc:     bp.DateTimeUtc,
		Systolic:        bp.Systolic,
		Diastolic:       bp.Diastolic,
		Pulse:           bp.Pulse,
		UserId:          bp.UserId,
		Dosage:          bp.Dosage,
		Comment:         bp.Comment,
		MedicationId:    bp.MedicationId,
		MedicationTaken: bp.MedicationTaken,
	}
}

func toDomain(bp *BloodPressureDbo) *domain.BloodPressure {
	return &domain.BloodPressure{
		ID:              bp.ID,
		UserId:          bp.UserId,
		DateTimeUtc:     bp.DateTimeUtc,
		Systolic:        bp.Systolic,
		Diastolic:       bp.Diastolic,
		Pulse:           bp.Pulse,
		Dosage:          bp.Dosage,
		Comment:         bp.Comment,
		MedicationId:    bp.MedicationId,
		Medication:      *medication.ToDomain(&bp.Medication),
		MedicationTaken: bp.MedicationTaken,
		CreatedAt:       bp.CreatedAt,
		UpdatedAt:       bp.UpdatedAt,
	}
}
