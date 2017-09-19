package codeflow

import (
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/plugins"
)

type Admin struct {
	Path string
}

func (x *Admin) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path+"/serviceSpecs", x.serviceSpecs),
		rest.Put(x.Path+"/serviceSpecs", x.updateServiceSpec),
		rest.Delete(x.Path+"/serviceSpecs/#slug", x.deleteServiceSpec),
		rest.Get(x.Path+"/serviceSpecs/#slug/services", x.serviceSpecServices),
		rest.Put(x.Path+"/deployallprojects", x.deployAllProjects),
		rest.Put(x.Path+"/deployloadbalancers", x.deployAllLoadBalancers),
	)
	log.Printf("Started the codeflow admin handler on %s\n", x.Path)
	return routes
}

// Deploys all projects (does not touch services)
func (x *Admin) deployAllProjects(w rest.ResponseWriter, r *rest.Request) {
	projects := []Project{}
	m := bson.M{"deleted": false}
	results := db.Collection("projects").Find(m)
	results.Query.Sort("-$natural").All(&projects)
	log.Printf("Deploying %d projects", len(projects))
	for _, p := range projects {
		log.Printf("Deploying %s...", p.Slug)
		release := Release{}
		if err := GetCurrentRelease(p.Id, &release); err != nil {
			log.Printf("GetCurrentRelease::Error: %s; Aborting release for %s", err, p.Slug)
		} else {
			log.Printf("Deploying %s release", release.Id)
			err := CreateDeploy(&release)
			if err != nil {
				log.Printf("ERROR '%s' creating deployment for %s", err, p.Slug)
			}
		}
	}
}

// Deploys all loadBalancers
func (x *Admin) deployAllLoadBalancers(w rest.ResponseWriter, r *rest.Request) {
	projects := []Project{}
	m := bson.M{"deleted": false}
	results := db.Collection("projects").Find(m)
	results.Query.Sort("-$natural").All(&projects)
	for _, p := range projects {
		log.Printf("Finding services for '%s' project...", p.Slug)
		l := bson.M{"projectId": p.Id, "state": "complete"}
		lbResults := db.Collection("extensions").Find(l)
		loadBalancer := &LoadBalancer{}
		for lbResults.Next(loadBalancer) {
			log.Printf("Deploying loadbalancer %s for %s", loadBalancer.Name, p.Slug)
			err := ExtensionUpdated(loadBalancer)
			if err != nil {
				log.Printf("ERROR '%s' updating loadBalancer %s", err, loadBalancer.Name)
			}
		}
	}
	// Throttle so that we can be "nice" to the amazon APIs
	log.Println("Sleeping 1 second")
	time.Sleep(time.Second * 1)
}

func (x *Admin) deleteServiceSpec(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")

	match := bson.M{"_id": bson.ObjectIdHex(slug)}

	if err := db.Collection("serviceSpecs").DeleteOne(match); err != nil {
		log.Printf("ServiceSpec:: Delete Error %s", err)
		return
	}
}

func (x *Admin) updateServiceSpec(w rest.ResponseWriter, r *rest.Request) {
	spec := ServiceSpec{}

	err := r.DecodeJsonPayload(&spec)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("serviceSpecs").Save(&spec); err != nil {
		log.Printf("ServiceSpec:: Save Error %s", err)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteJson(&spec)
}

func (x *Admin) serviceSpecs(w rest.ResponseWriter, r *rest.Request) {

	specs := []ServiceSpec{}
	spec := ServiceSpec{}

	results := db.Collection("serviceSpecs").Find(bson.M{})
	for results.Next(&spec) {
		specs = append(specs, spec)
	}

	w.WriteJson(specs)
}

func (x *Admin) serviceSpecServices(w rest.ResponseWriter, r *rest.Request) {
	services := []Service{}
	service := Service{}

	slug := r.PathParam("slug")
	match := bson.M{"specId": bson.ObjectIdHex(slug), "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running}}}

	results := db.Collection("services").Find(match)
	for results.Next(&service) {
		services = append(services, service)
	}

	w.WriteJson(services)
}
