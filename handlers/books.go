package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
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

func GetBookByID(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var resBook []responses.Book = make([]responses.Book, 1)
	var reqBook responses.Book

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	reqBook = request.Book

	query := "SELECT * FROM books WHERE book_id = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	err = db.QueryRowContext(ctx, query, reqBook.ID).Scan(&reqBook.ID, &reqBook.Title, &reqBook.Author)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusNotFound
	}

	resBook[0] = reqBook

	response.Success = true
	response.Data.Books = resBook

	return response, http.StatusOK

}

func GetBooks(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var books []responses.Book
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

	query := "SELECT books.book_id, books.title, books.author_fk FROM books INNER JOIN lib ON books.book_id = lib.book_fk WHERE lib_fk IN (SELECT lib_id FROM libs WHERE user_fk = ?)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*5)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, user.Username)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusBadRequest
	}
	for {
		var book responses.Book
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

	log.Print(books)

	response.Data.Books = books
	return response, http.StatusAccepted
}

func GetBooksBySearch(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user responses.ResponseUser
	var books []responses.Book
	var reqBooks responses.Book

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusBadRequest
	}

	user = request.User
	reqBooks = request.Book
	log.Println(reqBooks.ID)

	query := "SELECT book_id, title, author_fk FROM books WHERE MATCH(title, author_fk) AGAINST(? IN NATURAL LANGUAGE MODE)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*2)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, reqBooks.ID)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusNotFound
	}
	books = compileBooks(rows, books)

	query = "SELECT book_id, title, author FROM new_books WHERE MATCH(title, author) AGAINST(? IN NATURAL LANGUAGE MODE) AND user_fk = ?"
	ctx, cancel = context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	rows, err = db.QueryContext(ctx, query, reqBooks.ID, user.Username)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusNotFound
	}
	books = compileBooks(rows, books)

	response.Success = true
	response.Data.Books = books

	log.Println(books)

	return response, http.StatusOK

}

func GetBooksByLib(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var lib responses.Lib
	var books []responses.Book

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println(err.Error())
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusBadRequest
	}

	lib = request.Lib
	log.Println(lib)

	query := "SELECT books.book_id, books.title, books.author_fk FROM books INNER JOIN lib ON books.book_id = lib.book_fk WHERE lib.lib_fk = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, lib.LibId)
	if err != nil {
		log.Println(err.Error())
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusNotFound
	}

	books = compileBooks(rows, books)

	response.Success = true
	response.Data.Books = books

	return response, http.StatusOK

}

func SearchBooksByLib(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var lib responses.Lib
	var books []responses.Book

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse

		return response, http.StatusBadRequest
	}

	lib = request.Lib

	query := "SELECT books.book_id, books.title, books.author_fk FROM books INNER JOIN lib ON books.book_id = lib.book_fk WHERE lib.lib_fk = ?"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, query, lib.LibId)
	if err != nil {
		log.Println(err.Error())
		errResponse := responses.ResponseError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = errResponse
		return response, http.StatusNotFound
	}

	books = compileBooks(rows, books)

	response.Success = true
	response.Data.Books = books

	return response, http.StatusAccepted

}

func compileBooks(rows *sql.Rows, books []responses.Book) []responses.Book {
	for {
		var book responses.Book
		var id uint64
		rows.Next()
		err := rows.Scan(&id, &book.Title, &book.Author)
		if err != nil {
			// log.Println(err.Error())
			break
		}
		book.ID = fmt.Sprintf("%d", id)
		books = append(books, book)
	}
	return books
}

func CreateFromUser(r *http.Request, db *sql.DB) (responses.Response, int) {

	var request responses.Request
	var response responses.Response
	var user string
	var newBook responses.Book
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusBadRequest
	}
	user = request.User.Username
	newBook = request.Book

	query := "INSERT INTO new_books(contrib_id, book_id, title, author, user_fk) VALUES(UUID_TO_BIN(UUID()), ?, ?, ?, ?)"
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()
	_, err = db.QueryContext(ctx, query, newBook.ID, newBook.Title, newBook.Author, user)
	if err != nil {
		errResponse := responses.JSONError{
			Err: err.Error(),
		}
		response.Success = false
		response.Data.Err = responses.ResponseError(errResponse)
		return response, http.StatusInternalServerError
	}

	response.Success = true
	response.Data.Books = append(response.Data.Books, newBook)

	return response, http.StatusAccepted

}
