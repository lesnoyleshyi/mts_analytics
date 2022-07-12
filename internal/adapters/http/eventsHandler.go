package http

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gitlab.com/g6834/team17/analytics-service/internal/domain/entity"
	"go.uber.org/zap"
	"net/http"
)

func (a AdapterHTTP) routeEvents() http.Handler {
	r := chi.NewRouter()

	r.Get("/agreed", a.getSignedCount)
	r.Get("/canceled", a.getUnsignedCount)
	r.Get("/total_time", a.getSignitionTime)

	return r
}

func (a AdapterHTTP) getSignedCount(w http.ResponseWriter, r *http.Request) {
	// TODO pass correct context here
	count, err := a.events.GetSignedCount(context.TODO())
	if err != nil {
		a.logger.Error("can't get count of signed tasks", zap.Error(err))
		a.responder.RespondError(w, "error receiving count of signed tasks",
			http.StatusBadRequest)

		return
	}

	a.responder.RespondSuccess(w, fmt.Sprint(count), http.StatusOK)
}

func (a AdapterHTTP) getUnsignedCount(w http.ResponseWriter, r *http.Request) {
	count, err := a.events.GetUnsignedCount(context.TODO())
	if err != nil {
		a.logger.Error("can't get count of unsigned tasks", zap.Error(err))
		a.responder.RespondError(w, "error receiving count of unsigned tasks",
			http.StatusBadRequest)

		return
	}

	a.responder.RespondSuccess(w, fmt.Sprint(count), http.StatusOK)

	a.logger.Debug("request served")
}

func (a AdapterHTTP) getSignitionTime(w http.ResponseWriter, r *http.Request) {
	taskUUID := r.URL.Query().Get("id")
	if taskUUID == "" {
		a.logger.Warn("no task id provided")
		a.responder.RespondError(w, "provide task UUID as URL parameter in form ?id=1",
			http.StatusBadRequest)
		return
	}

	t, err := a.events.GetSignitionTime(context.TODO(), entity.Event{TaskUUID: taskUUID})
	if err != nil {
		a.logger.Error("can't get count of unsigned tasks", zap.Error(err))
		// not all errors should be 5XX. In case of wrong UUID it should return 400
		a.responder.RespondError(w, "error receiving signition time",
			http.StatusInternalServerError)
		return
	}

	a.responder.RespondSuccess(w, fmt.Sprintf("time in sec: %d", t), http.StatusOK)

	a.logger.Debug("request served")
}
