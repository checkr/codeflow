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

func (x *Projects) CurrentProject(r *rest.Request) (Project, error) {
	project := Project{}
	slug := r.PathParam("slug")

	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		return project, err
	}

	return project, nil
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
	projects := []Project{}
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("projects_page"))
	itemsCount := 0
	itemsLimit := 20

	collection := db.C("projects")

	itemsCount, err := collection.Find(bson.M{}).Count()
	if err != nil {
	}

	p := Pagination{
		Page:  currentPage,
		Limit: itemsLimit,
		Count: itemsCount,
	}.Calc()

	if err := collection.Find(bson.M{}).Limit(p.Limit).Skip(p.Offset).All(&projects); err != nil {
	}

	pageResults := PageResults{
		Records:    &projects,
		Pagination: p,
	}

	w.WriteJson(pageResults)
}

func (x *Projects) project(w rest.ResponseWriter, r *rest.Request) {
	project, err := x.CurrentProject(r)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteJson(project)
}

func (x *Projects) createProjects(w rest.ResponseWriter, r *rest.Request) {
	project := Project{}
	project.Id = bson.NewObjectId()

	err := r.DecodeJsonPayload(&project)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
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

	collection := db.C("projects")
	if err := collection.Insert(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Bookmark
	if project.Bokmarked {
		user, _ := CurrentUser(r)
		bookmark := Bookmark{
			ProjectId: project.Id,
			UserId:    user.Id,
		}
		bookmarkCol := db.C("bookmarks")
		if _, err := bookmarkCol.Upsert(bson.M{"userId": user.Id, "projectId": project.Id}, &bookmark); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		cf.UpdateBookmarks(&user)
	}

	go cf.CreateProject(&project)
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
	s := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"slug": s}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := Service{}

	err := r.DecodeJsonPayload(&service)
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

	if err := serviceCol.Insert(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go cf.CreateService(&service)
	w.WriteJson(service)
}

func (x *Projects) updateServices(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := Service{}

	err := r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.ProjectId = project.Id
	service.State = plugins.Running
	service.UpdatedAt = time.Now()

	if err := serviceCol.Update(bson.M{"_id": service.Id}, &service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go cf.UpdateService(&service)
	w.WriteJson(service)
}

func (x *Projects) deleteServices(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	service := Service{}

	err := r.DecodeJsonPayload(&service)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.State = plugins.Deleting
	service.UpdatedAt = time.Now()

	if err := serviceCol.Update(bson.M{"_id": service.Id}, &service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.DeleteService(&service)

	w.WriteJson(service)
}

func (x *Projects) extensions(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	extensionCol := db.C("extensions")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadBalancers := []LoadBalancer{}

	if err := extensionCol.Find(bson.M{"projectId": project.Id}).All(&loadBalancers); err != nil {
	}

	w.WriteJson(loadBalancers)
}

func (x *Projects) createExtensions(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}
	service := Service{}

	projectCol := db.C("projects")
	extensionCol := db.C("extensions")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadBalancer := LoadBalancer{}

	err := r.DecodeJsonPayload(&loadBalancer)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := serviceCol.Find(bson.M{"_id": loadBalancer.ServiceId}).One(&service); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadBalancer.Id = bson.NewObjectId()
	lbIDTrimmed := strings.ToLower(loadBalancer.Id.Hex())[0:5]
	loadBalancer.Name = fmt.Sprintf("%v-%v-%v-%v", project.Slug, service.Name, loadBalancer.Type, lbIDTrimmed)
	loadBalancer.ProjectId = project.Id
	loadBalancer.CreatedAt = time.Now()
	loadBalancer.UpdatedAt = time.Now()

	if err := extensionCol.Insert(&loadBalancer); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go cf.CreateExtensions(&loadBalancer)
	w.WriteJson(loadBalancer)
}

func (x *Projects) updateExtensions(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	extensionCol := db.C("extensions")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadBalancer := LoadBalancer{}

	err := r.DecodeJsonPayload(&loadBalancer)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loadBalancer.ProjectId = project.Id
	loadBalancer.UpdatedAt = time.Now()

	if err := extensionCol.Update(bson.M{"_id": loadBalancer.Id}, &loadBalancer); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go cf.UpdateExtensions(&loadBalancer)
	w.WriteJson(loadBalancer)
}

func (x *Projects) deleteExtensions(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}

	projectCol := db.C("projects")
	extensionCol := db.C("extensions")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	loadBalancer := LoadBalancer{}

	err := r.DecodeJsonPayload(&loadBalancer)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cf.DeleteExtensions(&loadBalancer)

	extensionCol.Remove(bson.M{"_id": loadBalancer.Id})

	w.WriteJson(loadBalancer)
}

func (x *Projects) features(w rest.ResponseWriter, r *rest.Request) {
	features := []Feature{}
	slug := r.PathParam("slug")
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("features_page"))
	itemsLimit := 5

	projectCol := db.C("projects")
	featureCol := db.C("features")

	project := Project{}
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p *Pagination

	currentlyDeployedFeature := Feature{}
	if ok := x.findCurrentlyDeployedFeature(&project, &currentlyDeployedFeature); !ok {
		itemsCount, err := featureCol.Find(bson.M{"projectId": project.Id}).Count()
		if err != nil {

		}

		p = Pagination{
			Page:  currentPage,
			Limit: itemsLimit,
			Count: itemsCount,
		}.Calc()

		if err := featureCol.Find(bson.M{"projectId": project.Id}).Limit(itemsLimit).Skip(p.Offset).Sort("-$natural").All(&features); err != nil {

		}
	} else {
		itemsCount, err := featureCol.Find(bson.M{"projectId": project.Id, "_id": bson.M{"$gt": currentlyDeployedFeature.Id}}).Count()
		if err != nil {

		}

		p = Pagination{
			Page:  currentPage,
			Limit: itemsLimit,
			Count: itemsCount,
		}.Calc()

		if err := featureCol.Find(bson.M{"projectId": project.Id, "_id": bson.M{"$gt": currentlyDeployedFeature.Id}}).Limit(p.Limit).Skip(p.Offset).Sort("-$natural").All(&features); err != nil {

		}
	}

	pageResults := PageResults{
		Records:    &features,
		Pagination: p,
	}

	w.WriteJson(&pageResults)
}

