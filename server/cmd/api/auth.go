package main

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	authv1 "nopresh.apetrovic.com/gen/proto/auth/v1"
	"nopresh.apetrovic.com/internal/data"
	user "nopresh.apetrovic.com/internal/domain"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type AuthServer struct {
	models data.Models
	logger *slog.Logger
	jwt    *auth.JWT
}

func (s *AuthServer) Register(
	_ context.Context,
	req *connect.Request[authv1.RegisterRequest],
) (*connect.Response[authv1.RegisterResponse], error) {
	newUser := user.New(req.Msg.Email, req.Msg.Name)

	params := auth.DefaultParams()

	encodedHash, err := auth.GenerateFromPassword(req.Msg.Password, params)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	newUser, err = s.models.Users.Insert(newUser, encodedHash)

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

	res.Header().Set("Set-Cookie",
		"jwt="+accessToken+"; HttpOnly; SameSite=Lax; Path=/")
	res.Header().Add("Set-Cookie",
		"refresh="+refreshToken.RefreshToken+"; HttpOnly; SameSite=Lax; Path=/")

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

	refreshExpiry := time.Now().Add(30 * time.Hour)

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

	res.Header().Set("Set-Cookie",
		"jwt="+accessToken+"; HttpOnly; SameSite=Lax; Path=/")
	res.Header().Add("Set-Cookie",
		"refresh="+refreshToken.RefreshToken+"; HttpOnly; SameSite=Lax; Path=/")

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

	res.Header().Add("Set-Cookie", "jwt=; HttpOnly; SameSite=Lax; Path=/; Max-Age=0")
	res.Header().Add("Set-Cookie", "refresh=; HttpOnly; SameSite=Lax; Path=/; Max-Age=0")

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
