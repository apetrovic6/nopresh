package settings

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
	"nopresh.apetrovic.com/internal/domain/settings"
)

type SettingsModel struct {
	DB *gorm.DB
}

type SettingsDbo struct {
	gorm.Model
	UserId              uint `gorm:"uniqueIndex"`
	DefaultMedicationId uint
}

func (s *SettingsDbo) TableName() string {
	return "settings"
}

func (s *SettingsModel) Insert(settings *settings.Settings) (*settings.Settings, error) {
	ctx := context.Background()

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
	}
}

func toDomain(s *SettingsDbo) *settings.Settings {
	return &settings.Settings{
		ID:                  s.ID,
		UserId:              s.UserId,
		DefaultMedicationId: s.DefaultMedicationId,
	}
}
