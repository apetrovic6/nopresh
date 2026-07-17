package settings

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nopresh.apetrovic.com/internal/domain/settings"
)

type SettingsModel struct {
	DB *gorm.DB
}

type SettingsDbo struct {
	gorm.Model
	TimeZone            string
	UserId              uint `gorm:"uniqueIndex"`
	DefaultMedicationId uint
}

func (s *SettingsDbo) TableName() string {
	return "settings"
}

func (s *SettingsModel) Insert(ctx context.Context, settings *settings.Settings) (*settings.Settings, error) {
	settingsDbo := toDbo(settings)

	err := gorm.G[SettingsDbo](s.DB).Create(ctx, settingsDbo)

	if err != nil {
		if strings.Contains(err.Error(), "idx_settings_user_id") {
			return nil, errors.New("Settings entry already exists")
		}

		return nil, err
	}

	return toDomain(settingsDbo), nil
}

func (s *SettingsModel) Update(ctx context.Context, userId uint, updateDto settings.UpdateDto) (*settings.Settings, error) {
	updates := map[string]any{}

	if updateDto.DefaultMedicationId != nil {
		updates["default_medication_id"] = *updateDto.DefaultMedicationId
	}

	if updateDto.TimeZone != nil {
		updates["time_zone"] = *updateDto.TimeZone
	}

	rowsAffected, err := gorm.G[SettingsDbo](s.DB).
		Where("user_id = ?", userId).
		Set(clause.Assignments(updates)).
		Update(ctx)

	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, errors.New("settings entry not found")
	}

	return nil, nil
}

func (s *SettingsModel) GetByUserId(userId uint) (*settings.Settings, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	settings, err := gorm.G[SettingsDbo](s.DB).Where("user_id = ?", userId).First(ctx)

	if err != nil {
		return nil, err
	}

	return toDomain(&settings), nil
}

func toDbo(settings *settings.Settings) *SettingsDbo {
	return &SettingsDbo{
		UserId:              settings.UserId,
		DefaultMedicationId: settings.DefaultMedicationId,
		TimeZone:            settings.TimeZone,
	}
}

func toDomain(s *SettingsDbo) *settings.Settings {
	return &settings.Settings{
		ID:                  s.ID,
		UserId:              s.UserId,
		DefaultMedicationId: s.DefaultMedicationId,
		TimeZone:            s.TimeZone,
	}
}
