package responses

import (
	"time"
)

type JSONResponse struct {
	Success bool
	Data    interface{}
	Meta    interface{}
}

type JSONError struct {
	Err string
}

type ResponseUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Jwt       string `json:"token"`
}

type Book struct {
	ID     string `json:"book_id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

type ResponseBook struct { // MAYBE DELETE
	Books []Book `json:"books"`
}

type Lib struct {
	LibId   string `json:"lib_id"`
	LibName string `json:"lib_name"`
}

type ResponseLib struct { // MAYBE DELETE
	Libs []Lib `json:"libs"`
}

type ResponseError struct {
	Err string `json:"error"`
}

type ResponseData struct {
	User  []ResponseUser `json:"user"`
	Books []Book         `json:"books"`
	Libs  []Lib          `json:"libs"`
	Err   ResponseError  `json:"error"`
}

type ResponseMeta struct {
	Timestamp time.Time `json:"timestamp"`
}

type Response struct {
	Success bool         `json:"Success"`
	Data    ResponseData `json:"Data"`
	Meta    ResponseMeta `json:"Meta"`
}

type Request struct {
	User ResponseUser `json:"user"`
	Book Book         `json:"book"`
	Lib  Lib          `json:"lib"`
}
