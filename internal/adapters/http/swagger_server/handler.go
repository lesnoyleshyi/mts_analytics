package swagger_server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"strings"
	// path to swagger docs. Is required by swagger package
	_ "gitlab.com/g6834/team17/analytics-service/docs"
	"net"
	"net/http"
)

func (a SwaggerAdapter) routes() http.Handler {
	r := chi.NewRouter()

	basePath := net.JoinHostPort(host, strings.TrimPrefix(port, ":"))

	httpSwagger.UIConfig(map[string]string{
		"showExtensions":        "true",
		"onComplete":            `() => { window.ui.setBasePath('v3'); }`,
		"defaultModelRendering": `"model"`,
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", basePath))))

	return r
}
