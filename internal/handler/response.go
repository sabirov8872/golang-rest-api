package handler

import (
	"encoding/json"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, data interface{}) {
	out, err := json.Marshal(data)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)
	if err != nil {
		newErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}
}

func newErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}
