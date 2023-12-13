package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/amrohan/expenso-go/internal/models"
)

func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, err error) {
	w.WriteHeader(statusCode)
	res := models.Response{Status: statusCode, Message: message, Data: nil}

	if err != nil {
		res.Data = err.Error()
	} else if data != nil {
		res.Data = data
	}
	json.NewEncoder(w).Encode(res)
}
