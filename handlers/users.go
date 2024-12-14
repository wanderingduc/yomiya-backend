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
	_, invalid, err := db.QueryContext(ctx, query, newUser.Username, newUser.Password)
	if err != nil || invalid != nil {
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

func AuthUser(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var toAuth ClientUser
	var pass string
	var response responses.JSONResponse

	err := json.NewDecoder(r.Body).Decode(&toAuth)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusBadRequest
	}

	query := "SELECT passwd FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, toAuth.Username).Scan(&pass)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusBadRequest
	}

	err = auth.CheckPassword([]byte(pass), []byte(toAuth.Password))
	if err != nil {
		errResponse := responses.JSONError{
			Err: "Invalid username or password",
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusBadRequest
	}

	token, err := auth.CreateJWT(toAuth.Username)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
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
		Data:    token,
		Meta:    nil,
	}

	return response, http.StatusAccepted

}
