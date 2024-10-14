package main

import (
	"errors"
	"fmt"
	"net/http"

	"greenlight.nesty.net/internal/data"
	"greenlight.nesty.net/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json: "name"`
		Email    string `json: "email"`
		Password string `json: password_hash"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newUser := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = newUser.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidUser(v, newUser); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.User.Insert(newUser)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.backgropund(func() {
		err = app.mailer.Send(newUser.Email, "user_welcome.tmpl", newUser)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, envelope{"user": newUser}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
