package main

import "net/http"

func (app *application) URLNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("route does not exist"))
}

// func (app *application) errorResponse(
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	status int,
// 	message any,
// ) {
// 	env := envelope{"error": message}

// 	err := app.writeJSON(w, status, env, nil)
// 	if err != nil {
// 		app.logError(r, err)
// 		w.WriteHeader(500)
// 	}
// }
