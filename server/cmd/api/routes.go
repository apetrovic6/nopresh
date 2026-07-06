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
	"nopresh.apetrovic.com/gen/proto/bloodpressure/v1/bloodpressurev1connect"
	"nopresh.apetrovic.com/gen/proto/medication/v1/medicationv1connect"
	"nopresh.apetrovic.com/gen/proto/settings/v1/settingsv1connect"
)

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

	router.Handle(bloodpressurev1connect.NewBloodPressureServiceHandler(&BloodPressureServer{
		models: app.models,
		logger: app.logger,
	}))

	router.Handle(medicationv1connect.NewMedicationServiceHandler(&MedicationServer{
		models: app.models,
		logger: app.logger,
	}))

	router.Handle(settingsv1connect.NewSettingsServiceHandler(&SettingsServer{
		models: app.models,
		logger: app.logger,
	}))

	middleware := authn.NewMiddleware(app.middlewares.Authenticate)
	wrapped := middleware.Wrap(router)

	return app.middlewares.WithCors(app.middlewares.WithTokenRefresh(wrapped))
}

func (app *app) registerRoutes() {
	api := app.routes()
	app.mux.Handle("/api/", http.StripPrefix("/api", api))
}
