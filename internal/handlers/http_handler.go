package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"mts_analytics/internal/domain"
	"net/http"
	"strings"
	"time"
)

const app_name = `analytics`
const host_ip = `127.0.0.1`

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
	r.Use(middleware.RequestID)

	r.With(validateToken).Get("/agreed", h.getSignedCount)
	r.With(validateToken).Get("/canceled", h.getNotSignedYetCount)
	r.With(validateToken).Get("/total_time", h.getSignitionTotalTime)

	return r
}

func (h handler) getSignedCount(w http.ResponseWriter, r *http.Request) {
	fn := "handlers_getSignedCount"
	count, err := h.s.GetSignedCount()
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":     app_name,
			"host_ip":      host_ip,
			"requestId":    middleware.GetReqID(r.Context()),
			"request_path": r.URL.Path,
			"logger_name":  fn,
		}).Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}
	w.Write([]byte(fmt.Sprintln(count)))
	log.WithField("func", fn).Debug("request served")
}

func (h handler) getNotSignedYetCount(w http.ResponseWriter, r *http.Request) {
	fn := "handlers_getNotSignedYetCount"
	count, err := h.s.GetNotSignedYetCount()
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":     app_name,
			"host_ip":      host_ip,
			"timestamp":    time.Now().Format(time.RFC3339),
			"requestId":    middleware.GetReqID(r.Context()),
			"request_path": r.URL.Path,
			"logger_name":  fn,
		}).Warn(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}
	w.Write([]byte(fmt.Sprintln(count)))
	log.WithField("func", fn).Debug("request served")
}

func (h handler) getSignitionTotalTime(w http.ResponseWriter, r *http.Request) {
	fn := "handler.getSignitionTotalTime"
	taskUUID := r.URL.Query().Get("id")
	if taskUUID == "" {
		log.WithFields(log.Fields{
			"app_name":     app_name,
			"host_ip":      host_ip,
			"timestamp":    time.Now().Format(time.RFC3339),
			"requestId":    middleware.GetReqID(r.Context()),
			"request_path": r.URL.Path,
			"logger_name":  fn,
		}).Debug("no task id provided")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provide task UUID as URL parameter"))
		return
	}
	t, err := h.s.GetSignitionTotalTime(taskUUID)
	if err != nil {
		log.WithFields(log.Fields{
			"app_name":     app_name,
			"host_ip":      host_ip,
			"timestamp":    time.Now().Format(time.RFC3339),
			"requestId":    middleware.GetReqID(r.Context()),
			"request_path": r.URL.Path,
			"logger_name":  fn,
		}).Debug(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error occurred"))
		return
	}
	w.Write([]byte(fmt.Sprintf("time in sec: %d", t)))
	log.WithField("func", fn).Debug("request served")
}

func validateToken(next http.Handler) http.Handler {
	fn := "handler_validateToken"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			log.WithFields(log.Fields{
				"app_name":     app_name,
				"host_ip":      host_ip,
				"timestamp":    time.Now().Format(time.RFC3339),
				"requestId":    middleware.GetReqID(r.Context()),
				"request_path": r.URL.Path,
				"logger_name":  fn,
			}).Debug("Empty auth token")
		} else {
			log.WithFields(log.Fields{
				"func":  "validateToken",
				"token": strings.TrimPrefix(token, "Bearer "),
			}).Debug("request served")
		}
		next.ServeHTTP(w, r)
	})
}
