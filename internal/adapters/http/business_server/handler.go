package business_server

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

// getSignedCount godoc
// @Tags count
// @Summary get signed tasks count
// @Description Get count of all signed tasks.
// @Produce plain
// @Success 200 {integer} int
// @Failure 400 {string} string
//nolint:godot
// @Router /agreed [get]
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

// getUnsignedCount godoc
// @Tags count
// @Summary get unsigned tasks count
// @Description Get count of all tasks which are not signed (rejected and "in process").
// @Produce plain
// @Success 200 {integer} int
// @Failure 400 {string} string
//nolint:godot
// @Router /canceled [get]
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

// getSignitionTime godoc
// @Tags time
// @Summary get signition time of task
// @Description Get total signition time in seconds of particular task by its id
// @Produce plain
// TODO inadequate representation in swagger UI! I don't want 'string' in 'Example value'
// @Success 200 {string} string "example: 'time in sec: 100500'"
// @Failure 400 {string} string "returns in case user input is invalid"
// @Failure 500 {string} string "returns in case server can't retrieve signition time"
// @Router /total_time [get]
//nolint:godot
// @Param id query int true "uuid of task"
func (a AdapterHTTP) getSignitionTime(w http.ResponseWriter, r *http.Request) {
	taskUUID := r.URL.Query().Get("id")
	if taskUUID == "" {
		a.logger.Warn("no task id provided")
		a.responder.RespondError(w, "provide task UUID as URL parameter in form ?id=1",
			http.StatusBadRequest)
		return
	}

	t, err := a.events.GetSignitionTime(context.TODO(), entity.Event{TaskUUID: taskUUID}) //nolint:exhaustruct
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
