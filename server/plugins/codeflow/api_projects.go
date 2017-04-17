package codeflow

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/davecgh/go-spew/spew"
	"github.com/extemporalgenome/slug"
	"github.com/maxwellhealth/bongo"
	"golang.org/x/crypto/ssh"
)

type Projects struct {
	Path string
}

func (x *Projects) Register(api *rest.Api) []*rest.Route {
	var routes []*rest.Route
	routes = append(routes,
		rest.Get(x.Path, x.projects),
		rest.Get(x.Path+"/#slug/services", x.services),
		rest.Post(x.Path+"/#slug/services", x.createServices),
		rest.Put(x.Path+"/#slug/services", x.updateServices),
		rest.Delete(x.Path+"/#slug/services", x.deleteServices),
		rest.Get(x.Path+"/#slug/extensions", x.extensions),
		rest.Post(x.Path+"/#slug/extensions", x.createExtensions),
		rest.Put(x.Path+"/#slug/extensions", x.updateExtensions),
		rest.Delete(x.Path+"/#slug/extensions", x.deleteExtensions),
		rest.Get(x.Path+"/#slug/settings", x.settings),
		rest.Put(x.Path+"/#slug/settings", x.updateSettings),
		rest.Post(x.Path+"/#slug/releases", x.createReleases),
		rest.Get(x.Path+"/#slug/releases", x.releases),
		rest.Get(x.Path+"/#slug/releases/current", x.currentRelease),
		rest.Post(x.Path+"/#slug/releases/rollback", x.createReleasesRollback),
		rest.Get(x.Path+"/#slug/features", x.features),
		rest.Get(x.Path+"/#slug", x.project),
		rest.Post(x.Path, x.createProjects),
		rest.Get(x.Path+"/#slug/releases/#id/build", x.releaseBuild),
		rest.Post(x.Path+"/#slug/releases/#id/build", x.updateReleaseBuild),
		rest.Get(x.Path+"/#slug/serviceSpecs", x.serviceSpecs),
	)

	log.Printf("Started the codeflow projects handler on %s\n", x.Path)
	return routes
}

func (x *Projects) serviceSpecs(w rest.ResponseWriter, r *rest.Request) {
	specs := []ServiceSpec{}
	spec := ServiceSpec{}

	results := db.Collection("serviceSpecs").Find(bson.M{})
	for results.Next(&spec) {
		specs = append(specs, spec)
	}

	w.WriteJson(specs)
}

func (x *Projects) projects(w rest.ResponseWriter, r *rest.Request) {
	projects := []Project{}
	pageResults := PageResults{}
	query := r.URL.Query()
	search := query.Get("q")
	currentPage, _ := strconv.Atoi(query.Get("projects_page"))
	perPage := 20
	m := bson.M{}

	if search != "" {
		m = bson.M{"name": bson.RegEx{search, "i"}}
	}

	results := db.Collection("projects").Find(m)
	if pagination, err := results.Paginate(perPage, currentPage); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		pageResults.Pagination = *pagination
	}

	results.Query.Sort("-$natural").All(&projects)
	pageResults.Records = projects

	w.WriteJson(pageResults)
}

func (x *Projects) project(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteJson(project)
}

