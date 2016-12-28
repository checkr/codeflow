package codeflow

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/checkr/codeflow/server/plugins/codeflow/db"
	"github.com/spf13/viper"
	"gopkg.in/go-playground/validator.v9"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbSession *mgo.Session
var db *mgo.Database
var validate *validator.Validate
var cf *Codeflow

type CodeflowHandler interface {
	Register(api *rest.Api) []*rest.Route
}

func init() {
	agent.RegisterPlugin("codeflow", func() agent.Plugin { return NewCodeflow() })
}

type Codeflow struct {
	ServiceAddress string `mapstructure:"service_address"`
	events         chan agent.Event

	Projects  *Projects
	Auth      *Auth
	Users     *Users
	Bookmarks *Bookmarks
}

func NewCodeflow() *Codeflow {
	return &Codeflow{}
}

func (x *Codeflow) SampleConfig() string {
	return `
  ## Address and port to host Codeflow listener on
	service_address: ":3000"
	builds:
		path: "/builds"
	`
}

func (x *Codeflow) Description() string {
	return "A Codeflow Event collector"
}

func (x *Codeflow) Listen() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	api.Use(&rest.CorsMiddleware{
		RejectNonCorsRequests: false,
		OriginValidator: func(origin string, request *rest.Request) bool {
			allowedOrigins := viper.GetStringSlice("allowed_origins")
			if agent.SliceContains(origin, allowedOrigins) {
				return true
			}
			return false
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{
			"Authorization", "Accept", "Content-Type", "X-Custom-Header", "Origin"},
		AccessControlAllowCredentials: true,
		AccessControlMaxAge:           3600,
	})

	var routes []*rest.Route
	var handlerRoutes []*rest.Route

	for _, handler := range x.AvailableCodeflowHandlers() {
		handlerRoutes = handler.Register(api)
		routes = append(routes, handlerRoutes...)
	}
	router, err := rest.MakeRouter(routes...)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s", x.ServiceAddress), api.MakeHandler()))
}

// Looks for fields which implement CodeflowHandler interface
func (x *Codeflow) AvailableCodeflowHandlers() []CodeflowHandler {
	handlers := make([]CodeflowHandler, 0)
	s := reflect.ValueOf(x).Elem()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		if !f.CanInterface() {
			continue
		}

		if handler, ok := f.Interface().(CodeflowHandler); ok {
			if !reflect.ValueOf(handler).IsNil() {
				handlers = append(handlers, handler)
			}
		}
	}

	return handlers
}

func (x *Codeflow) Start(events chan agent.Event) error {
	var err error

	x.events = events
	validate = validator.New()
	cf = x

	mongoConfig := codeflow_db.Config{
		URI: viper.GetString("mongo.uri"),
		SSL: viper.GetBool("mongo.ssl"),
		Creds: mgo.Credential{
			Username: viper.GetString("mongo.username"),
			Password: viper.GetString("mongo.password"),
		},
	}

	dbSession, err = codeflow_db.NewConnection(mongoConfig)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %s", err))
	}

	db = dbSession.DB(viper.GetString("mongo.database"))

	go x.Listen()
	log.Printf("Started Codeflow service on %s\n", x.ServiceAddress)
	return nil
}

func (x *Codeflow) Stop() {
	log.Println("Stopping Codeflow service")
}

func (x *Codeflow) Subscribe(e agent.Event) []string {
	return []string{
		"plugins.GitPing",
		"plugins.GitCommit",
		"plugins.DockerBuild:status",
		"plugins.HeartBeat",
		"plugins.LoadBalancer:status",
		"plugins.DockerDeploy:status",
	}
}