func (x *Projects) releases(w rest.ResponseWriter, r *rest.Request) {
	releases := []Release{}

	slug := r.PathParam("slug")
	currentPage, _ := strconv.Atoi(r.URL.Query().Get("releases_page"))
	itemsCount := 0
	itemsLimit := 5

	projectCol := db.C("projects")
	releaseCol := db.C("releases")
	workflowCol := db.C("workflows")
	featureCol := db.C("features")
	userCol := db.C("users")

	project := Project{}
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itemsCount, err := releaseCol.Find(bson.M{"projectId": project.Id}).Count()
	if err != nil {

	}

	p := Pagination{
		Page:  currentPage,
		Limit: itemsLimit,
		Count: itemsCount,
	}.Calc()

	if err := releaseCol.Find(bson.M{"projectId": project.Id}).Limit(p.Limit).Skip(p.Offset).Sort("-$natural").All(&releases); err != nil {

	}

	for idx, rel := range releases {
		headFeature := Feature{}
		tailFeature := Feature{}
		user := User{}

		if err := featureCol.Find(bson.M{"_id": rel.HeadFeatureId}).One(&headFeature); err != nil {

		}

		if err := featureCol.Find(bson.M{"_id": rel.TailFeatureId}).One(&tailFeature); err != nil {

		}

		if err := userCol.Find(bson.M{"_id": rel.UserId}).One(&user); err != nil {

		}

		releases[idx].HeadFeature = headFeature
		releases[idx].TailFeature = tailFeature
		releases[idx].User = user

		workflows := []Flow{}
		if err := workflowCol.Find(bson.M{"releaseId": rel.Id}).All(&workflows); err != nil {
		}
		releases[idx].Workflow = workflows
	}

	pageResults := PageResults{
		Records:    &releases,
		Pagination: p,
	}

	w.WriteJson(&pageResults)
}

func (x *Projects) currentRelease(w rest.ResponseWriter, r *rest.Request) {
	release := Release{}

	slug := r.PathParam("slug")

	projectCol := db.C("projects")
	releaseCol := db.C("releases")
	workflowCol := db.C("workflows")
	featureCol := db.C("features")
	userCol := db.C("users")

	project := Project{}
	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := releaseCol.Find(bson.M{"projectId": project.Id, "state": plugins.Complete}).Sort("-$natural").One(&release); err != nil {
		w.WriteJson(&release)
		return
	}

	headFeature := Feature{}
	tailFeature := Feature{}
	user := User{}

	if err := featureCol.Find(bson.M{"_id": release.HeadFeatureId}).One(&headFeature); err != nil {

	}

	if err := featureCol.Find(bson.M{"_id": release.TailFeatureId}).One(&tailFeature); err != nil {

	}

	if err := userCol.Find(bson.M{"_id": release.UserId}).One(&user); err != nil {

	}

	release.HeadFeature = headFeature
	release.TailFeature = tailFeature
	release.User = user

	workflows := []Flow{}
	if err := workflowCol.Find(bson.M{"releaseId": release.Id}).All(&workflows); err != nil {
	}
	release.Workflow = workflows

	w.WriteJson(&release)
}

func (x *Projects) findHeadFeature(p *Project, f *Feature) bool {
	featureCol := db.C("features")
	releaseCol := db.C("releases")
	release := Release{}

	if err := featureCol.Find(bson.M{"projectId": p.Id}).Sort("$natural").One(&f); err != nil {
		return false
	}

	if err := releaseCol.Find(bson.M{"projectId": p.Id, "state": plugins.Complete}).Sort("-$natural").One(&release); err != nil {
		log.Println("No release found, selecting first feature")
	} else {
		if err := featureCol.Find(bson.M{"_id": release.HeadFeatureId}).One(&f); err != nil {
			log.Panic(err)
		}
	}

	return true
}

