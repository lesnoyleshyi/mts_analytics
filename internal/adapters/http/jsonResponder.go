package http

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type JSONResponder struct {
	logger *zap.Logger
}

func NewJSONResponder(logger *zap.Logger) JSONResponder {
	return JSONResponder{logger: logger}
}

func (r JSONResponder) RespondSuccess(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "{\"result\":\"%s\"}", msg); err != nil {
		r.logger.Warn("error writing response", zap.Error(err))
	}
}

func (r JSONResponder) RespondError(w http.ResponseWriter, msg string, status int) {
	http.Error(w, fmt.Sprintf("{\"error\":\"%s\"}", msg), status)
}
