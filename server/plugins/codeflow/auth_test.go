package codeflow

import (
	"net/http"
	"testing"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ant0ine/go-json-rest/rest/test"
	"github.com/jarcoal/httpmock"
)

func TestOktaToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://.okta.com/oauth2/v1/keys",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, OktaKeys())
			if err != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		},
	)

	auth := Auth{}
	api := rest.NewApi()

	var routes []*rest.Route
	var handlerRoutes []*rest.Route

	handlerRoutes = auth.Register(api)
	routes = append(routes, handlerRoutes...)

	router, err := rest.MakeRouter(routes...)
	if err != nil {
		t.Fatal(err.Error())
	}

	api.SetApp(router)

	handler := api.MakeHandler()
	if handler == nil {
		t.Fatal("the http.Handler must have been created")
	}

	recorded := test.RunRequest(t, handler, test.MakeSimpleRequest("POST", "/callback/okta", OktaIdToken()))
	recorded.CodeIs(400)
}
