package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yomiya/backend/api/auth"
	"yomiya/backend/api/responses"

	"golang.org/x/crypto/bcrypt"
)

type ClientUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ServerUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

func CreateUser(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var newUser responses.ResponseUser

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	newUser = request.User

	query := "INSERT INTO users(username, passwd, created_at, updated_at) VALUES(?, ?, CURRENT_TIMESTAMP(), CURRENT_TIMESTAMP())"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	hashPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: "Could not create user",
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusInternalServerError
	}
	_, err = db.QueryContext(ctx, query, newUser.Username, hashPass)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: "Could not create user",
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusInternalServerError
	}

	token, err := auth.CreateJWT(newUser.Username)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusInternalServerError
	}

	newUser.Jwt = token

	response.Success = true
	response.Data.User = []responses.ResponseUser{newUser}

	return response, http.StatusCreated

}

func GetUserByID(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user responses.ResponseUser
	// var resUser ServerUser
	// var findUser string

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	user = request.User

	query := "SELECT * FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, user.Username).Scan(&user.Username, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(err.Error())
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	response.Success = true
	response.Data.User = []responses.ResponseUser{user}

	return response, http.StatusAccepted

}

func AuthUser(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var toAuth responses.ResponseUser
	var pass string

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}
	toAuth = request.User

	query := "SELECT passwd FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, toAuth.Username).Scan(&pass)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	err = auth.CheckPassword([]byte(pass), []byte(toAuth.Password))
	if err != nil {
		errResponse := responses.ResponseError{
			Err: "Invalid username or password",
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	token, err := auth.CreateJWT(toAuth.Username)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusInternalServerError
	}

	toAuth.Jwt = token

	response.Success = true
	response.Data.User = []responses.ResponseUser{toAuth}

	return response, http.StatusAccepted

}

func AuthToken(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var token string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}
	token = request.User.Jwt

	err = auth.CheckToken(token)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	response.Success = true

	return response, http.StatusAccepted
}

func ReportBug(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user responses.ResponseUser

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	user = request.User

	query := "INSERT INTO reports(user_fk, bug) VALUES(?, ?)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	_, err = db.QueryContext(ctx, query, user.Username, user.Jwt)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	response.Success = true

	return response, http.StatusAccepted

}
