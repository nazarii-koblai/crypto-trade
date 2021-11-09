package web

import (
	"encoding/json"
	"net/http"

	"github.com/crypto-trade/model"
)

const (
	contentType     = "Content-Type"
	applicationJSON = "application/json"
)

// RespondWithError used to log and respond with error.
func RespondWithError(w http.ResponseWriter, status int, err error) {
	w.Header().Add(contentType, applicationJSON)
	w.WriteHeader(status)
	b, _ := json.Marshal(model.ErrorResponse{
		Message: err.Error(),
	})
	w.Write(b)
}

// RespondWithSuccess used to log and respond with success.
func RespondWithSuccess(w http.ResponseWriter, status int, body []byte) {
	w.Header().Add(contentType, applicationJSON)
	w.WriteHeader(status)
	w.Write(body)
}
