package net

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type Error struct {
	Error string `json:"error,omitempty"`
}

func Success[K any](w http.ResponseWriter, result K) {
	addCorsResponseHeaders(w)
	writeResponse(w, http.StatusOK, result)
}

func HaltBadRequest(w http.ResponseWriter, error string) {
	addCorsResponseHeaders(w)
	writeResponse(w, http.StatusBadRequest, Error{Error: error})
}

func HaltUnauthorized(w http.ResponseWriter, error string) {
	addCorsResponseHeaders(w)
	writeResponse(w, http.StatusUnauthorized, Error{Error: error})
}

func ParseBody[K any](r *http.Request) (*K, error) {
	var result K
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&result)
	if err != nil {
		if err.Error() == "unexpected EOF" {
			return nil, errors.New("validation error")
		} else {
			return nil, errors.New("invalid body")
		}
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(result)
	if err != nil {
		return nil, errors.New("validation error")
	}

	return &result, nil
}

func writeResponse[K any](w http.ResponseWriter, status int, result K) {
	w.WriteHeader(status)
	response, _ := json.Marshal(result)
	_, _ = w.Write(response)
}

func addCorsResponseHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
