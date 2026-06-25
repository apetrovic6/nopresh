package main

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
	authv1 "nopresh.apetrovic.com/gen/proto/auth/v1"
	"nopresh.apetrovic.com/internal/data"
	user "nopresh.apetrovic.com/internal/domain/user"
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
	user := user.New(req.Msg.Email, req.Msg.Name)

	params := auth.DefaultParams()

	encodedHash, err := auth.GenerateFromPassword(req.Msg.Password, params)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	user, err = s.models.Users.Insert(user, encodedHash)

	if err != nil {
		return nil, err
	}

	token, err := s.jwt.CreateToken(user.ID, user.Name, user.Email)

	res := connect.NewResponse(&authv1.RegisterResponse{Token: token})

	res.Header().Set("Set-Cookie",
		"jwt="+token+"; HttpOnly; SameSite=Lax; Path=/")

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

	token, err := s.jwt.CreateToken(user.ID, user.Name, user.Email)

	if err != nil {
		s.logger.Error("error during creation of jwt token ", "error", err.Error())
		return nil, err
	}

	res := connect.NewResponse(&authv1.LoginResponse{Token: token})

	res.Header().Set("Set-Cookie",
		"jwt="+token+"; HttpOnly; SameSite=Lax; Path=/")

	return res, nil
}

func (s *AuthServer) Me(
	ctx context.Context,
	_ *connect.Request[emptypb.Empty],
) (*connect.Response[authv1.MeResponse], error) {
	claims, ok := authn.GetInfo(ctx).(*auth.UserClaims)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("missing auth info"))
	}

	return connect.NewResponse(&authv1.MeResponse{
		Name:  claims.Name,
		Email: claims.Email,
	}), nil
}
