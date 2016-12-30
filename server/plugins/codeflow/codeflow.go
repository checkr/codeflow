package codeflow

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/davecgh/go-spew/spew"
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
	Stats     *Stats
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

	mongoConfig := MongoConfig{
		URI: viper.GetString("mongo.uri"),
		SSL: viper.GetBool("mongo.ssl"),
		Creds: mgo.Credential{
			Username: viper.GetString("mongo.username"),
			Password: viper.GetString("mongo.password"),
		},
	}

	dbSession, err = NewMongoConnection(mongoConfig)
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
		var project Project
		var err error

		if project, err = GetProjectByRepository(payload.Repository); err != nil {
			log.Println(err.Error())
			return err
		}

		project.Pinged = true
		if err = UpdateProject(project.Id, &project); err != nil {
			log.Println(err.Error())
			return err
		}

		return nil
	}

	if e.Name == "plugins.GitCommit" {
		payload := e.Payload.(plugins.GitCommit)
		var err error
		var project Project
		var feature Feature

		// Only build master branch
		if payload.Ref != "refs/heads/master" {
			return nil
		}

		if project, err = GetProjectByRepository(payload.Repository); err != nil {
			log.Println(err.Error())
			return err
		}

		if !project.Pinged {
			project.Pinged = true
			if err = UpdateProject(project.Id, &project); err != nil {
				log.Println(err.Error())
				return err
			}
		}

		if feature, err = GetFeatureByHash(payload.Hash); err != nil {
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

				if err := CreateFeature(&feature); err != nil {
					log.Println(err.Error())
					return err
				}

				x.FeatureCreated(&feature, e)
				return nil
			default:
				log.Println(err.Error())
				return err
			}
		}

		log.Printf("Feature `%v:%v` already exists", project.Repository, payload.Hash)
	}

	if e.Name == "plugins.DockerBuild:status" {
		payload := e.Payload.(plugins.DockerBuild)
		x.DockerBuildStatus(&payload)
	}

	return nil
}

func (x *Codeflow) ProjectCreated(p *Project) {
	wsMsg := plugins.WebsocketMsg{
		Channel: "projects/new",
		Payload: p,
	}
	event := agent.NewEvent(wsMsg, nil)
	x.events <- event
}

