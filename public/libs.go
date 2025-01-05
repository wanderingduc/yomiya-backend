package main

import (
	"encoding/json"
	"log"
	"net/http"
	"yomiya/backend/api/handlers"
)

func (app *application) createLib(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.CreateLib(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) getLibs(w http.ResponseWriter, r *http.Request) {

	log.Println("Fetching libs")

	response, status := handlers.GetLibs(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) getLibsBySearch(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.GetLibsBySearch(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) addBookToLib(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.AddBookToLib(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) deleteLib(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.DeleteLib(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}
