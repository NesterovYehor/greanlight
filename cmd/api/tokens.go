package main

import (
	"errors"
	"net/http"
	"time"

	"greenlight.nesty.net/internal/data"
	"greenlight.nesty.net/internal/validator"
)

func (app *application) createAuthentidcationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	token, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, envelope{"authentication_token": token}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json: "email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no matching email address found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
	}

	if !user.Activated {
		v.AddError("email", "user account must be activated")
		app.inactiveAccountResponse(w, r)
		return
	}

	token, err := app.models.Token.New(user.ID, 30*time.Minute, data.ScopePasswordReset)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.backgropund(func() {
		data := map[string]any{
			"passwordResetToken": token.Plaintext,
		}

		err := app.mailer.Send(user.Email, "token_password_reset.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	env := envelope{"message": "an email will be sent to you containing password reset instructions"}

	err = app.writeJSON(w, env, http.StatusAccepted, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