func (x *Codeflow) Process(e agent.Event) error {
	log.Printf("Process Codeflow event: %s", e.Name)

	if e.Name == "plugins.DockerDeploy:status" {
		dockerDeploy := e.Payload.(plugins.DockerDeploy)
		x.DockerDeployStatus(&dockerDeploy)
	}

	if e.Name == "plugins.LoadBalancer:status" {
		payload := e.Payload.(plugins.LoadBalancer)
		x.LoadBalancerStatus(&payload)
	}

	if e.Name == "plugins.GitPing" {
		payload := e.Payload.(plugins.GitPing)
		project := Project{}

		collection := db.C("projects")
		if err := collection.Find(bson.M{"repository": payload.Repository}).One(&project); err != nil {
			return err
		}

		project.Pinged = true
		if err := collection.Update(bson.M{"_id": project.Id}, project); err != nil {
			return err
		}

		return nil
	}

	if e.Name == "plugins.GitCommit" {
		payload := e.Payload.(plugins.GitCommit)

		// Only build master branch
		if payload.Ref != "refs/heads/master" {
			return nil
		}

		project := Project{}
		feature := Feature{}

		projectCol := db.C("projects")
		featureCol := db.C("features")

		if err := projectCol.Find(bson.M{"repository": payload.Repository}).One(&project); err != nil {
			return err
		}

		if !project.Pinged {
			project.Pinged = true
			if err := projectCol.Update(bson.M{"_id": project.Id}, project); err != nil {
				return err
			}
		}

		if err := featureCol.Find(bson.M{"hash": payload.Hash}).One(&feature); err != nil {
			switch err {
			case mgo.ErrNotFound:
				feature = Feature{
					ProjectId:  project.Id,
					Message:    payload.Message,
					User:       payload.User,
					Hash:       payload.Hash,
					ParentHash: payload.ParentHash,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}

				if err := featureCol.Insert(feature); err != nil {
					return err
				}

				x.CreateFeature(&feature, e)
				return nil
			default:
				return err
			}
		}

		if err := featureCol.Update(bson.M{"_id": feature.Id}, feature); err != nil {
			return err
		}

		x.CreateFeature(&feature, e)
	}

	if e.Name == "plugins.DockerBuild:status" {
		payload := e.Payload.(plugins.DockerBuild)
		x.DockerBuildStatus(&payload)
	}

	return nil
}

func CurrentUser(r *rest.Request) (User, error) {
	user := User{}

	userCol := db.C("users")
	if err := userCol.Find(bson.M{"username": r.Env["REMOTE_USER"]}).One(&user); err != nil {
	}

	return user, nil
}

func (x *Codeflow) ProjectChange(p *Project, name string, msg string) {
	projectChangeCol := db.C("projectChanges")
	projectChange := ProjectChange{
		Id:        bson.NewObjectId(),
		ProjectId: p.Id,
		Name:      name,
		Message:   msg,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := projectChangeCol.Insert(&projectChange); err != nil {
		log.Println(err)
	}
}

func (x *Codeflow) CreateProject(p *Project) {
	wsMsg := plugins.WebsocketMsg{
		Channel: "projects/new",
		Payload: p,
	}
	event := agent.NewEvent(wsMsg, nil)
	x.events <- event
}

func (x *Codeflow) CreateService(s *Service) {
	project := Project{}
	projectCol := db.C("projects")

	if err := projectCol.Find(bson.M{"_id": s.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	var services []*Service
	services = append(services, s)

	secretCol := db.C("secrets")
	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false, "type": bson.M{"$in": []plugins.Type{plugins.Env, plugins.File}}}).All(&secrets); err != nil {

	}
}

func (x *Codeflow) UpdateService(s *Service) {
	project := Project{}
	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"_id": s.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	var services []*Service
	services = append(services, s)

	secretCol := db.C("secrets")
	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false, "type": bson.M{"$in": []plugins.Type{plugins.Env, plugins.File}}}).All(&secrets); err != nil {

	}

	x.ProjectChange(&project, "serviceUpdate", fmt.Sprintf("Service `%v` was updated", s.Name))
}

func (x *Codeflow) DeleteService(s *Service) {
	project := Project{}
	projectCol := db.C("projects")
	if err := projectCol.Find(bson.M{"_id": s.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	var services []*Service
	services = append(services, s)

	secretCol := db.C("secrets")
	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false, "type": bson.M{"$in": []plugins.Type{plugins.Env, plugins.File}}}).All(&secrets); err != nil {

	}
}

