package handlers // NEEDS TESTING

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"yomiya/backend/api/responses"
)

func CreateAdmin(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user responses.ResponseUser

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = true
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	user = request.User

	query := "INSERT INTO admins(username, password, created_at) VALUES(?, ?, CURRENT_TIMESTAMP())"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	_, err = db.QueryContext(ctx, query, user.Username, user.Password)
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

func DeleteUser(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user responses.ResponseUser

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = true
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	user = request.User

	query := "DELETE FROM users WHERE username = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	_, err = db.QueryContext(ctx, query, user.Username)
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
