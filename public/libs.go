package main

import (
	"encoding/json"
	"net/http"
	"yomiya/backend/api/handlers"
)

func (app *application) getLibs(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.GetLibs(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
