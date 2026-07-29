package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	authv1 "nopresh.apetrovic.com/gen/proto/auth/v1"
	"nopresh.apetrovic.com/internal/data"
	"nopresh.apetrovic.com/internal/domain/settings"
	"nopresh.apetrovic.com/internal/domain/user"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type AuthServer struct {
	models data.Models
	logger *slog.Logger
	jwt    *auth.JWT
	config *config
}

func setJwtCookies(header http.Header, accessToken, refreshToken string) {
	header.Set("Set-Cookie", (&http.Cookie{
		Name:     "jwt",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   900,
	}).String())

	header.Add("Set-Cookie", (&http.Cookie{
		Name:     "refresh",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   2592000,
	}).String())
}

func (s *AuthServer) Register(
	ctx context.Context,
	req *connect.Request[authv1.RegisterRequest],
) (*connect.Response[authv1.RegisterResponse], error) {
	newUser := user.New(req.Msg.Email, req.Msg.Name)

	params := auth.DefaultParams()

	encodedHash, err := auth.GenerateFromPassword(req.Msg.Password, params)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	newUser, err = s.models.Users.Insert(ctx, newUser, encodedHash)

	if err != nil {
		return nil, err
	}

	newSettings := settings.Settings{
		UserId:   newUser.ID,
		TimeZone: s.config.tz,
	}

	_, err = s.models.Settings.Insert(ctx, &newSettings)

	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.CreateToken(newUser.ID, newUser.Name, newUser.Email)

	if err != nil {
		s.logger.Error("error during creation of jwt token ", "error", err.Error())
		return nil, err
	}

	refreshExpiry := time.Now().Add(30 * time.Hour)

	refreshTokenString, uuid, err := s.jwt.CreateRefreshToken(newUser.ID, newUser.Name, newUser.Email, refreshExpiry)

	if err != nil {
		s.logger.Error("error during creation of jwt token ", "error", err.Error())
		return nil, err
	}

	refreshToken, err := s.models.RefreshToken.Insert(refreshTokenString, uuid, newUser.Email, refreshExpiry)

	if err != nil {
		s.logger.Error("error during creation of the refresh token ", "error", err.Error())
		return nil, err
	}

	res := connect.NewResponse(&authv1.RegisterResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
	})

	setJwtCookies(res.Header(), accessToken, refreshTokenString)

	return res, nil
}

func (s *AuthServer) Login(
	_ context.Context,
	req *connect.Request[authv1.LoginRequest],
) (*connect.Response[authv1.LoginResponse], error) {
	user, err := s.models.Users.GetByEmail(req.Msg.Email)

	if err != nil {
		return nil, err
	}

	match, err := auth.ComparePasswordAndHash(req.Msg.Password, user.Hashed_pasword)

	if err != nil {
		return nil, err
	}

	if !match {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("Incorrect Credentials"))
	}

	accessToken, err := s.jwt.CreateToken(user.ID, user.Name, user.Email)

	if err != nil {
		s.logger.Error("error during creation of jwt token ", "error", err.Error())
		return nil, err
	}

	refreshExpiry := time.Now().Add(30 * 24 * time.Hour)

	refreshTokenString, uuid, err := s.jwt.CreateRefreshToken(user.ID, user.Name, user.Email, refreshExpiry)

	if err != nil {
		s.logger.Error("error during creation of jwt token ", "error", err.Error())
		return nil, err
	}

	refreshToken, err := s.models.RefreshToken.Insert(refreshTokenString, uuid, user.Email, refreshExpiry)

	if err != nil {
		s.logger.Error("error during creation of the refresh token ", "error", err.Error())
		return nil, err
	}

	res := connect.NewResponse(&authv1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.RefreshToken,
	})

	setJwtCookies(res.Header(), accessToken, refreshTokenString)

	return res, nil
}

func (s *AuthServer) Logout(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[emptypb.Empty], error) {
	info, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
	}

	err := s.models.RefreshToken.RevokeToken(info.RefreshClaims.RegisteredClaims.ID, info.RefreshClaims.Email)

	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
	}

	res := connect.NewResponse(&emptypb.Empty{})

	deleteCookie := func(name string) string {
		return (&http.Cookie{
			Name:     name,
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			Path:     "/",
			MaxAge:   -1,
		}).String()
	}

	res.Header().Add("Set-Cookie", deleteCookie("jwt"))
	res.Header().Add("Set-Cookie", deleteCookie("refresh"))

	return res, nil
}

func (s *AuthServer) Me(
	ctx context.Context,
	_ *connect.Request[emptypb.Empty],
) (*connect.Response[authv1.MeResponse], error) {
	info, ok := authn.GetInfo(ctx).(*auth.AuthInfo)

	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
	}

	return connect.NewResponse(&authv1.MeResponse{
		Name:  info.JwtClaims.Name,
		Email: info.JwtClaims.Email,
	}), nil
}
