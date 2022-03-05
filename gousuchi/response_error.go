package gousuchi

import (
	"fmt"
	"net/http"

	"github.com/indece-official/go-gousu/v2/gousu/logger"
)

type ResponseError struct {
	Request       *http.Request
	StatusCode    int
	PublicMessage string
	DetailedError error
}

var _ IResponse = (*ResponseError)(nil)

func (r *ResponseError) GetRequest() *http.Request {
	return r.Request
}

func (r *ResponseError) Write(w http.ResponseWriter) IResponse {
	w.WriteHeader(r.StatusCode)
	fmt.Fprintf(w, r.PublicMessage)

	return nil
}

func (r *ResponseError) Log(log *logger.Log) {
	if r.StatusCode >= 500 {
		log.Errorf("%s %s - %d %s", r.Request.Method, r.Request.RequestURI, r.StatusCode, r.DetailedError)
	} else {
		log.Warnf("%s %s - %d %s", r.Request.Method, r.Request.RequestURI, r.StatusCode, r.DetailedError)
	}
}

func InternalServerError(request *http.Request, detailedMessage string, args ...interface{}) *ResponseError {
	return &ResponseError{
		Request:       request,
		StatusCode:    http.StatusInternalServerError,
		PublicMessage: "Internal server error",
		DetailedError: fmt.Errorf(detailedMessage, args...),
	}
}

func NotFound(request *http.Request, detailedMessage string, args ...interface{}) *ResponseError {
	return &ResponseError{
		Request:       request,
		StatusCode:    http.StatusNotFound,
		PublicMessage: "Not found",
		DetailedError: fmt.Errorf(detailedMessage, args...),
	}
}

func BadRequest(request *http.Request, detailedMessage string, args ...interface{}) *ResponseError {
	return &ResponseError{
		Request:       request,
		StatusCode:    http.StatusBadRequest,
		PublicMessage: "Bad request",
		DetailedError: fmt.Errorf(detailedMessage, args...),
	}
}

func Unauthorized(request *http.Request, detailedMessage string, args ...interface{}) *ResponseError {
	return &ResponseError{
		Request:       request,
		StatusCode:    http.StatusUnauthorized,
		PublicMessage: "Unauthorized",
		DetailedError: fmt.Errorf(detailedMessage, args...),
	}
}

func Forbidden(request *http.Request, detailedMessage string, args ...interface{}) *ResponseError {
	return &ResponseError{
		Request:       request,
		StatusCode:    http.StatusForbidden,
		PublicMessage: "Forbidden",
		DetailedError: fmt.Errorf(detailedMessage, args...),
	}
}
