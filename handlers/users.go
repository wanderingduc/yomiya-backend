package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"yomiya/backend/api/responses"
)

type ClientUser struct {
	username string
	password string
}

type ServerUser struct {
	username  string
	password  string
	createdAt string
}

func GetUserByID(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var resUser ServerUser
	var findUser string
	json.NewDecoder(r.Body).Decode(&findUser)
	query := "SELECT * FROM users WHERE uname = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err := db.QueryRowContext(ctx, query, findUser).Scan(&resUser.username, &resUser.password, &resUser.createdAt)
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
