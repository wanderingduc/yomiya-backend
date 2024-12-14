package main

import (
	"database/sql"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dev/v1/test", app.getUser)

	return mux
}

func (app *application) run(mux http.Handler) error {

	server := &http.Server{
		Addr:    app.conf.addr,
		Handler: mux,
	}

	return http.ListenAndServe(server.Addr, server.Handler)

}
