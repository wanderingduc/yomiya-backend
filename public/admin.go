package main // NEEDS TESTING

import (
	"encoding/json"
	"net/http"
	"yomiya/backend/api/handlers"
)

func (app *application) createAdmin(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.CreateAdmin(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) deleteUser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.DeleteUser(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
