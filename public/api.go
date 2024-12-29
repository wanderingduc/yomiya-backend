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
	// auth = middleware.CheckRole(auth, app.db) // NEEDS TESTING

	mux := http.NewServeMux()
	mux.HandleFunc("POST /dev/v1/users/create", app.createUser)
	mux.HandleFunc("POST /dev/v1/users/auth", app.authUser)
	mux.HandleFunc("POST /dev/v1/users/authtoken", app.authToken)
	mux.HandleFunc("/dev/v1/test", app.testGet)
	mux.HandleFunc("POST /dev/v1/users/userid", app.getUser)
	mux.Handle("/dev/v1/auth/", auth)

	// BOOK ROUTES SHOULD BE ADDED TO AUTH MUX AFTER COMPLETION

	mux.HandleFunc("POST /dev/v1/books/create", app.createBook)
	mux.HandleFunc("POST /dev/v1/books/bookid", app.getBookByID)
	mux.HandleFunc("POST /dev/v1/books/search", app.getBooksBySearch)
	mux.HandleFunc("POST /dev/v1/books/get", app.getBooks)
	mux.HandleFunc("POST /dev/v1/books/lib", app.searchBooksByLib)

	// LIB ROUTES SHOULD BE ADDED TO AUTH MUX AFTER COMPLETION

	mux.HandleFunc("POST /dev/v1/libs/get", app.getLibs)
	mux.HandleFunc("POST /dev/v1/libs/libid", app.getBooksByLib)
	mux.HandleFunc("POST /dev/v1/libs/addbook", app.addBookToLib)
	mux.HandleFunc("DELETE /dev/v1/libs/libid", app.deleteLib)

	// ADMIN ROUTES SHOULD BE ADDED TO AUTH MUX AFTER COMPLETION

	mux.HandleFunc("DELETE /dev/v1/admin/user", app.deleteUser) // NEEDS TESTING

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
