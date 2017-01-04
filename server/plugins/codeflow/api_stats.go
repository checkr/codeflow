package codeflow

import (
	"log"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

type Stats struct {
	Path string
}

func (x *Stats) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path, x.stats),
	)
	log.Printf("Started the codeflow stats handler on %s\n", x.Path)
	return routes
}

func (x *Stats) stats(w rest.ResponseWriter, r *rest.Request) {
	stats := Statistics{}

	if err := CollectStats(false, &stats); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(stats)
}
