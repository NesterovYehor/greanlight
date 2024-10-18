package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps    float64
		burst  int
		enable bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

func (cnf *config) New() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	flag.IntVar(&cnf.port, "port", 4000, "API server port")
	flag.StringVar(&cnf.env, "env", "development", "Eniroment (developmet|staging|production)")
	flag.StringVar(&cnf.db.dsn, "db-dsn", dbURL, "PostgreSQL DSN")

	flag.IntVar(&cnf.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL DSN")
	flag.IntVar(&cnf.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL DSN")
	flag.StringVar(&cnf.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL DSN")

	flag.Float64Var(&cnf.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cnf.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cnf.limiter.enable, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cnf.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cnf.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cnf.smtp.username, "smtp-username", "e9e77e413749f7", "SMTP username")
	flag.StringVar(&cnf.smtp.password, "smtp-password", "636272c0c7428b", "SMTP password")
	flag.StringVar(&cnf.smtp.sender, "smtp-sender", "Test User <test@mailtrap.io>", "SMTP sender")

	flag.Parse()
}
