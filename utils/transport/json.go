package transport

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, resp interface{}, httpCode int) {
	responseBody, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(responseBody)
}

func JsonResponseError(w http.ResponseWriter, resp interface{}, err error, httpCode int) {
	responseBody, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshalErr.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(responseBody)
}