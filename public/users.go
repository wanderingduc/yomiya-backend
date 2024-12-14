package main

import (
	"encoding/json"
	"log"
	"net/http"
	"yomiya/backend/api/handlers"
)

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.GetUserByID(r, app.db)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)

}

func (app *application) createuser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.CreateUser(r, app.db)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
}

func (app *application) authUser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.AuthUser(r, app.db)

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)

}
