package main

import (
	"context"
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

func (app *app) extractClaims(accessToken, refreshToken string) (*auth.UserClaims, *auth.UserClaims, error) {
	accessClaims, accessErr := app.jwt.VerifyToken(accessToken)
	refreshClaims, refreshErr := app.jwt.VerifyToken(refreshToken)

	if refreshErr != nil {
		return nil, nil, auth.ErrCouldntVerifyRefreshToken
	}

	if accessErr != nil {
		return nil, refreshClaims, auth.ErrCouldntVerifyAccessToken
	}

	return accessClaims, refreshClaims, nil
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

	return WithCors(app.withTokenRefresh(wrapped))
}
