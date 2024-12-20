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

func CreateUser(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var response responses.JSONResponse
	var newUser ClientUser
	var check string

	json.NewDecoder(r.Body).Decode(&newUser)

	query := "SELECT username FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	exist := db.QueryRowContext(ctx, query, newUser.Username).Scan(&check)
	if exist == nil {
		errResponse := responses.JSONError{
			Err: "User already exists",
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusForbidden
	}

	query = "INSERT INTO users(username, passwd, created_at) VALUES(?, ?, CURRENT_TIMESTAMP())"
	ctx, cancel = context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	hashPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := responses.JSONError{
			Err: "Could not create user",
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusInternalServerError
	}
	_, err = db.QueryContext(ctx, query, newUser.Username, hashPass)
	if err != nil {
		errResponse := responses.JSONError{
			Err: "Could not create user",
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusInternalServerError
	}

	response = responses.JSONResponse{
		Success: true,
		Data:    newUser,
		Meta:    nil,
	}

	return response, http.StatusAccepted

}

func GetUserByID(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var resUser ServerUser
	var findUser string
	json.NewDecoder(r.Body).Decode(&findUser)
	query := "SELECT * FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err := db.QueryRowContext(ctx, query, findUser).Scan(&resUser.Username, &resUser.Password, &resUser.CreatedAt)
	if err != nil {
		log.Println(err.Error())
		eRes := responses.JSONError{
			Err: err.Error(),
		}
		response := responses.JSONResponse{
			Success: false,
			Data:    eRes,
			Meta:    nil,
		}
		return response, http.StatusBadRequest
	}

	response := responses.JSONResponse{
		Success: true,
		Data:    resUser,
		Meta:    nil,
	}

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

	response.Success = true
	response.Data.User.Username = toAuth.Username
	response.Data.User.Jwt = token

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