func (x *Projects) findCurrentlyDeployedFeature(p *Project, f *Feature) bool {
	featureCol := db.C("features")
	releaseCol := db.C("releases")
	release := Release{}

	if err := releaseCol.Find(bson.M{"projectId": p.Id, "state": plugins.Complete}).Sort("-$natural").One(&release); err != nil {
		return false
	} else {
		if err := featureCol.Find(bson.M{"_id": release.HeadFeatureId}).One(&f); err != nil {
			log.Panic(err)
		}
	}

	return true
}

func (x *Projects) createReleases(w rest.ResponseWriter, r *rest.Request) {
	project, _ := x.CurrentProject(r)
	user, _ := CurrentUser(r)

	releaseCol := db.C("releases")
	secretCol := db.C("secrets")

	headFeature := Feature{}
	tailFeature := Feature{}

	err := r.DecodeJsonPayload(&headFeature)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false, "type": bson.M{"$in": []plugins.Type{plugins.Env, plugins.File}}}).All(&secrets); err != nil {

	}

	release := Release{
		Id:            bson.NewObjectId(),
		ProjectId:     project.Id,
		HeadFeatureId: headFeature.Id,
		HeadFeature:   headFeature,
		UserId:        user.Id,
		User:          user,
		State:         plugins.Waiting,
		Secrets:       secrets,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if ok := x.findCurrentlyDeployedFeature(&project, &tailFeature); !ok {
		release.TailFeatureId = headFeature.Id
		release.TailFeature = headFeature
	} else {
		release.TailFeatureId = tailFeature.Id
		release.TailFeature = tailFeature
	}

	if err := releaseCol.Insert(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.CreateRelease(&release)

	w.WriteJson(release)
}

func (x *Projects) createReleasesRollback(w rest.ResponseWriter, r *rest.Request) {
	project, _ := x.CurrentProject(r)
	user, _ := CurrentUser(r)

	releaseCol := db.C("releases")
	featureCol := db.C("features")
	headFeature := Feature{}
	tailFeature := Feature{}
	release := Release{}

	err := r.DecodeJsonPayload(&release)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := releaseCol.Find(bson.M{"_id": release.Id, "projectId": project.Id}).One(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := featureCol.Find(bson.M{"_id": release.HeadFeatureId}).One(&headFeature); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	if err := featureCol.Find(bson.M{"_id": release.TailFeatureId}).One(&tailFeature); err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	release.Id = bson.NewObjectId()
	release.HeadFeature = headFeature
	release.TailFeature = tailFeature
	release.User = user
	release.UserId = user.Id
	release.State = plugins.Waiting
	release.CreatedAt = time.Now()
	release.UpdatedAt = time.Now()

	if err := releaseCol.Insert(&release); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cf.CreateRelease(&release)

	w.WriteJson(release)
}

func (x *Projects) settings(w rest.ResponseWriter, r *rest.Request) {
	slug := r.PathParam("slug")
	project := Project{}
	secrets := []Secret{}

	secretCol := db.C("secrets")
	projectCol := db.C("projects")

	if err := projectCol.Find(bson.M{"slug": slug}).One(&project); err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false}).All(&secrets); err != nil {

	}

	settings := ProjectSettings{
		ProjectId: project.Id,
		GitSshUrl: project.GitSshUrl,
		Secrets:   secrets,
	}

	w.WriteJson(settings)
}

func (x *Projects) updateSettings(w rest.ResponseWriter, r *rest.Request) {
	projectSlug := r.PathParam("slug")
	project := Project{}
	projectSettingsRequest := ProjectSettings{}

	err := r.DecodeJsonPayload(&projectSettingsRequest)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	projectCol := db.C("projects")
	secretCol := db.C("secrets")

	if err := projectCol.Find(bson.M{"slug": projectSlug}).One(&project); err != nil {
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

		if err := projectCol.Update(bson.M{"_id": project.Id}, bson.M{"$set": bson.M{
			"name":       repository,
			"repository": repository,
			"slug":       slug.Slug(repository),
			"gitSshUrl":  projectSettingsRequest.GitSshUrl,
		}}); err != nil {
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

			if err := secretCol.Update(bson.M{"_id": secret.Id}, &secret); err != nil {
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
			if err := secretCol.Find(bson.M{"_id": secret.Id}).One(&s); err != nil {
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
			s.DeletedAt = time.Now()

			if err := secretCol.Update(bson.M{"_id": secret.Id}, &s); err != nil {
				rest.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}

		secret.Id = bson.NewObjectId()
		secret.ProjectId = project.Id
		secret.CreatedAt = time.Now()

		if err := secretCol.Insert(&secret); err != nil {
			rest.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false}).All(&secrets); err != nil {

	}

	projectSettingsResponse := ProjectSettings{
		ProjectId: projectSettingsRequest.ProjectId,
		GitSshUrl: projectSettingsRequest.GitSshUrl,
		Secrets:   secrets,
		UpdatedAt: time.Now(),
	}

	w.WriteJson(projectSettingsResponse)
}
