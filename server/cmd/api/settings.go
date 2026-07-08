package main

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	settingsv1 "nopresh.apetrovic.com/gen/proto/settings/v1"
	"nopresh.apetrovic.com/internal/data"
	"nopresh.apetrovic.com/internal/domain/settings"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type SettingsServer struct {
	models data.Models
	logger *slog.Logger
}

type SettReq[T any] = *connect.Request[T]
type SettRes[T any] = *connect.Response[T]

func (s *SettingsServer) CreateSettings(
	ctx context.Context,
	req SettReq[settingsv1.CreateSettingsRequest],
) (SettRes[settingsv1.CreateSettingsResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	settingsEntry := settings.New(
		0,
		userCtx.JwtClaims.ID,
		uint(req.Msg.DefaultMedicationId),
	)

	settingsEntry, err := s.models.Settings.Insert(settingsEntry)

	if err != nil {
		s.logger.Error("couldn't save user settings entry",
			"id", settingsEntry.ID,
			"userId", settingsEntry.UserId,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't save new settings entry for the user"))
	}

	res := connect.NewResponse(&settingsv1.CreateSettingsResponse{
		Settings: settingsFromDomainObject(settingsEntry),
	})

	return res, nil
}

func (s *SettingsServer) GetSettings(
	ctx context.Context,
	req SettReq[settingsv1.GetSettingsRequest],
) (SettRes[settingsv1.GetSettingsResponse], error) {
	userCtx, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, ConnErrMissingAuthInfo
	}

	settings, err := s.models.Settings.GetByUserId(userCtx.JwtClaims.ID)

	if err != nil {
		s.logger.Error("couldn't retrieve user settings entry",
			"userId", userCtx.JwtClaims.ID,
		)
		return nil, connect.NewError(connect.CodeInternal, errors.New("couldn't retrieve settings entry for the user"))
	}

	res := connect.NewResponse(&settingsv1.GetSettingsResponse{
		Settings: settingsFromDomainObject(settings),
	})

	return res, nil
}

func settingsFromDomainObject(s *settings.Settings) *settingsv1.Settings {
	return &settingsv1.Settings{
		DefaultMedicationId: uint32(s.DefaultMedicationId),
	}
}
