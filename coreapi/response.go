package coreapi

import (
	"encoding/json"
	"errors"
	"net/http"
)

// A Response contains information to be written
type Response struct {
	Status int
	Body   []byte
	Header http.Header
	Error  error
}

// WriteTo writes the Response to a http.ResponseWriter
func (r *Response) WriteTo(out http.ResponseWriter) {
	header := out.Header()
	for k, v := range r.Header {
		header[k] = v
	}
	out.WriteHeader(r.Status)
	if r.Body != nil && len(r.Body) > 0 {
		if _, err := out.Write(r.Body); err != nil {
			panic(err) // can't write to response
		}
	}
}

// JSONResponse builds a Response with the body formatted as JSON
func JSONResponse(status int, body interface{}) *Response {
	var b []byte
	var err error
	if b, err = json.Marshal(body); err != nil {
		return JSONErrorResponse(500, errors.New("Error marshalling JSON"))
	}
	resp := &Response{
		Body:   b,
		Status: status,
		Header: make(http.Header),
	}
	resp.Header.Set("Content-Type", "application/json")
	return resp
}

// JSONErrorResponse returns a Response object formatted to display given error as JSON
func JSONErrorResponse(status int, err error) *Response {
	renderedError := errorResponse{Type: "error", Status: status, Message: err.Error()}
	res := JSONResponse(status, renderedError)
	res.Error = err
	return res
}

// EmptyResponse returns a Response object without a body
func EmptyResponse(status int) *Response {
	return &Response{Status: status, Header: make(http.Header)}
}

type errorResponse struct {
	Type    string `json:"type"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
