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
)

func (app *app) authenticate(_ context.Context, req *http.Request) (any, error) {
	allowList := map[string]struct{}{
		// Procedure constants are available in the generated code.
		authv1connect.AuthServiceRegisterProcedure:                       {},
		authv1connect.AuthServiceLoginProcedure:                          {},
		"/grpc.reflection.v1.ServerReflection/ServerReflectionInfo":      {},
		"/grpc.reflection.v1alpha.ServerReflection/ServerReflectionInfo": {},
	}

	token := ""
	procedure, _ := authn.InferProcedure(req.URL)

	if cookie, err := req.Cookie("jwt"); err == nil {
		token = cookie.Value
	} else if t, ok := authn.BearerToken(req); ok {
		token = t
	}

	if token == "" {
		if _, ok := allowList[procedure]; ok {
			return nil, nil
		}
		err := authn.Errorf("invalid authorization")
		err.Meta().Set("WWW-Authenticate", "Bearer")
		return nil, err
	}

	return app.jwt.VerifyToken(token)
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

	path, handler := greetv1connect.NewGreetServiceHandler(&GreetServer{},
		connect.WithInterceptors(validate.NewInterceptor()),
	)

	router.Handle(path, handler)

	middleware := authn.NewMiddleware(app.authenticate)

	// router.Handle(path, middleware.Wrap(handler))

	router.Handle(authv1connect.NewAuthServiceHandler(&AuthServer{models: app.models, logger: app.logger, jwt: app.jwt},
		connect.WithInterceptors(validate.NewInterceptor()),
	))

	wrapped := middleware.Wrap(router)

	return wrapped
}
