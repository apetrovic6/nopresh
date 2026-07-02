package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"nopresh.apetrovic.com/internal/data"
	bp "nopresh.apetrovic.com/internal/data/bloodpressure"
	"nopresh.apetrovic.com/internal/data/medication"
	rt "nopresh.apetrovic.com/internal/data/refreshToken"
	"nopresh.apetrovic.com/internal/data/settings"
	user "nopresh.apetrovic.com/internal/data/user"
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
	mux         *http.ServeMux
	config      config
	logger      *slog.Logger
	models      data.Models
	middlewares Middlewares
	jwt         *auth.JWT
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

	err = db.AutoMigrate(&user.UserDbo{}, &bp.BloodPressureDbo{}, &rt.RefreshTokenDbo{}, &medication.MedicationDbo{}, &settings.SettingsDbo{})

	if err != nil {
		logger.Error("couldn't migrate db")
		panic(err)
	}

	jwt := auth.NewJWT(cfg.jwtSecret, cfg.jwtDuration)
	app := &app{
		mux:    http.NewServeMux(),
		logger: logger,
		config: cfg,
		middlewares: Middlewares{
			jwt: jwt,
		},
		models: data.NewModels(db),
		jwt:    jwt,
	}

	app.registerReflection()
	app.registerRoutes()

	p := new(http.Protocols)

	p.SetHTTP1(true)

	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	srv := http.Server{
		Addr:      "127.0.0.1:5000",
		Handler:   app.mux,
		Protocols: p,
	}

	app.logger.Info("starting server", "addr", srv.Addr) //	"env", app.config.env

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
