package codeflow

import (
	"log"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

type MockEvents struct {
	Path string
}

func (x *MockEvents) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path+"/#event", x.mockEvents),
	)
	log.Printf("Started the codeflow user handler on %s\n", x.Path)
	return routes
}

func (x *MockEvents) mockEvents(w rest.ResponseWriter, r *rest.Request) {
	event := r.PathParam("event")
	response := `{"status": "ok"}`

	switch event {
	case "ProjectCreate":
		event := plugins.ProjectCreateMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "GitPing":
		event := plugins.GitPingMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "GitCommit":
		event := plugins.GitCommitMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "GitStatusCirclePending":
		event := plugins.GitStatusCirclePendingMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "GitStatusCircleFailed":
		event := plugins.GitStatusCircleFailedMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "GitStatusCircleSuccess":
		event := plugins.GitStatusCircleSuccessMock()
		cf.Events <- agent.NewEvent(event, nil)
	case "DockerDeployCreate":
		event := plugins.DockerDeployCreateMock()
		cf.Events <- agent.NewEvent(event, nil)
	default:
		response = `{"status": "not_found"}`
	}

	w.WriteJson(response)
}
