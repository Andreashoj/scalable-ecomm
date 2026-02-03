package rest

import (
	"andreasho/scalable-ecomm/pgk/errors"
	"encoding/json"
	"fmt"
	"net/http"
)

func Response(w http.ResponseWriter, payload interface{}, statusCode int) error { // Worth abstracting this ?
	json, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed encoding payload: %s", err)
	}

	w.WriteHeader(statusCode)
	w.Write(json)
	return nil
}

func ErrorResponse(w http.ResponseWriter, statusCode int, errorMessage errors.ErrorMessage) {
	w.WriteHeader(statusCode)
	w.Write([]byte(errorMessage))
}
