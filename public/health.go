package main

import "net/http"

func (app *application) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Health OK"))
}
