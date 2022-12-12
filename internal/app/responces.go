package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JsonErrTmpl is the template to use when returning a JSON error. It is
// rendered using Printf, not json.Encode, so values must be escaped by the
// caller.
const jsonErrTmpl = `{"error":"%s"}`

// Standart form for all responses
type Response struct {
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}) *Response {
	return &Response{data}
}

// Standart form for all Error responses
type ErrorResponse struct {
	ErrorMsg string `json:"error"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{ErrorMsg: msg}
}

// Response with http Status and data
func ResponseWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	payload := NewResponse(data)
	response, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, jsonErrTmpl, http.StatusText(http.StatusInternalServerError))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(response)
}

//Response with http Status and Error
func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	payload := NewErrorResponse(msg)
	response, err := json.Marshal(payload)
	if err != nil {
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, jsonErrTmpl, http.StatusText(http.StatusInternalServerError))
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(response)
}

//Response with http Status and media content
func RespondWithPicture(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}
