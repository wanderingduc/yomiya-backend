package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
	"yomiya/backend/api/responses"
)

type Lib struct {
	Libname string
}

func GetLibs(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var response responses.JSONResponse
	var user string

	err := json.NewDecoder(r.Body).Decode(&user)
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

	query := "SELECT* FROM libs WHERE user_fk = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, user)
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

	for {
		rows.Next()
		if rows == nil {
			break
		}

	}

	return response, http.StatusAccepted
}
