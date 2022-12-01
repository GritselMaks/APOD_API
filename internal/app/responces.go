package app

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Data interface{} `json:"data"`
}

func NewResponse(data interface{}) *Response {
	return &Response{data}
}

type ErrorResponse struct {
	ErrorMsg string `json:"error"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{ErrorMsg: msg}
}

// Forms a response with http Status and data
func ResponseWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	payload := NewResponse(data)
	response, err := json.Marshal(payload)
	if err != nil {
		//FIXME
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func RespondWithError(w http.ResponseWriter, statusCode int, msg string) {
	payload := NewErrorResponse(msg)
	response, err := json.Marshal(payload)
	if err != nil {
		//FIXME
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func RespondWithPicture(w http.ResponseWriter, statusCode int, data []byte) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}
