package view

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Response represents response model that will be converted to json.
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// RespondWithJSON writes json response format.
func RespondWithJSON(w http.ResponseWriter, statusCode int, statusMessage string, data interface{}) {
	if statusCode/100 != 2 {
		statusMessage = http.StatusText(statusCode) + ":" + statusMessage
	}

	responseJSON, _ := json.Marshal(Response{
		Status:  statusCode,
		Message: statusMessage,
		Data:    data,
	})

	// Set response header.
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseJSON)))
	w.WriteHeader(statusCode)

	_, _ = w.Write(responseJSON)
}
