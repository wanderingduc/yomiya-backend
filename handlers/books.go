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

type Book struct {
	ID     string `json:"book_id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type UserBook struct {
	ID     string `json:"book_id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	User   string `json:"username"`
}

func GetBookByID(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var response responses.JSONResponse
	var resBook Book
	var reqBook Book

	err := json.NewDecoder(r.Body).Decode(&reqBook)
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

	query := "SELECT * FROM new_books WHERE book_id = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, reqBook.ID).Scan(&resBook.ID, &resBook.Title, &resBook.Author)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response = responses.JSONResponse{
			Success: false,
			Data:    errResponse,
			Meta:    nil,
		}
		return response, http.StatusNotFound
	}

	response = responses.JSONResponse{
		Success: true,
		Data:    resBook,
		Meta:    nil,
	}

	return response, http.StatusOK

}

func GetBooks(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var response responses.JSONResponse
	var books []Book
	var user ClientUser
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

	query := "SELECT books.book_id, books.title, books.author_fk FROM books INNER JOIN lib ON books.book_id = lib.book_fk WHERE lib_fk IN (SELECT lib_id FROM libs WHERE user_fk = ?)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, user.Username)
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
	for {
		var book Book
		rows.Next()
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			log.Println(err.Error())
			break
			// errResponse := responses.JSONError{
			// 	Err: err.Error(),
			// }
			// response = responses.JSONResponse{
			// 	Success: false,
			// 	Data:    errResponse,
			// 	Meta:    nil,
			// }
			// return response, http.StatusInternalServerError
		}
		books = append(books, book)
	}

	response = responses.JSONResponse{
		Success: true,
		Data:    books,
		Meta:    nil,
	}

	return response, http.StatusAccepted
}

func CreateFromUser(r *http.Request, db *sql.DB) (responses.JSONResponse, int) {

	var response responses.JSONResponse
	var newBook UserBook
	err := json.NewDecoder(r.Body).Decode(&newBook)
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

	query := "INSERT INTO new_books(contrib_id, book_id, title, author, user_fk) VALUES(UUID_TO_BIN(UUID()), ?, ?, ?, ?)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	_, err = db.QueryContext(ctx, query, newBook.ID, newBook.Title, newBook.Author, newBook.User)
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
		Data:    newBook,
		Meta:    nil,
	}

	return response, http.StatusAccepted

}
