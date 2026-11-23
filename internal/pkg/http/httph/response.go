package httph

import (
	"encoding/json"
	"errors"
	"net/http"
)

type httpCoder interface {
	error
	HTTPStatus() int
}

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
	var hc httpCoder
	if errors.As(err, &hc) {
		ErrorApply(r, hc)
		ErrorApplyStatusCode(r, hc.HTTPStatus())
		sendError(w, hc)
		return
	}

	defaultErr := &defaultError{err: err}
	ErrorApply(r, defaultErr)
	ErrorApplyStatusCode(r, defaultErr.HTTPStatus())
	sendError(w, defaultErr)
}

func sendError(w http.ResponseWriter, err httpCoder) {
	SendJSON(w, err.HTTPStatus(), map[string]string{"error": err.Error()})
}

type defaultError struct {
	err error
}

func (e *defaultError) Error() string {
	return "internal server error"
}

func (e *defaultError) HTTPStatus() int {
	return http.StatusInternalServerError
}
