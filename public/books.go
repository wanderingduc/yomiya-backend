package main

import (
	"encoding/json"
	"log"
	"net/http"
	"yomiya/backend/api/handlers"
)

func (app *application) createBook(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.CreateFromUser(r, app.db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) getBookByID(w http.ResponseWriter, r *http.Request) {

	log.Println("Getting book")

	response, status := handlers.GetBookByID(r, app.db)

	log.Println(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) getBooks(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.GetBooks(r, app.db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (app *application) getBooksBySearch(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.GetBooksBySearch(r, app.db)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}
