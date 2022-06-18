package handlers

import (
	"fmt"
	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
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

func NewProfiler() http.Handler {
	r := chi.NewRouter()
	r.Use(validateToken)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		url := fmt.Sprintf("%s/debug/pprof", req.Host)
		_, _ = fmt.Fprintf(w, NotFoundPage, url, url)
	})

	r.Route("/debug/pprof", func(r chi.Router) {
		r.HandleFunc("/cmdline", pprof.Cmdline)
		r.HandleFunc("/profile", pprof.Profile)
		r.HandleFunc("/symbol", pprof.Symbol)
		r.HandleFunc("/trace", pprof.Trace)
		r.HandleFunc("/*", pprof.Index)
	})

	r.HandleFunc("/che", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("chto u vas proishodit")
	})

	return r
}
