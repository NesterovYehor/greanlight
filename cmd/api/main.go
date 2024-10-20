package main

import (
	"context"
	"database/sql"
	"expvar"
	"os"
	"runtime"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"greenlight.nesty.net/internal/data"
	"greenlight.nesty.net/internal/jsonlog"
	"greenlight.nesty.net/internal/mailer"
)


type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	cnf := config{}
	cnf.New()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cnf)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cnf,
		logger: logger,
		models: data.NewModel(db),
		mailer: mailer.New(cnf.smtp.port, cnf.smtp.host, cnf.smtp.username, cnf.smtp.password, cnf.smtp.sender), // Corrected line
		wg:     sync.WaitGroup{},
	}

	logger.PrintInfo("database connection pool established", nil)

	err = app.server()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cnf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cnf.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cnf.db.maxOpenConns)
	db.SetMaxIdleConns(cnf.db.maxIdleConns)

	duration, err := time.ParseDuration(cnf.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
