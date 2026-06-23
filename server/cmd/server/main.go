package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/validate"
	greetv1 "nopresh.apetrovic.com/gen/greet/v1"
	"nopresh.apetrovic.com/gen/greet/v1/greetv1connect"
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	_ context.Context,
	req *greetv1.GreetRequest,
) (*greetv1.GreetResponse, error) {
	res := &greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Name),
	}
	return res, nil
}

func main() {
	api := http.NewServeMux()
	api.Handle(greetv1connect.NewGreetServiceHandler(&GreetServer{},
		connect.WithInterceptors(validate.NewInterceptor()),
	))

	mux := http.NewServeMux()

	mux.Handle("/api/", http.StripPrefix("/api", api))

	// path, handler := greetv1connect.NewGreetServiceHandler(
	// 	greeter,
	// 	// Validation via Protovalidate is almost always recommended
	// 	connect.WithInterceptors(validate.NewInterceptor()),
	// )
	// mux.Handle(path, handler)

	p := new(http.Protocols)

	p.SetHTTP1(true)

	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)
	s := http.Server{
		Addr:      "127.0.0.1:5000",
		Handler:   mux,
		Protocols: p,
	}

	err := s.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
