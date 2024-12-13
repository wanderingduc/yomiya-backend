package responses

type JSONResponse struct {
	Success bool
	Data    interface{}
	Meta    interface{}
}

type JSONError struct {
	Err string
}
