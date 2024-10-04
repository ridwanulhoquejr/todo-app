package main

import "net/http"

func (app *application) singleTodoHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": err.Error()}, nil)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"id": id}, nil)
	// fmt.Fprintf(w, "succesfully ping\n")
}
