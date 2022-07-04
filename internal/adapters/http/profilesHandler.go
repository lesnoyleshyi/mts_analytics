package http

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/pprof"
)

const NotFoundPage = `
<!DOCTYPE html>
<html>
<body>
	<p>Unknown Url</p>
	<p>Visit <a href="http://%s">%s</a> to get all available profiles</p>
</body>
</html>
`

func (a ProfileAdapter) routeProfiles() http.Handler {
	r := chi.NewRouter()

	//r.Handle("/debug/loglevel", app.DynamicLogLevel)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		url := fmt.Sprintf("%s/debug/pprof/", req.Host)
		_, _ = fmt.Fprintf(w, NotFoundPage, url, url)
	})

	r.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.HandleFunc("/*", pprof.Index)
	})

	return r
}
