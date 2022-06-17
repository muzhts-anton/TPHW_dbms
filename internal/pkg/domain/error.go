package domain

type ErrorResp struct {
	Message string `json:"message"`
}

const (
	ErrorBadRequest          = "Bad request"
	ErrorNotFound            = "Item is not found"
	ErrorConflict            = "Already exist"
	ErrorInternalServerError = "Internal Server Error"
)

//easyjson:skip
type NetError struct {
	Err        error
	Statuscode int
	Message    string
}
