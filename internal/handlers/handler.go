package handlers

import (
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"mts_analytics/internal/domain"
	"net/http"
	"strings"
)

type filter interface {
	Match()
}

type agregate interface {
	GetFromRequest(r *http.Request)
}

type service interface {
	Save(event domain.Event) error
	//Get(filter)
	GetSignedCount() (int, error)
	GetNotSignedYetCount() (int, error)
	GetSignitionTotalTime(taskUUID string) (seconds int, err error)
}

type handler struct {
	s service
}

func New(s service) *handler {
	return &handler{s: s}
}

func (h handler) NewMux() http.Handler {
	r := chi.NewRouter()

	r.With(validateToken).Get("/signed", h.getSignedCount)

	return r
}

func (h handler) getSignedCount(w http.ResponseWriter, r *http.Request) {
	log.WithField("func", "handler.getSignedCount").Debug("request served")
}

func validateToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.WithField("func", "validateToken").Warningf("Empty auth token")
		} else {
			log.WithFields(log.Fields{
				"func":  "validateToken",
				"token": strings.TrimPrefix(token, "Bearer "),
			}).Debug("request served")
		}
		next.ServeHTTP(w, r)
	})
}