func (x *Codeflow) CreateRelease(r *Release) {
	project := Project{}

	projectCol := db.C("projects")
	workflowCol := db.C("workflows")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"_id": r.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	var services []*Service

	if err := serviceCol.Find(bson.M{"projectId": r.ProjectId}).All(&services); err != nil {
		panic(err)
	}

	// Create Workflow
	// Pull this info from project spec/req
	workflows := []string{"build:DockerImage"}

	for _, str := range workflows {
		s := strings.Split(str, ":")
		t, n := s[0], s[1]
		flow := Flow{
			Id:        bson.NewObjectId(),
			ReleaseId: r.Id,
			Type:      t,
			Name:      n,
			State:     plugins.Waiting,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := workflowCol.Insert(&flow); err != nil {
			fmt.Println(err)
		}

		r.Workflow = append(r.Workflow, flow)
	}

	x.CheckWorkflows(r)
}

func (x *Codeflow) CheckWorkflows(r *Release) {
	var workflowStatus plugins.State = plugins.Complete
	projectCol := db.C("projects")
	featureCol := db.C("features")
	releaseCol := db.C("releases")
	workflowCol := db.C("workflows")
	userCol := db.C("users")

	project := Project{}
	headFeature := Feature{}
	tailFeature := Feature{}
	user := User{}
	workflows := []Flow{}

	if err := projectCol.Find(bson.M{"_id": r.ProjectId}).One(&project); err != nil {
		fmt.Println("Project not found!")
		return
	}

	if err := featureCol.Find(bson.M{"_id": r.HeadFeatureId}).One(&headFeature); err != nil {
		fmt.Println("HeadFeature not found!")
		return
	}

	if err := featureCol.Find(bson.M{"_id": r.TailFeatureId}).One(&tailFeature); err != nil {
		fmt.Println("TailFeature not found!")
		return
	}

	if err := userCol.Find(bson.M{"_id": r.UserId}).One(&user); err != nil {
		fmt.Println("User not found!")
		return
	}

	workflowCol.Find(bson.M{"releaseId": r.Id}).All(&workflows)

	r.HeadFeature = headFeature
	r.TailFeature = tailFeature
	r.User = user
	r.Workflow = workflows

	for idx, _ := range r.Workflow {
		flow := &r.Workflow[idx]

		switch flow.Type {
		case "build":
			buildCol := db.C("builds")
			build := Build{}
			if err := buildCol.Find(bson.M{"featureHash": r.HeadFeature.Hash, "type": flow.Name}).One(&build); err != nil {
				flow.State = plugins.Failed
			} else {
				switch build.State {
				case "waiting":
					flow.State = plugins.Waiting
				case "complete":
					flow.State = plugins.Complete
				case "failed":
					flow.State = plugins.Failed
					r.State = plugins.Failed
				default:
					flow.State = plugins.Running
				}
			}

			flow.UpdatedAt = time.Now()

			if err := workflowCol.Update(bson.M{"_id": flow.Id}, &flow); err != nil {
				panic(err)
			}

			// Update release state
			releaseCol.Update(bson.M{"_id": r.Id}, bson.M{"$set": bson.M{"state": r.State}})

			if flow.State != plugins.Complete {
				workflowStatus = flow.State
			}

			break
		}
	}

	if workflowStatus == plugins.Complete {
		x.CreateDeploy(r)
	}
}

func (x *Codeflow) CreateExtensions(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Create)
}

func (x *Codeflow) UpdateExtensions(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Update)
}

func (x *Codeflow) DeleteExtensions(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Destroy)
}

