package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"greenlight.nesty.net/internal/data"
	"greenlight.nesty.net/internal/jsonlog"
)

const version = "1.0.0"

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
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var conf config

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")

	flag.IntVar(&conf.port, "port", 4000, "API server port")
	flag.StringVar(&conf.env, "env", "development", "Eniroment (developmet|staging|production)")
	flag.StringVar(&conf.db.dsn, "db-dsn", dbURL, "PostgreSQL DSN")

	flag.IntVar(&conf.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL DSN")
	flag.IntVar(&conf.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL DSN")
	flag.StringVar(&conf.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL DSN")

	flag.Float64Var(&conf.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&conf.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&conf.limiter.enable, "limiter-enable", true, "Enable rate limiter")

	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(conf)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: conf,
		logger: logger,
		models: data.NewModel(db),
	}

	err = app.server()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(conf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.db.dsn)

	db.SetMaxOpenConns(conf.db.maxOpenConns)
	db.SetMaxIdleConns(conf.db.maxIdleConns)

	duration, err := time.ParseDuration(conf.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
