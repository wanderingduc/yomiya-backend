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

type Lib struct {
	Libname string
}

func GetLibs(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user string
	var libs []responses.Lib

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	user = request.User.Username

	query := "SELECT * FROM libs WHERE user_fk = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, user)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusInternalServerError
	}

	for {
		var lib responses.Lib
		var a string
		rows.Next()
		err := rows.Scan(&lib.LibId, &lib.LibName, &a)
		if err != nil {
			log.Println(err.Error())
			break
		}
		libs = append(libs, lib)

	}

	response.Success = true
	response.Data.Libs = libs

	return response, http.StatusAccepted
}
