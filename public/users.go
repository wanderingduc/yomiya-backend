package main

import (
	"encoding/json"
	"log"
	"net/http"
	"yomiya/backend/api/handlers"
	"yomiya/backend/api/responses"
)

func (app *application) getUser(w http.ResponseWriter, r *http.Request) {

	// response, status := handlers.GetUserByID(r, app.db)

	// w.WriteHeader(status)
	// w.Header().Set("Content-Type", "application/json")
	// err := json.NewEncoder(w).Encode(response)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("hello")

}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.CreateUser(r, app.db)

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) authUser(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.AuthUser(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (app *application) authToken(w http.ResponseWriter, r *http.Request) {

	response, status := handlers.AuthToken(r, app.db)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *application) testGet(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	var res responses.JSONResponse
	res.Success = true
	res.Data = "asdfasdf"
	json.NewEncoder(w).Encode(res)

}
