package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ridwanulhoquejr/todo-app/internal/validator"
)

// our response wrapper!
type envelope map[string]any

// our response wrapper!
// i use struct bcz map type json.Marshal ordered alphatically, where struct won't
type APIResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code,omitempty"`
	Details any    `json:"details,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// readJSON: convert json data to go-struct data
func (app *application) readJSON(
	w http.ResponseWriter, r *http.Request, dst any) error {
	// Use http.MaxBytesReader() to limit the size of the request body to 1MB.
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {

		// If there is an error during decoding, start the triage...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {

		// Use the errors.As() function to check whether the error has the type
		// *json.SyntaxError. If it does, then return a plain-english error message
		// which includes the location of the problem.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		// In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
		// for syntax errors in the JSON. So we check for this using errors.Is() and
		// return a generic error message. There is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		// Likewise, catch any *json.UnmarshalTypeError errors. These occur when the
		// JSON value is the wrong type for the target destination. If the error relates
		// to a specific field, then we include that in our error message to make it
		// easier for the client to debug.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		// An io.EOF error will be returned by Decode() if the request body is empty. We
		// check for this with errors.Is() and return a plain-english error message
		// instead.
		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		// A json.InvalidUnmarshalError error will be returned if we pass a non-nil
		// pointer to Decode(). We catch this and panic, rather than returning an error
		// to our handler. At the end of this chapter we'll talk about panicking
		// versus returning errors, and discuss why it's an appropriate thing to do in
		// this specific situation.
		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		// For anything else, return the error message as-is.
		default:
			return err
		}
	}

	// Call Decode() again, using a pointer to an empty anonymous struct as the
	// destination. If the request body only contained a single JSON value this will
	// return an io.EOF error. So if we get anything else, we know that there is
	// additional data in the request body and we return our own custom error message.
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func createAPIResponse(statusCode int, data any) APIResponse {
	// Create a default response structure
	rsp := APIResponse{
		Code: statusCode,
	}

	// Set status and response fields based on the status code
	switch {
	case statusCode >= 400 && statusCode < 500:
		rsp.Status = "fail"
		rsp.Details = data
	case statusCode >= 500:
		rsp.Status = "error"
		rsp.Details = data
	default:
		rsp.Status = "success"
		rsp.Data = data
	}

	return rsp
}

// writeJSON: convert go-struct to json data
func (app *application) writeJSON(
	w http.ResponseWriter, status int, data any, headers http.Header) error {

	response := createAPIResponse(status, data)

	js, err := json.Marshal(response)
	if err != nil {
		return err
	}

	// At this point, we know that we won't encounter any more errors before writing the
	// response, so it's safe to add any headers that we want to include. We loop
	// through the header map and add each header to the http.ResponseWriter header map.
	// Note that it's OK if the provided header map is nil. Go doesn't throw an error
	// if you try to range over (or generally, read from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "applicaiton/json")
	// for sending the status code via header
	w.WriteHeader(status)
	w.Write([]byte(js))

	return nil
}

// The readInt() helper reads a string value from the query string and converts it to an
// integer before returning. If no matching key could be found it returns the provided
// default value.
func (app *application) readInt(
	qs url.Values, key string, defaultValue int, v *validator.Validator) int {

	// extract the query string from url.Values
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}

	// convert query string value to int!
	i, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be a number")
		return defaultValue
	}
	return i
}

// read time.Time
func (app *application) readTime(qs url.Values, key string, defaultValue time.Time) time.Time {

	t := qs.Get(key)

	// empty or null ? return default
	if t == "" || t == "null" {
		return defaultValue
	}

	// Parse the time string into a time.Time object
	parsedTime, err := time.Parse("2006-01-02", t)
	if err != nil {
		app.logger.PrintInfo("error while parsing datetime sort", nil)
		return defaultValue
	}

	return parsedTime
}

// read boolean
// TODO: need to complete this for completed field
// func (app *application) readBoolean(
// 	qs url.Values, key, defaultValue string, v *validator.Validator) string {

// 	return ""
// }

// read string
func (app *application) readString(
	qs url.Values, key, defaultValue string) string {

	s := qs.Get(key)
	if s == "" || s == "null" {
		return defaultValue
	}
	return s
}

// read CSV
// func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {

// 	csv := qs.Get(key)
// 	if csv == "" {
// 		return defaultValue
// 	}
// 	return strings.Split(csv, ",")
// }

func (app *application) readIDParam(r *http.Request) (int64, error) {

	param := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		return 0, errors.New("invalid path parameter given for id")
	}

	return id, nil
}
