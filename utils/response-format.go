package util

import (
	"encoding/json"
	"net/http"
    "math/rand"
	"time"
)


type ErrorFormat struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}

func RandomGenerateOtp() int32{
    rand.Seed(time.Now().UnixNano())
	randomNumber := int32(rand.Intn(900000) + 100000)
    return randomNumber
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
    w.WriteHeader(statusCode)
    errResponse := ErrorFormat{
        Status:  statusCode,
        Message: message,
    }
    json.NewEncoder(w).Encode(errResponse)
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