func (x *Projects) createProjects(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	project := Project{}
	bookmark := Bookmark{}

	if err := r.DecodeJsonPayload(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := CurrentUser(r, &user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	project.Secret = agent.RandomString(30)

	// priv *rsa.PrivateKey;
	priv, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = priv.Validate()
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get der format. priv_der []byte
	priv_der := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}

	// Public Key generation
	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	project.RsaPrivateKey = string(pem.EncodeToMemory(&priv_blk))
	project.RsaPublicKey = string(ssh.MarshalAuthorizedKey(pub))

	// TODO: make this dynamic based on type of the project
	project.Workflows = []string{"build/DockerImage"}

	if err := db.Collection("projects").Save(&project); err != nil {
		if vErr, ok := err.(*bongo.ValidationError); ok {
			rest.Error(w, "Validation error", http.StatusBadRequest)
			spew.Dump(vErr)
			return
		} else {
			log.Printf("Projects::Save::Error: %v", err.Error())
		}
		return
	}

	// Bookmark
	if project.Bokmarked {
		bookmark.ProjectId = project.Id
		bookmark.UserId = user.Id

		if err := db.Collection("bookmarks").Save(&bookmark); err != nil {
			log.Printf("Bookmarks::Save::Error: %v", err.Error())
		} else {
			BookmarksUpdated(&user)
		}
	}

	if project.GitProtocol == "HTTPS" {
		GitSyncProjects([]bson.ObjectId{project.Id})
	}

	ProjectCreated(&project)

	w.WriteJson(project)
}

func (x *Projects) services(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	services := []Service{}
	service := Service{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := db.Collection("services").Find(bson.M{"projectId": project.Id, "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running}}})
	for results.Next(&service) {
		services = append(services, service)
	}

	w.WriteJson(services)
}

func (x *Projects) createServices(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	service := Service{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.Name = slug.Slug(service.Name)
	service.State = plugins.Waiting
	service.ProjectId = project.Id

	if err := db.Collection("services").Save(&service); err != nil {
		log.Printf("Services::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ServiceCreated(&service)

	w.WriteJson(service)
}

func (x *Projects) updateServices(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	service := Service{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.ProjectId = project.Id

	if err := db.Collection("services").Save(&service); err != nil {
		log.Printf("Services::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ServiceUpdated(&service)

	w.WriteJson(service)
}

func (x *Projects) deleteServices(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	service := Service{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.State = plugins.Deleted

	if err := db.Collection("services").Save(&service); err != nil {
		log.Printf("Services::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ServiceDeleted(&service)

	w.WriteJson(service)
}

func (x *Projects) extensions(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	extensions := []LoadBalancer{}
	extension := LoadBalancer{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := db.Collection("extensions").Find(bson.M{"projectId": project.Id, "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running, plugins.Complete, plugins.Failed}}})
	for results.Next(&extension) {
		extensions = append(extensions, extension)
	}

	w.WriteJson(extensions)
}

func (x *Projects) createExtensions(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	extension := LoadBalancer{}
	service := Service{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&extension); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("services").FindById(extension.ServiceId, &service); err != nil {
		log.Printf("Services::FindById::Error: %v, serviceId: %v", err.Error(), extension.ServiceId)
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extension.Id = bson.NewObjectId()
	extension.State = plugins.Waiting
	extension.Name = fmt.Sprintf("%v-%v", extension.Type, extension.Id.Hex())
	extension.ProjectId = project.Id

	if err := db.Collection("extensions").Save(&extension); err != nil {
		log.Printf("Extensions::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ExtensionCreated(&extension)

	w.WriteJson(extension)
}

func (x *Projects) updateExtensions(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	extension := LoadBalancer{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&extension); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extension.ProjectId = project.Id

	if err := db.Collection("extensions").Save(&extension); err != nil {
		log.Printf("Extensions::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ExtensionUpdated(&extension)

	w.WriteJson(extension)
}

func (x *Projects) deleteExtensions(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	extension := LoadBalancer{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&extension); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extension.State = plugins.Deleting

	if err := db.Collection("extensions").Save(&extension); err != nil {
		log.Printf("Extensions::Save::Error: %v", err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ExtensionDeleted(&extension)

	w.WriteJson(extension)
}

func (x *Projects) features(w rest.ResponseWriter, r *rest.Request) {
	features := []Feature{}
	feature := Feature{}
	project := Project{}
	pageResults := PageResults{}
	currentRelease := Release{}
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("features_page"))
	perPage := 5

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	query := bson.M{"projectId": project.Id}
	if err := GetCurrentRelease(project.Id, &currentRelease); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			// nothing to do
		} else {
			log.Printf("GetCurrentRelease::Error: %s", err.Error())
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		// Find only undeployed features
		query = bson.M{"projectId": project.Id, "created": bson.M{"$gt": currentRelease.HeadFeature.Created}}
	}

	results := db.Collection("features").Find(query)

	if pagination, err := results.Paginate(perPage, currentPage); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		pageResults.Pagination = *pagination
	}

	results.Query.Sort("-created")
	for results.Next(&feature) {
		features = append(features, feature)
	}

	pageResults.Records = features

	w.WriteJson(pageResults)
}

func (x *Projects) releases(w rest.ResponseWriter, r *rest.Request) {
	releases := []Release{}
	release := Release{}
	project := Project{}
	pageResults := PageResults{}
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("releases_page"))
	perPage := 5

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results := db.Collection("releases").Find(bson.M{"projectId": project.Id})
	if pagination, err := results.Paginate(perPage, currentPage); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		pageResults.Pagination = *pagination
	}

	results.Query.Sort("-$natural")
	for results.Next(&release) {
		releases = append(releases, release)
	}

	pageResults.Records = releases

	w.WriteJson(pageResults)
}

func (x *Projects) currentRelease(w rest.ResponseWriter, r *rest.Request) {
	release := Release{}
	project := Project{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := GetCurrentRelease(project.Id, &release); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			// Nothing to do here
		} else {
			log.Printf("GetCurrentRelease::Error: %s", err.Error())
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteJson(&release)
}

func (x *Projects) createReleases(w rest.ResponseWriter, r *rest.Request) {
	user := User{}
	headFeature := Feature{}
	release := Release{}

	if err := CurrentUser(r, &user); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.DecodeJsonPayload(&headFeature); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := NewRelease(headFeature, user, &release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := ReleaseCreated(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(release)
}

func (x *Projects) createReleasesRollback(w rest.ResponseWriter, r *rest.Request) {
	release := Release{}

	if err := r.DecodeJsonPayload(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("releases").FindById(release.Id, &release); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Releases::FindById::DocumentNotFoundError: _id: `%v`", release.Id)
		} else {
			log.Printf("Releases::FindById::Error: %s", err.Error())
		}
		return
	}

	release.Id = bson.NewObjectId()
	release.State = plugins.Waiting

	if err := db.Collection("releases").Save(&release); err != nil {
		log.Printf("Releases::Save::Error: %v", err.Error())
		log.Printf("Releases::FindById::Error: %s", err.Error())
		return
	}

	if err := ReleaseCreated(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(release)
}

func (x *Projects) settings(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	secrets := []Secret{}
	secret := Secret{}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	results := db.Collection("secrets").Find(bson.M{"projectId": project.Id, "deleted": false})
	results.Query.Sort("$natural")
	for results.Next(&secret) {
		secrets = append(secrets, secret)
	}

	settings := ProjectSettings{
		ProjectId:             project.Id,
		GitUrl:                project.GitUrl,
		GitProtocol:           project.GitProtocol,
		Secrets:               secrets,
		ContinuousIntegration: project.ContinuousIntegration,
		ContinuousDelivery:    project.ContinuousDelivery,
		NotifyChannels:        project.NotifyChannels,
	}

	w.WriteJson(settings)
}

func (x *Projects) updateSettings(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	projectSettingsRequest := ProjectSettings{}

	err := r.DecodeJsonPayload(&projectSettingsRequest)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
	}

	project.GitUrl = projectSettingsRequest.GitUrl
	project.GitProtocol = projectSettingsRequest.GitProtocol

	project.ContinuousIntegration = projectSettingsRequest.ContinuousIntegration
	if project.ContinuousIntegration {
		for idx, wf := range project.Workflows {
			if strings.HasPrefix(wf, "ci/") {
				project.Workflows = append(project.Workflows[:idx], project.Workflows[idx+1:]...)
			}
		}
		// TODO: make this dynamic
		project.Workflows = append(project.Workflows, "ci/circleci")
	} else {
		for idx, wf := range project.Workflows {
			if strings.HasPrefix(wf, "ci/") {
				project.Workflows = append(project.Workflows[:idx], project.Workflows[idx+1:]...)
			}
		}
	}

	project.ContinuousDelivery = projectSettingsRequest.ContinuousDelivery

	project.NotifyChannels = projectSettingsRequest.NotifyChannels

	if err := db.Collection("projects").Save(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete secrets
	for _, secret := range projectSettingsRequest.DeletedSecrets {
		if secret.Id.Valid() {
			secret.ProjectId = project.Id
			secret.Deleted = true

			if err := db.Collection("secrets").Save(&secret); err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Update secrets
	for _, secret := range projectSettingsRequest.Secrets {
		if secret.Value == "" {
			continue
		}

		if secret.Id.Valid() {
			s := Secret{}

			if err := db.Collection("secrets").FindById(secret.Id, &s); err != nil {
				if _, ok := err.(*bongo.DocumentNotFoundError); ok {
					log.Printf("Secrets::FindById::DocumentNotFoundError: _id: `%v`", secret.Id)
				} else {
					log.Printf("Secrets::FindById::Error: %s", err.Error())
				}
				rest.Error(w, err.Error(), http.StatusBadRequest)
			}

			if s.Type == secret.Type && s.Key == secret.Key && s.Value == secret.Value {
				continue
			}

			if s.Type == "" {
				s.Type = plugins.Env
			}

			s.Deleted = true

			if err := db.Collection("secrets").Save(&s); err != nil {
				rest.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		secret.Id = bson.NewObjectId()
		secret.ProjectId = project.Id

		if err := db.Collection("secrets").Save(&secret); err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	secrets := []Secret{}
	secret := Secret{}
	results := db.Collection("secrets").Find(bson.M{"projectId": project.Id, "deleted": false})
	results.Query.Sort("$natural")
	for results.Next(&secret) {
		secrets = append(secrets, secret)
	}

	projectSettingsResponse := ProjectSettings{
		ProjectId:             projectSettingsRequest.ProjectId,
		GitUrl:                projectSettingsRequest.GitUrl,
		GitProtocol:           projectSettingsRequest.GitProtocol,
		ContinuousIntegration: project.ContinuousIntegration,
		ContinuousDelivery:    project.ContinuousDelivery,
		Secrets:               secrets,
		NotifyChannels:        project.NotifyChannels,
	}

	w.WriteJson(projectSettingsResponse)
}

func (x *Projects) releaseBuild(w rest.ResponseWriter, r *rest.Request) {
	release := Release{}
	project := Project{}
	build := Build{}
	releaseId := r.PathParam("id")

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("releases").FindById(bson.ObjectIdHex(releaseId), &release); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Releases::FindOne::DocumentNotFoundError: releaseId: `%v`", releaseId)
		} else {
			log.Printf("Releases::FindOne::Error: %s", err.Error())
		}

		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("builds").FindOne(bson.M{"featureHash": release.HeadFeature.Hash}, &build); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Builds::FindOne::DocumentNotFoundError: featureHash: `%v`", release.HeadFeature.Hash)
		} else {
			log.Printf("Builds::FindOne::Error: %s", err.Error())
		}
	}

	w.WriteJson(build)
}

func (x *Projects) updateReleaseBuild(w rest.ResponseWriter, r *rest.Request) {
	release := Release{}
	project := Project{}
	build := Build{}
	releaseId := r.PathParam("id")

	if err := CurrentProject(r, &project); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("releases").FindById(bson.ObjectIdHex(releaseId), &release); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Releases::FindOne::DocumentNotFoundError: releaseId: `%v`", releaseId)
		} else {
			log.Printf("Releases::FindOne::Error: %s", err.Error())
		}

		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	release.State = plugins.Waiting
	if err := db.Collection("releases").Save(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := db.Collection("builds").FindOne(bson.M{"featureHash": release.HeadFeature.Hash}, &build); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Builds::FindOne::DocumentNotFoundError: featureHash: `%v`", release.HeadFeature.Hash)
		} else {
			log.Printf("Builds::FindOne::Error: %s", err.Error())
		}
	}

	DockerBuildRebuild(&release)

	w.WriteJson(build)
}