func (x *Codeflow) ServiceCreated(s *Service) {
	var err error
	var project Project

	if project, err = GetProjectById(s.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	// TODO: Magic
	spew.Dump(project, s)
}

func (x *Codeflow) ServiceUpdated(s *Service) {
	var err error
	var project Project

	if project, err = GetProjectById(s.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	// TODO: Magic
	spew.Dump(project, s)
}

func (x *Codeflow) ServiceDeleted(s *Service) {
	var err error
	var project Project

	if project, err = GetProjectById(s.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	// TODO: Magic
	spew.Dump(project, s)
}

func (x *Codeflow) ReleaseCreated(r *Release) {
	var err error

	if _, err = GetProjectById(r.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	// TODO: Read required workflows from db
	workflows := []string{"build:DockerImage"}

	for _, str := range workflows {
		s := strings.Split(str, ":")
		t, n := s[0], s[1]
		flow := Flow{
			ReleaseId: r.Id,
			Type:      t,
			Name:      n,
			State:     plugins.Waiting,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := CreateWorkflow(&flow); err != nil {
			log.Println(err.Error())
			return
		}

		r.Workflow = append(r.Workflow, flow)
	}

	x.CheckWorkflows(r)
	x.ReleaseUpdated(r)
}

func (x *Codeflow) CheckWorkflows(r *Release) {
	var err error
	var workflowStatus plugins.State = plugins.Complete

	if err = PopulateRelease(r); err != nil {
		log.Println(err.Error())
		return
	}

	for idx, _ := range r.Workflow {
		flow := &r.Workflow[idx]

		switch flow.Type {
		case "build":
			var build Build
			if build, err = GetBuildByHashAndType(r.HeadFeature.Hash, flow.Name); err != nil {
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

			if err := UpdateFlow(flow.Id, flow); err != nil {
				log.Println(err.Error())
				return
			}

			if err := UpdateRelease(r.Id, r); err != nil {
				log.Println(err.Error())
				return
			}

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

func (x *Codeflow) ExtensionCreated(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Create)
}

func (x *Codeflow) ExtensionUpdated(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Update)
}

func (x *Codeflow) ExtensionDeleted(lb *LoadBalancer) {
	x.loadBalancer(lb, plugins.Destroy)
}

func (x *Codeflow) loadBalancer(lb *LoadBalancer, action plugins.Action) {
	var err error
	var project Project
	var service Service

	if _, err = GetProjectById(lb.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	if service, err = GetServiceById(lb.ServiceId); err != nil {
		log.Println("Load balancer is not attached to any service or service was deleted")
		log.Println(err.Error())
		return
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
	}

	event := agent.NewEvent(loadBalancerEvent, nil)
	x.events <- event
}

func (x *Codeflow) FeatureCreated(feature *Feature, e agent.Event) {
	var err error
	var project Project
	var secrets []Secret

	if project, err = GetProjectById(feature.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	build := Build{
		FeatureHash: feature.Hash,
		Type:        "DockerImage",
		State:       plugins.Waiting,
	}

	if err = CreateOrUpdateBuildByFeatureHash(feature.Hash, &build); err != nil {
		log.Println(err.Error())
		return
	}

	if secrets, err = GetSecretsByProjectIdAndType(project.Id, plugins.Build); err != nil {
		log.Println(err.Error())
		return
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
	var err error
	var build Build
	var feature Feature

	if _, err = GetProjectByRepository(br.Project.Repository); err != nil {
		log.Println(err.Error())
		return
	}

	if build, err = GetBuildByHash(br.Feature.Hash); err != nil {
		log.Println(err.Error())
		return
	}

	build.Image = br.Image
	build.State = br.State
	build.BuildLog = br.BuildLog
	build.BuildError = br.BuildError

	if err := UpdateBuild(build.Id, &build); err != nil {
		log.Println(err.Error())
		return
	}

	if feature, err = GetFeatureByHash(br.Feature.Hash); err != nil {
		log.Println(err.Error())
		return
	}

	x.UpdateInProgessReleases(&feature)
}

func (x *Codeflow) UpdateInProgessReleases(f *Feature) {
	var err error
	var releases []Release

	if releases, err = GetReleasesByFeatureIdAndState(f.Id, plugins.Waiting); err != nil {
		log.Println(err.Error())
		return
	}

	for _, rel := range releases {
		x.CheckWorkflows(&rel)
		x.ReleaseUpdated(&rel)
	}
}

func (x *Codeflow) CreateDeploy(r *Release) {
	var err error
	var project Project
	var build Build
	var releaseServices []Service

	if project, err = GetProjectById(r.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	if releaseServices, err = GetReleaseServices(r.ProjectId); err != nil {
		log.Println(err.Error())
		fmt.Println("Service not found! Can't deploy without service")
		return
	}

	if build, err = GetBuildByHashAndState(r.HeadFeature.Hash, plugins.Complete); err != nil {
		log.Println(err.Error())
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
	var err error
	var project Project
	var extension LoadBalancer

	if project, err = GetProjectByRepository(lb.Project.Repository); err != nil {
		log.Println(err.Error())
		return
	}

	if extension, err = GetExtensionByProjectIdAndName(project.Id, lb.Name); err != nil {
		fmt.Printf("LoadBalancer %s not found!\n", lb.Name)
		log.Println(err.Error())
		return
	}

	extension.DNSName = lb.DNSName
	extension.State = lb.State
	extension.StateMessage = lb.StateMessage

	if err = UpdateExtension(extension.Id, &extension); err != nil {
		log.Println(err.Error())
		return
	}
}

func (x *Codeflow) DockerDeployStatus(e *plugins.DockerDeploy) {
	var err error
	var release Release

	if release, err = GetReleaseById(bson.ObjectIdHex(e.Release.Id)); err != nil {
		log.Println(err.Error())
		return
	}

	release.State = e.State

	if err = UpdateRelease(release.Id, &release); err != nil {
		log.Println(err.Error())
		return
	}

	if release.State == plugins.Complete {

	}

	x.ReleaseUpdated(&release)
}

func (x *Codeflow) ReleaseUpdated(r *Release) {
	var err error
	var project Project

	if project, err = GetProjectById(r.ProjectId); err != nil {
		log.Println(err.Error())
		return
	}

	if err = PopulateRelease(r); err != nil {
		log.Println(err.Error())
		return
	}

	wsMsg := plugins.WebsocketMsg{
		Channel: fmt.Sprintf("update/projects/%v/release", project.Slug),
		Payload: r,
	}

	x.events <- agent.NewEvent(wsMsg, nil)
}

func (x *Codeflow) BookmarksUpdated(u *User) {
	var err error
	var bookmarks []Bookmark

	if bookmarks, err = GetUserBookmarks(u.Id); err != nil {
		log.Println(err.Error())
		return
	}

	wsMsg := plugins.WebsocketMsg{
		Channel: "bookmarks/" + u.Id.Hex(),
		Payload: bookmarks,
	}

	event := agent.NewEvent(wsMsg, nil)
	x.events <- event
}
