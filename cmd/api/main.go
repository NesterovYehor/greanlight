package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
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

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(conf)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database connection pool established")

	app := &application{
		config: conf,
		logger: logger,
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Printf("starting %s server on %s", conf.env, srv.Addr)

	err = srv.ListenAndServe()
	logger.Fatal(err)
}

func openDB(conf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.db.dsn)
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
