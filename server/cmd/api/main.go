package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"connectrpc.com/grpcreflect"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nopresh.apetrovic.com/gen/greet/v1/greetv1connect"
	"nopresh.apetrovic.com/gen/proto/auth/v1/authv1connect"
	"nopresh.apetrovic.com/internal/data"
	"nopresh.apetrovic.com/internal/utils/auth"
)

type config struct {
	port        int
	env         string
	db          db
	jwtSecret   string
	jwtDuration time.Duration
}

type db struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type app struct {
	config config
	logger *slog.Logger
	models data.Models
	jwt    *auth.JWT
}

func main() {
	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "Evironment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", 5000, "5000")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("NOPRESH_DB_DSN"), "Postgres DSN")
	flag.StringVar(&cfg.jwtSecret, "JWT Secret Key", "supersecretkey", "Postgres DSN")
	flag.DurationVar(&cfg.jwtDuration, "JWT Token Duration", 15*time.Minute, "Postgres DSN")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := gorm.Open(postgres.Open(cfg.db.dsn), &gorm.Config{})

	if err != nil {
		logger.Error(err.Error())
	}

	db.AutoMigrate(&data.UserDbo{}, &data.RefreshTokenDbo{})

	app := &app{
		logger: logger,
		config: cfg,
		models: data.NewModels(db),
		jwt:    auth.NewJWT(cfg.jwtSecret, cfg.jwtDuration),
	}

	api := app.routes()

	reflector := grpcreflect.NewStaticReflector(
		greetv1connect.GreetServiceName,
		authv1connect.AuthServiceName,
	)

	mux := http.NewServeMux()
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle("/api/", http.StripPrefix("/api", api))

	p := new(http.Protocols)

	p.SetHTTP1(true)

	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	srv := http.Server{
		Addr:      "127.0.0.1:5000",
		Handler:   mux,
		Protocols: p,
	}

	app.logger.Info("starting server", "addr", srv.Addr) //	"env", app.config.env

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
