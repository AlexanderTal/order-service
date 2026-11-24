package httph

import (
	"encoding/json"
	"errors"
	"net/http"
)

func SendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func SendEmpty(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorApply(r, err)

	var hc interface {
		HTTPStatus() int
	}
	if errors.As(err, &hc) {
		status := hc.HTTPStatus()
		ErrorApplyStatusCode(r, status)
		sendError(w, status, err)
		return
	}
	sendError(w, http.StatusInternalServerError, err)
}

func sendError(w http.ResponseWriter, status int, err error) {
	SendJSON(w, status, map[string]string{"error": err.Error()})
}
