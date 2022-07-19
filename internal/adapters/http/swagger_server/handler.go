package swagger_server

import (
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	// path to swagger docs. Is required by swagger package
	_ "gitlab.com/g6834/team17/analytics-service/api/swagger"
	"net/http"
)

func (a SwaggerAdapter) routes() http.Handler {
	r := chi.NewRouter()

	// TODO it works bad: every path, inserted in "explore" form is valid
	// perhaps because of /*
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))

	return r
}
