package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"connectrpc.com/authn"
	"connectrpc.com/connect"
	"connectrpc.com/validate"

	greetv1 "nopresh.apetrovic.com/gen/greet/v1"
	"nopresh.apetrovic.com/gen/greet/v1/greetv1connect"
	"nopresh.apetrovic.com/gen/proto/auth/v1/authv1connect"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type AuthInfo struct {
	RefreshClaims *auth.UserClaims
	JwtClaims     *auth.UserClaims
	JWTToken      string
	RefreshToken  string
}

func extractTokens(req *http.Request) (string, string, error) {
	jwtToken := ""
	refreshToken := ""

	if cookie, err := req.Cookie("jwt"); err == nil {
		jwtToken = cookie.Value
	} else if t, ok := authn.BearerToken(req); ok {
		jwtToken = t
	}

	if jwtToken == "" {
		return "", "", authn.Errorf("invalid authorization")
	}

	if cookie, err := req.Cookie("refresh"); err == nil {
		refreshToken = cookie.Value
	}

	if cookie, err := req.Cookie("refresh"); err == nil {
		refreshToken = cookie.Value
	} else if t, ok := authn.BearerToken(req); ok {
		refreshToken = t
	}

	if refreshToken == "" {
		return "", "", authn.Errorf("invalid authorization")
	}

	return jwtToken, refreshToken, nil

}

func (app *app) extractClaims(accessToken, refreshToken string) (*auth.UserClaims, *auth.UserClaims, error) {
	accessClaims, err := app.jwt.VerifyToken(accessToken)

	if err != nil {
		return nil, nil, errors.New("couldn't verify access token")
	}

	refreshClaims, err := app.jwt.VerifyToken(refreshToken)

	if err != nil {
		return nil, nil, errors.New("couldn't verify refresh token")
	}

	return accessClaims, refreshClaims, nil

}

func (app *app) authenticate(_ context.Context, req *http.Request) (any, error) {
	allowList := map[string]struct{}{
		authv1connect.AuthServiceRegisterProcedure: {},
		authv1connect.AuthServiceLoginProcedure:    {},
		// authv1connect.AuthServiceLogoutProcedure:                         {},
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
	}

	procedure, _ := authn.InferProcedure(req.URL)

	// if !(procedure == authv1connect.AuthServiceLogoutProcedure) {
	if _, ok := allowList[procedure]; ok {
		return nil, nil
	}
	// }

	token, refreshToken, err := extractTokens(req)

	if err != nil {
		return nil, authn.Errorf("invalid authorization")
	}

	jwtClaims, refreshClaims, err := app.extractClaims(token, refreshToken)

	if err != nil {
		return nil, authn.Errorf("invalid authorization")
	}

	return &AuthInfo{
		RefreshClaims: refreshClaims,
		JwtClaims:     jwtClaims,
		JWTToken:      token,
		RefreshToken:  refreshToken,
	}, nil
}

type GreetServer struct{}

func (s *GreetServer) Greet(
	_ context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {

	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})

	return res, nil
}

func (app *app) routes() http.Handler {
	router := http.NewServeMux()

	router.Handle(greetv1connect.NewGreetServiceHandler(&GreetServer{},
		connect.WithInterceptors(validate.NewInterceptor()),
	))

	router.Handle(authv1connect.NewAuthServiceHandler(&AuthServer{models: app.models, logger: app.logger, jwt: app.jwt},
		connect.WithInterceptors(validate.NewInterceptor()),
	))

	middleware := authn.NewMiddleware(app.authenticate)
	wrapped := middleware.Wrap(router)

	return WithCors(wrapped)
}
