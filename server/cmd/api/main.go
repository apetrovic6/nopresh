package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

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

type app struct {
	mux         *http.ServeMux
	config      *config
	logger      *slog.Logger
	models      data.Models
	middlewares Middlewares
	jwt         *auth.JWT
}

func main() {
	var cfg config

	cfg.loadFlags()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := gorm.Open(postgres.Open(cfg.db.getDsn(cfg.tz)), &gorm.Config{})

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
		mux:         http.NewServeMux(),
		logger:      logger,
		config:      &cfg,
		middlewares: NewMiddleware(jwt, logger, cfg.domains),
		models:      data.NewModels(db),
		jwt:         jwt,
	}

	app.registerReflection()
	app.registerRoutes()
	app.registerWeb()

	p := new(http.Protocols)

	p.SetHTTP1(true)

	// Use h2c so we can serve HTTP/2 without TLS.
	p.SetUnencryptedHTTP2(true)

	srv := http.Server{
		Addr:      cfg.getHost(),
		Handler:   app.mux,
		Protocols: p,
	}

	app.logger.Info("starting server", "addr", srv.Addr) //	"env", app.config.env

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
