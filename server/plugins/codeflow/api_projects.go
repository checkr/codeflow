package codeflow

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/extemporalgenome/slug"
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
	)

	log.Printf("Started the codeflow projects handler on %s\n", x.Path)
	return routes
}

func (x *Projects) projects(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var pageResults PageResults
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("projects_page"))
	itemsLimit := 20

	if pageResults, err = GetProjectsWithPagination(currentPage, itemsLimit); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(pageResults)
}

func (x *Projects) project(w rest.ResponseWriter, r *rest.Request) {
	project, err := CurrentProject(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(project)
}

func (x *Projects) createProjects(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var user User
	var project Project

	err = r.DecodeJsonPayload(&project)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user, err = CurrentUser(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = validate.Struct(project)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	gitSshUrlRegex, _ := regexp.Compile("(?:git|ssh|git@[\\w\\.]+):((?:\\/\\/)?[\\w\\.@:\\/~_-]+)\\.git(?:\\/?|\\#[\\d\\w\\.\\-_]+?)$")

	if !gitSshUrlRegex.MatchString(project.GitSshUrl) {
		rest.Error(w, "Wrong Git SSH url", http.StatusBadRequest)
		return
	}

	repository := gitSshUrlRegex.FindStringSubmatch(project.GitSshUrl)[1]
	project.Name = repository
	project.Repository = repository
	project.Slug = slug.Slug(repository)
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

	if err := CreateProject(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Bookmark
	if project.Bokmarked {
		if err := CreateUserBookmark(user.Id, project.Id); err != nil {
			log.Println(err.Error())
		} else {
			cf.BookmarksUpdated(&user)
		}
	}

	cf.ProjectCreated(&project)
	w.WriteJson(project)
}

func (x *Projects) services(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	services := []Service{}

	serviceCol := db.C("services")
	if err := serviceCol.Find(bson.M{"projectId": project.Id}).All(&services); err != nil {
	}

	w.WriteJson(services)
}

func (x *Projects) createServices(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var service Service

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.Id = bson.NewObjectId()
	service.Name = slug.Slug(service.Name)
	service.State = plugins.Waiting
	service.ProjectId = project.Id
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()

	if err := CreateService(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ServiceCreated(&service)
	w.WriteJson(service)
}

func (x *Projects) updateServices(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var service Service

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.ProjectId = project.Id
	service.State = plugins.Running
	service.UpdatedAt = time.Now()

	if err := UpdateService(service.Id, &service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ServiceUpdated(&service)
	w.WriteJson(service)
}

func (x *Projects) deleteServices(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var service Service

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if service.ProjectId != project.Id {
		rest.Error(w, fmt.Sprintf("Service does not belong to %v", project.Repository), http.StatusInternalServerError)
		return
	}

	service.State = plugins.Deleting

	if err := UpdateService(service.Id, &service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ServiceDeleted(&service)
	w.WriteJson(service)
}

func (x *Projects) extensions(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var extensions []LoadBalancer

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if extensions, err = GetExtensions(project.Id); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(extensions)
}

func (x *Projects) createExtensions(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var service Service

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extension := LoadBalancer{}

	err = r.DecodeJsonPayload(&extension)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if service, err = GetServiceById(extension.ServiceId); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extension.Id = bson.NewObjectId()
	lbIDTrimmed := strings.ToLower(extension.Id.Hex())[0:5]
	extension.Name = fmt.Sprintf("%v-%v-%v-%v", project.Slug, service.Name, extension.Type, lbIDTrimmed)
	extension.ProjectId = project.Id
	extension.CreatedAt = time.Now()
	extension.UpdatedAt = time.Now()

	if err := CreateExtension(&extension); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ExtensionCreated(&extension)
	w.WriteJson(extension)
}

func (x *Projects) updateExtensions(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extension := LoadBalancer{}

	err = r.DecodeJsonPayload(&extension)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extension.ProjectId = project.Id

	if err := UpdateExtension(extension.Id, &extension); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ExtensionUpdated(&extension)
	w.WriteJson(extension)
}

func (x *Projects) deleteExtensions(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	extension := LoadBalancer{}

	err = r.DecodeJsonPayload(&extension)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if extension.ProjectId != project.Id {
		rest.Error(w, fmt.Sprintf("Extension does not belong to %v", project.Repository), http.StatusInternalServerError)
		return
	}

	extension.State = plugins.Deleting

	if err := UpdateExtension(extension.Id, &extension); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ExtensionDeleted(&extension)

	w.WriteJson(extension)
}

func (x *Projects) features(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var pageResults PageResults
	var project Project
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("features_page"))
	itemsLimit := 5

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if pageResults, err = GetFeaturesWithPagination(project.Id, currentPage, itemsLimit); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(&pageResults)
}

func (x *Projects) releases(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var pageResults PageResults
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("releases_page"))
	itemsLimit := 5

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if pageResults, err = GetReleasesWithPagination(project.Id, currentPage, itemsLimit); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(&pageResults)
}

func (x *Projects) currentRelease(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var release Release
	var project Project

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if release, err = GetCurrentRelease(project.Id); err != nil {
		w.WriteJson(&release)
		return
	}

	w.WriteJson(&release)
}

func (x *Projects) createReleases(w rest.ResponseWriter, r *rest.Request) {
	var secrets []Secret
	var envSecrets []Secret
	var fileSecrets []Secret
	var headFeature Feature
	var tailFeature Feature

	project, _ := CurrentProject(r)
	user, _ := CurrentUser(r)

	err := r.DecodeJsonPayload(&headFeature)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if envSecrets, err = GetSecretsByProjectIdAndType(project.Id, plugins.Env); err != nil {
		log.Println(err.Error())
		return
	}

	if fileSecrets, err = GetSecretsByProjectIdAndType(project.Id, plugins.File); err != nil {
		log.Println(err.Error())
		return
	}

	secrets = append(secrets, envSecrets...)
	secrets = append(secrets, fileSecrets...)

	release := Release{
		ProjectId:     project.Id,
		HeadFeatureId: headFeature.Id,
		HeadFeature:   headFeature,
		UserId:        user.Id,
		User:          user,
		State:         plugins.Waiting,
		Secrets:       secrets,
	}

	if err := GetCurrentlyDeployedFeature(project.Id, &tailFeature); err != nil {
		release.TailFeatureId = headFeature.Id
		release.TailFeature = headFeature
	} else {
		release.TailFeatureId = tailFeature.Id
		release.TailFeature = tailFeature
	}

	if err := CreateRelease(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ReleaseCreated(&release)
	w.WriteJson(release)
}

func (x *Projects) createReleasesRollback(w rest.ResponseWriter, r *rest.Request) {
	var release Release

	err := r.DecodeJsonPayload(&release)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = PopulateRelease(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := CreateRelease(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.ReleaseCreated(&release)
	w.WriteJson(release)
}

func (x *Projects) settings(w rest.ResponseWriter, r *rest.Request) {
	var err error
	var project Project
	var secrets []Secret

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if secrets, err = GetSecretsByProjectId(project.Id); err != nil {
		log.Println(err.Error())
		return
	}

	settings := ProjectSettings{
		ProjectId: project.Id,
		GitSshUrl: project.GitSshUrl,
		Secrets:   secrets,
	}

	w.WriteJson(settings)
}

func (x *Projects) updateSettings(w rest.ResponseWriter, r *rest.Request) {
	var project Project
	var projectSettingsRequest ProjectSettings

	err := r.DecodeJsonPayload(&projectSettingsRequest)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if project, err = CurrentProject(r); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update project
	if projectSettingsRequest.GitSshUrl != "" && project.GitSshUrl != projectSettingsRequest.GitSshUrl {
		gitSshUrlRegex, _ := regexp.Compile("(?:git|ssh|git@[\\w\\.]+):((?:\\/\\/)?[\\w\\.@:\\/~_-]+)\\.git(?:\\/?|\\#[\\d\\w\\.\\-_]+?)$")

		if !gitSshUrlRegex.MatchString(project.GitSshUrl) {
			rest.Error(w, "Wrong Git SSH url", http.StatusBadRequest)
			return
		}

		repository := gitSshUrlRegex.FindStringSubmatch(projectSettingsRequest.GitSshUrl)[1]

		project.Name = repository
		project.Repository = repository
		project.Slug = slug.Slug(repository)
		project.GitSshUrl = projectSettingsRequest.GitSshUrl

		if err := UpdateProject(project.Id, &project); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Delete secrets
	for _, secret := range projectSettingsRequest.DeletedSecrets {
		if secret.Id.Valid() {
			secret.ProjectId = project.Id
			secret.Deleted = true
			secret.DeletedAt = time.Now()

			if err := UpdateSecret(secret.Id, &secret); err != nil {
				rest.Error(w, err.Error(), http.StatusBadRequest)
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
			if s, err = GetSecretById(secret.Id); err != nil {
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if s.Type == secret.Type && s.Key == secret.Key && s.Value == secret.Value {
				continue
			}

			if s.Type == "" {
				s.Type = plugins.Env
			}

			s.Deleted = true

			if err := UpdateSecret(secret.Id, &s); err != nil {
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		secret.Id = bson.NewObjectId()
		secret.ProjectId = project.Id
		secret.CreatedAt = time.Now()

		if err := CreateSecret(&secret); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	secrets := []Secret{}
	if secrets, err = GetSecretsByProjectId(project.Id); err != nil {
		log.Println(err.Error())
		return
	}

	projectSettingsResponse := ProjectSettings{
		ProjectId: projectSettingsRequest.ProjectId,
		GitSshUrl: projectSettingsRequest.GitSshUrl,
		Secrets:   secrets,
		UpdatedAt: time.Now(),
	}

	w.WriteJson(projectSettingsResponse)
}
