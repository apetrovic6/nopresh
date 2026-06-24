package main

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
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
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {

	user := user.New(req.Email, req.Name)

	params := auth.DefaultParams()

	encodedHash, err := auth.GenerateFromPassword(req.Password, params)

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	user, err = s.models.Users.Insert(user, encodedHash)

	if err != nil {
		return nil, err
	}

	res := &authv1.RegisterResponse{Token: "ugala bugala"}

	return res, nil
}

func (s *AuthServer) Login(
	_ context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	user, err := s.models.Users.GetByEmail(req.Email)

	if err != nil {
		return nil, err
	}

	match, err := auth.ComparePasswordAndHash(req.Password, user.Hashed_pasword)

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

	res := &authv1.LoginResponse{Token: token}
	return res, nil
}
