package main

import (
	"database/sql"
	"net/http"
	"yomiya/backend/api/middleware"

	"github.com/rs/cors"
)

type application struct {
	conf config
	db   *sql.DB
}

type config struct {
	addr   string
	dbConf dbConfig
}

type dbConfig struct {
	addr        string
	MaxOpenConn int
	MaxIdleConn int
	MaxIdleTime string
}

func (app *application) mount() http.Handler {

	auth := app.mountAuth()
	auth = middleware.CheckToken(auth)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /dev/v1/create", app.createUser)
	mux.HandleFunc("POST /dev/v1/auth", app.authUser)
	mux.Handle("/dev/v1/auth/", auth)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Authorized", "Origin"},
	})

	router := c.Handler(mux)

	return router
}

func (app *application) mountAuth() http.Handler {
	authHandler := http.NewServeMux()
	authHandler.HandleFunc("POST /test", app.getUser)
	return http.StripPrefix("/dev/v1/auth", authHandler)
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:    app.conf.addr,
		Handler: mux,
	}

	return http.ListenAndServe(server.Addr, server.Handler)

}
