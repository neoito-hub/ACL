package common_services

import (
	"encoding/json"
	"net/http"
)

func BuildErrorResponse(err bool, msg string, data interface{}, status int) HandlerResponse {
	var handlerResp HandlerResponse

	byteSlice, _ := json.Marshal(data)
	stringBody := string(byteSlice)

	handlerResp.Err = err
	handlerResp.Status = status
	handlerResp.Data = stringBody

	return handlerResp

}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dataPayload, err := json.Marshal(payload)
	if err != nil {
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	w.Write(dataPayload)

}

func BuildResponse(err bool, msg string, data interface{}, status int) HandlerResponse {
	var handlerResp HandlerResponse

	byteSlice, _ := json.Marshal(data)
	stringBody := string(byteSlice)

	handlerResp.Err = err
	handlerResp.Status = status
	handlerResp.Data = stringBody

	return handlerResp

}