func (x *Codeflow) loadBalancer(lb *LoadBalancer, action plugins.Action) {
	project := Project{}
	service := Service{}

	projectCol := db.C("projects")
	serviceCol := db.C("services")

	if err := projectCol.Find(bson.M{"_id": lb.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	if err := serviceCol.Find(bson.M{"_id": lb.ServiceId}).One(&service); err != nil {
		fmt.Println("Load balancer is not attached to any service or sercice was already deleted")
	}

	var serviceListeners []plugins.Listener

	for _, listener := range service.Listeners {
		serviceListeners = append(serviceListeners, plugins.Listener{
			Port:     int32(listener.Port),
			Protocol: listener.Protocol,
		})
	}

	var listenerPairs []plugins.ListenerPair

	for _, listenerPair := range lb.ListenerPairs {
		listenerPairs = append(listenerPairs, plugins.ListenerPair{
			Source: plugins.Listener{
				Port:     int32(listenerPair.Source.Port),
				Protocol: listenerPair.Destination.Protocol,
			},
			Destination: plugins.Listener{
				Port:     int32(listenerPair.Destination.Port),
				Protocol: listenerPair.Destination.Protocol,
			},
		})
	}

	loadBalancerEvent := plugins.LoadBalancer{
		Name:   lb.Name,
		Action: action,
		Type:   StringToLoadBalancerType(lb.Type),
		Project: plugins.Project{
			Slug:       project.Slug,
			Repository: project.Repository,
		},
		Service: plugins.Service{
			Name:      service.Name,
			Command:   service.Command,
			Listeners: serviceListeners,
			Replicas:  int64(service.Count),
		},
		ListenerPairs: listenerPairs,
		Environment:   "development",
	}

	event := agent.NewEvent(loadBalancerEvent, nil)
	x.events <- event
}

func (x *Codeflow) CreateFeature(feature *Feature, e agent.Event) {
	project := Project{}

	projectCol := db.C("projects")
	secretCol := db.C("secrets")
	buildCol := db.C("builds")

	if err := projectCol.Find(bson.M{"_id": feature.ProjectId}).One(&project); err != nil {
		panic(err)
	}

	build := Build{
		FeatureHash: feature.Hash,
		Type:        "DockerImage",
		State:       plugins.Waiting,
		CreatedAt:   time.Now(),
	}

	_, err := buildCol.Upsert(bson.M{"featureHash": feature.Hash}, &build)
	if err != nil {
		fmt.Println(err)
	}

	secrets := []Secret{}
	if err := secretCol.Find(bson.M{"projectId": project.Id, "deleted": false, "type": plugins.Build}).All(&secrets); err != nil {

	}

	var buildArgs []plugins.Arg
	for _, secret := range secrets {
		arg := plugins.Arg{
			Key:   secret.Key,
			Value: secret.Value,
		}
		buildArgs = append(buildArgs, arg)
	}

	dockerBuildEvent := plugins.DockerBuild{
		Action: plugins.Create,
		State:  plugins.Waiting,
		Project: plugins.Project{
			Slug:       project.Slug,
			Repository: project.Repository,
		},
		Git: plugins.Git{
			SshUrl:        project.GitSshUrl,
			RsaPrivateKey: project.RsaPrivateKey,
			RsaPublicKey:  project.RsaPublicKey,
		},
		Feature: plugins.Feature{
			Hash:       feature.Hash,
			ParentHash: feature.ParentHash,
			User:       feature.User,
			Message:    feature.Message,
		},
		Registry: plugins.DockerRegistry{
			Host:     viper.GetString("plugins.docker_build.registry_host"),
			Username: viper.GetString("plugins.docker_build.registry_username"),
			Password: viper.GetString("plugins.docker_build.registry_password"),
			Email:    viper.GetString("plugins.docker_build.registry_user_email"),
		},
		BuildArgs: buildArgs,
	}

	x.events <- e.NewEvent(dockerBuildEvent, nil)

	wsMsg := plugins.WebsocketMsg{
		Channel: fmt.Sprintf("create/projects/%v/feature", project.Slug),
		Payload: feature,
	}

	x.events <- agent.NewEvent(wsMsg, nil)
}

func (x *Codeflow) DockerBuildStatus(br *plugins.DockerBuild) {
	project := Project{}
	build := Build{}

	projectCol := db.C("projects")
	buildCol := db.C("builds")

	if err := projectCol.Find(bson.M{"repository": br.Project.Repository}).One(&project); err != nil {
		panic(err)
	}

	buildCol.Update(bson.M{"featureHash": br.Feature.Hash}, bson.M{"$set": bson.M{
		"image":      br.Image,
		"state":      br.State,
		"buildLog":   br.BuildLog,
		"buildError": br.BuildError,
		"updatedAt":  time.Now(),
	}})

	if err := buildCol.Find(bson.M{"featureHash": br.Feature.Hash}).One(&build); err != nil {
		panic(err)
	}

	// Check if there is any releases pending and call their workflow
	featureCol := db.C("features")
	releaseCol := db.C("releases")
	feature := Feature{}
	releases := []Release{}

	if err := featureCol.Find(bson.M{"hash": br.Feature.Hash}).One(&feature); err != nil {
		fmt.Println("Feature not found!")
		return
	}

	if err := releaseCol.Find(bson.M{"headFeatureId": feature.Id, "state": plugins.Waiting}).All(&releases); err != nil {
		fmt.Println("No pending releases found.")
		return
	}

	for _, rel := range releases {
		x.CheckWorkflows(&rel)
		x.ReleaseUpdated(&rel)
	}
}

func (x *Codeflow) CreateDeploy(r *Release) {
	releaseServices := []Service{}
	build := Build{}
	project := Project{}

	serviceCol := db.C("services")
	buildCol := db.C("builds")
	projectCol := db.C("projects")

	if err := projectCol.Find(bson.M{"_id": r.ProjectId}).One(&project); err != nil {
		fmt.Println("Project not found!")
		return
	}

	if err := serviceCol.Find(bson.M{"projectId": r.ProjectId, "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running, plugins.Deleting}}}).All(&releaseServices); err != nil {
		fmt.Println("Service not found! Can't deploy without service")
		return
	}

	if err := buildCol.Find(bson.M{"featureHash": r.HeadFeature.Hash, "type": "DockerImage", "state": plugins.Complete}).One(&build); err != nil {
		fmt.Println("Build not found!")
		return
	}

	headFeature := plugins.Feature{
		Hash:       r.HeadFeature.Hash,
		ParentHash: r.HeadFeature.ParentHash,
		User:       r.HeadFeature.User,
		Message:    r.HeadFeature.Message,
	}

	tailFeature := plugins.Feature{
		Hash:       r.TailFeature.Hash,
		ParentHash: r.TailFeature.ParentHash,
		User:       r.TailFeature.User,
		Message:    r.TailFeature.Message,
	}

	var services []plugins.Service

	for _, service := range releaseServices {
		var listeners []plugins.Listener

		for _, listener := range service.Listeners {
			listeners = append(listeners, plugins.Listener{
				Port:     int32(listener.Port),
				Protocol: listener.Protocol,
			})
		}

		var action plugins.Action
		switch service.State {
		case plugins.Waiting:
			action = plugins.Create
		case plugins.Deleting:
			action = plugins.Destroy
		default:
			action = plugins.Update
		}

		services = append(services, plugins.Service{
			Action:    action,
			Name:      service.Name,
			Command:   service.Command,
			Listeners: listeners,
			Replicas:  int64(service.Count),
		})
	}

	var secrets []plugins.Secret

	for _, secret := range r.Secrets {
		newS := plugins.Secret{
			Key:   secret.Key,
			Value: secret.Value,
			Type:  secret.Type,
		}
		secrets = append(secrets, newS)
	}

	dockerDeployEvent := plugins.DockerDeploy{
		Action: plugins.Create,
		Project: plugins.Project{
			Slug:       project.Slug,
			Repository: project.Repository,
		},
		Release: plugins.Release{
			Id:          r.Id.Hex(),
			HeadFeature: headFeature,
			TailFeature: tailFeature,
		},
		Docker: plugins.Docker{
			Image: build.Image,
			Registry: plugins.DockerRegistry{
				Host:     "",
				Username: "",
				Password: "",
				Email:    "",
			},
		},
		Services:    services,
		Secrets:     secrets,
		Environment: "development",
	}

	event := agent.NewEvent(dockerDeployEvent, nil)
	x.events <- event
}

func (x *Codeflow) LoadBalancerStatus(lb *plugins.LoadBalancer) {
	projectCol := db.C("projects")
	extensionCol := db.C("extensions")
	var extension LoadBalancer
	var project Project

	if err := projectCol.Find(bson.M{"repository": lb.Project.Repository}).One(&project); err != nil {
		fmt.Println("Project not found!")
		return
	}

	if err := extensionCol.Find(bson.M{"projectId": project.Id, "name": lb.Name}).One(&extension); err != nil {
		fmt.Printf("LoadBalancer %s not found!\n", lb.Name)
		return
	}

	extensionCol.Update(bson.M{"_id": extension.Id}, bson.M{"$set": bson.M{"dnsName": lb.DNSName, "state": lb.State, "stateMessage": lb.StateMessage}})
}

func (x *Codeflow) DockerDeployStatus(e *plugins.DockerDeploy) {
	projectCol := db.C("projects")
	featureCol := db.C("features")
	releaseCol := db.C("releases")

	feature := Feature{}
	release := Release{}
	project := Project{}

	if err := projectCol.Find(bson.M{"repository": e.Project.Repository}).One(&project); err != nil {
		fmt.Println("Project not found!")
		return
	}

	if err := featureCol.Find(bson.M{"hash": e.Release.HeadFeature.Hash}).One(&feature); err != nil {
		fmt.Println("Feature not found!")
		return
	}

	if err := releaseCol.Find(bson.M{"_id": bson.ObjectIdHex(e.Release.Id)}).One(&release); err != nil {
		fmt.Println("Release not found!")
		return
	}

	// Update release state
	releaseCol.Update(bson.M{"_id": release.Id}, bson.M{"$set": bson.M{"state": e.State}})
	release.State = e.State

	if release.State == plugins.Complete {

	}

	x.ReleaseUpdated(&release)
}

func (x *Codeflow) ReleaseUpdated(r *Release) {
	projectCol := db.C("projects")
	featureCol := db.C("features")
	workflowCol := db.C("workflows")
	userCol := db.C("users")

	project := Project{}
	headFeature := Feature{}
	tailFeature := Feature{}
	user := User{}
	workflows := []Flow{}

	if err := projectCol.Find(bson.M{"_id": r.ProjectId}).One(&project); err != nil {
		fmt.Println("Project not found!")
		return
	}

	if err := featureCol.Find(bson.M{"_id": r.HeadFeatureId}).One(&headFeature); err != nil {
		fmt.Println("HeadFeature not found!")
		return
	}

	if err := featureCol.Find(bson.M{"_id": r.TailFeatureId}).One(&tailFeature); err != nil {
		fmt.Println("TailFeature not found!")
		return
	}

	if err := userCol.Find(bson.M{"_id": r.UserId}).One(&user); err != nil {
		fmt.Println("User not found!")
		return
	}

	r.HeadFeature = headFeature
	r.TailFeature = tailFeature
	r.User = user

	if err := workflowCol.Find(bson.M{"releaseId": r.Id}).All(&workflows); err != nil {

	}

	r.Workflow = workflows

	wsMsg := plugins.WebsocketMsg{
		Channel: fmt.Sprintf("update/projects/%v/release", project.Slug),
		Payload: r,
	}

	x.events <- agent.NewEvent(wsMsg, nil)
}

func (x *Codeflow) UpdateBookmarks(u *User) {
	bookmarks := []Bookmark{}
	bookmarkCol := db.C("bookmarks")

	if err := bookmarkCol.Find(bson.M{"userId": u.Id}).All(&bookmarks); err != nil {
		fmt.Println("Bookmarks error: " + err.Error())
	}

	projectCol := db.C("projects")
	for idx, bookmark := range bookmarks {
		project := Project{}
		if err := projectCol.Find(bson.M{"_id": bookmark.ProjectId}).One(&project); err != nil {

		}
		bookmarks[idx].Name = project.Name
		bookmarks[idx].Slug = project.Slug
	}

	wsMsg := plugins.WebsocketMsg{
		Channel: "bookmarks/" + u.Id.Hex(),
		Payload: bookmarks,
	}
	event := agent.NewEvent(wsMsg, nil)
	x.events <- event
}

// Calc produces a Paging calculated over
// Page, Limit, Count and DefaultLimit values of given Paging
func (p Pagination) Calc() *Pagination {
	if p.Page < 1 {
		p.Page = 1
	}

	if p.Limit < 1 {
		if p.DefaultLimit > 0 {
			p.Limit = p.DefaultLimit
		} else {
			p.Limit = 10
		}
	}

	if p.Count < p.Limit {
		p.TotalPages = 1
	} else {
		p.TotalPages = int(math.Ceil(float64(p.Count) / float64(p.Limit)))
	}

	if p.Page > p.TotalPages {
		p.Page = p.TotalPages
	}

	p.Offset = (p.Page - 1) * p.Limit

	return &p
}

func StringToLoadBalancerType(s string) plugins.Type {
	switch strings.ToLower(s) {
	case "internal":
		return plugins.Internal
	case "external":
		return plugins.External
	case "office":
		return plugins.Office
	}
	return plugins.Internal
}
