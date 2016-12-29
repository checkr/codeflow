package codeflow

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
)

type Users struct {
	Path string
}

func (x *Users) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path+"/me", x.me),
	)
	log.Printf("Started the codeflow user handler on %s\n", x.Path)
	return routes
}

func (x *Users) me(w rest.ResponseWriter, r *rest.Request) {
	user, _ := CurrentUser(r)
	w.WriteJson(user)
}
