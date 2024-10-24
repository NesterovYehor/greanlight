package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdwonError := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.logger.PrintInfo("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdwonError <- err
		}
		app.logger.PrintInfo("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.wg.Wait()
		shutdwonError <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": strconv.Itoa(app.config.port),
		"env":  app.config.env,
	})

	err := srv.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdwonError
	if err != nil {
		return err
	}
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
