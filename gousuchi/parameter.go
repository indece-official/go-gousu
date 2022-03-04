package gousuchi

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"gopkg.in/guregu/null.v4"
)

// QueryParamString loads a parameter from the url's query
//
// If the parameter is empty a BadRequest-Response is returned, else
// the response is nil
func QueryParamString(request *http.Request, name string) (string, IResponse) {
	value := request.URL.Query().Get(name)
	if value == "" {
		return "", BadRequest(request, "Empty query param %s", name)
	}

	return value, nil
}

// OptionalQueryParamString loads a parameter from the url's query
//
// Response is always nil and only returned for compatiblity here
func OptionalQueryParamString(request *http.Request, name string) (null.String, IResponse) {
	value := request.URL.Query().Get(name)
	if value == "" {
		return null.String{}, nil
	}

	return null.StringFrom(value), nil
}

// QueryParamInt64 loads a parameter from the url's query and parses it as int64
//
// If the parameter is not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func QueryParamInt64(request *http.Request, name string) (int64, IResponse) {
	valueStr := request.URL.Query().Get(name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalQueryParamInt64 loads a parameter from the url's query and parses it as int64
//
// If the parameter is not empty and not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func OptionalQueryParamInt64(request *http.Request, name string) (null.Int, IResponse) {
	valueStr := request.URL.Query().Get(name)
	if valueStr == "" {
		return null.Int{}, nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return null.Int{}, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.IntFrom(value), nil
}

// QueryParamBool loads a parameter from the url's query and parses it as bool
//
// If the parameter is not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func QueryParamBool(request *http.Request, name string) (bool, IResponse) {
	valueStr := request.URL.Query().Get(name)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalQueryParamBool loads a parameter from the url's query and parses it as bool
//
// If the parameter is not empty and not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func OptionalQueryParamBool(request *http.Request, name string) (null.Bool, IResponse) {
	valueStr := request.URL.Query().Get(name)
	if valueStr == "" {
		return null.Bool{}, nil
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return null.Bool{}, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.BoolFrom(value), nil
}

// URLParamString loads a parameter from the url
//
// If the parameter empty a BadRequest-Response is returned, else
// the response is nil
func URLParamString(request *http.Request, name string) (string, IResponse) {
	value := chi.URLParam(request, name)
	if value == "" {
		return "", BadRequest(request, "Empty url param %s", name)
	}

	return value, nil
}

// OptionalURLParamString loads a parameter from the url
//
// Response is always nil and only returned for compatiblity here
func OptionalURLParamString(request *http.Request, name string) (null.String, IResponse) {
	value := chi.URLParam(request, name)
	if value == "" {
		return null.String{}, nil
	}

	return null.StringFrom(value), nil
}

// URLParamInt64 loads a parameter from the url and parses it as int64
//
// If the parameter is not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func URLParamInt64(request *http.Request, name string) (int64, IResponse) {
	valueStr := chi.URLParam(request, name)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return 0, BadRequest(request, "Invalid url param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalURLParamInt64 loads a parameter from the url and parses it as int64
//
// If the parameter is not empty and not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func OptionalURLParamInt64(request *http.Request, name string) (null.Int, IResponse) {
	valueStr := chi.URLParam(request, name)
	if valueStr == "" {
		return null.Int{}, nil
	}

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return null.Int{}, BadRequest(request, "Invalid query param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.IntFrom(value), nil
}

// URLParamBool loads a parameter from the url and parses it as bool
//
// If the parameter is not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func URLParamBool(request *http.Request, name string) (bool, IResponse) {
	valueStr := chi.URLParam(request, name)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, BadRequest(request, "Invalid url param %s (value: '%s'): %s", name, valueStr, err)
	}

	return value, nil
}

// OptionalURLParamBool loads a parameter from the url and parses it as bool
//
// If the parameter is not empty and not a valid bool a BadRequest-Response is returned, else
// the response is nil. Accepted values are 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False
func OptionalURLParamBool(request *http.Request, name string) (null.Bool, IResponse) {
	valueStr := chi.URLParam(request, name)
	if valueStr == "" {
		return null.Bool{}, nil
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return null.Bool{}, BadRequest(request, "Invalid url param %s (value: '%s'): %s", name, valueStr, err)
	}

	return null.BoolFrom(value), nil
}

// URLParamInt64Slice loads a parameter from the url and parses it as a comma-separated list of int64
//
// If the parameter is not a valid int64 a BadRequest-Response is returned, else
// the response is nil
func URLParamInt64Slice(request *http.Request, name string) ([]int64, IResponse) {
	values := []int64{}

	valueStr := chi.URLParam(request, name)
	valueStrParts := strings.Split(valueStr, ",")

	for i, valueStrPart := range valueStrParts {
		value, err := strconv.ParseInt(valueStrPart, 10, 64)
		if err != nil {
			return []int64{}, BadRequest(request, "Invalid url param %s (%d. value: '%s'): %s", name, i+1, valueStrPart, err)
		}
		values = append(values, value)
	}

	return values, nil
}
