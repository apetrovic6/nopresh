package main

import (
	"connectrpc.com/grpcreflect"
	"nopresh.apetrovic.com/gen/greet/v1/greetv1connect"
	"nopresh.apetrovic.com/gen/proto/auth/v1/authv1connect"
	"nopresh.apetrovic.com/gen/proto/bloodpressure/v1/bloodpressurev1connect"
	"nopresh.apetrovic.com/gen/proto/medication/v1/medicationv1connect"
)

func (app *app) registerReflection() {

	reflector := grpcreflect.NewStaticReflector(
		greetv1connect.GreetServiceName,
		authv1connect.AuthServiceName,
		medicationv1connect.MedicationServiceName,
		bloodpressurev1connect.BloodPressureServiceName,
	)

	app.mux.Handle(grpcreflect.NewHandlerV1(reflector))
	app.mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
}
