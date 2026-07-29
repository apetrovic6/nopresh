package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
)

type config struct {
	host        string
	port        int
	env         string
	db          db
	tz          string
	jwtSecret   string
	jwtDuration time.Duration
}

func (c *config) getHost() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

type db struct {
	name         string
	host         string
	port         int
	user         string
	password     string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

func (db *db) getDsn(timezone string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", db.host, db.user, db.password, db.name, db.port, timezone)
}

func (cfg *config) loadFlags() {
	port, err := strconv.Atoi(os.Getenv("PORT"))

	if err != nil {
		msg := fmt.Sprintf("PORT not valid %v", err)
		panic(msg)
	}

	durationStr := os.Getenv("JWT_DURATION")

	if len(durationStr) == 0 {
		durationStr = "15m"
	}

	jwtDuration, err := time.ParseDuration(durationStr)

	if err != nil {
		panic(err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		msg := fmt.Sprintf("DB_PORT not valid %v", err)
		panic(msg)
	}

	flag.StringVar(&cfg.host, "host", os.Getenv("HOST"), "127.0.0.1")
	flag.StringVar(&cfg.tz, "timezone", os.Getenv("APP_TZ"), "Timezone for the app and db connection string")
	// flag.StringVar(&cfg.domains, "domain", os.Getenv("DOMAINS"), "Domains from which the frontend(s) will be hosted on. Domains should be separated by a comma (https://domain1.com,https://domain2.com)")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENVIRONMENT"), "Evironment (development|staging|production)")
	flag.IntVar(&cfg.port, "port", port, "5000")

	flag.StringVar(&cfg.db.name, "db name", os.Getenv("DB_NAME"), "Postgres Name")
	flag.StringVar(&cfg.db.host, "db host", os.Getenv("DB_HOST"), "Postgres Host")
	flag.IntVar(&cfg.db.port, "db port", dbPort, "Postgres Port")
	flag.StringVar(&cfg.db.user, "db user", os.Getenv("DB_USER"), "Postgres User")
	flag.StringVar(&cfg.db.password, "db password", os.Getenv("DB_PASSWORD"), "Postgres Password")

	flag.StringVar(&cfg.jwtSecret, "JWT Secret Key", os.Getenv("JWT_SECRET_KEY"), "JWT Secret Key")
	flag.DurationVar(&cfg.jwtDuration, "JWT Token Duration", jwtDuration, "JWT Token Duration")

	flag.Parse()
}
