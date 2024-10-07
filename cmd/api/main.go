package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var conf config

	flag.IntVar(&conf.port, "port", 4000, "API server port")
	flag.StringVar(&conf.env, "env", "development", "Eniroment (developmet|staging|production)")
    flag.Parse()


     logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)


    app := &application{
        config: conf,
        logger: logger,
    }


    mux := http.NewServeMux()
    mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)


    srv := &http.Server{
        Addr: fmt.Sprintf(":%d", conf.port),
        Handler: mux,
        IdleTimeout: time.Minute,
        ReadTimeout: 10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }


    logger.Printf("starting %s server on %s", conf.env, srv.Addr)

    err := srv.ListenAndServe()
    logger.Fatal(err)
}

