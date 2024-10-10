package main

import (
	"fmt"
	"net/http"
)

func (app *application) URLNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("route does not exist"))
}

func (app *application) errorResponse(
	w http.ResponseWriter, r *http.Request, status int, message any,
) {

	err := app.writeJSON(w, status, message, nil)
	if err != nil {
		w.WriteHeader(status)
	}
}

// not permitted
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account doesn't have the necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// not authorized
func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// not activated
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}

// invalid token response
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {

	// 	Note: We’re including a WWW-Authenticate: Bearer header here to help inform or
	// remind the client that we expect them to authenticate using a bearer token.
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

// invalid credentials
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (app *application) failedValidationResponse(
	w http.ResponseWriter, r *http.Request, errors map[string]string,
) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) badRequestResponse(
	w http.ResponseWriter, r *http.Request, err error,
) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) serverErrorResponse(
	w http.ResponseWriter, r *http.Request, err error,
) {
	msg := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(
	w http.ResponseWriter, r *http.Request,
) {

	msg := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) methodNotAllowedResponse(
	w http.ResponseWriter, r *http.Request,
) {

	msg := fmt.Sprintf("the %s mehtod is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, msg)
}
