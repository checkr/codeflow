package codeflow

import (
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

var cf *Codeflow
var db *bongo.Connection

type CodeflowHandler interface {
	Register(api *rest.Api) []*rest.Route
}

func init() {
	agent.RegisterPlugin("codeflow", func() agent.Plugin { return NewCodeflow() })
}

type Codeflow struct {
	ServiceAddress string `mapstructure:"service_address"`
	Events         chan agent.Event

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
			allowedOrigins := viper.GetStringSlice("plugins.codeflow.allowed_origins")
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

	x.Events = events
	cf = x

	config := &bongo.Config{
		ConnectionString: viper.GetString("plugins.codeflow.mongodb.uri"),
		Database:         viper.GetString("plugins.codeflow.mongodb.database"),
	}

	db, err = bongo.Connect(config)
	if err != nil {
		log.Fatal(err)
	}

	go x.Listen()

	log.Printf("Started Codeflow service on %s\n", x.ServiceAddress)

	return nil
}

func (x *Codeflow) Stop() {
	log.Println("Stopping Codeflow service")
}

func (x *Codeflow) Subscribe() []string {
	return []string{
		"plugins.GitPing",
		"plugins.GitCommit",
		"plugins.GitStatus",
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
		DockerDeployStatus(&dockerDeploy)
	}

	if e.Name == "plugins.LoadBalancer:status" {
		payload := e.Payload.(plugins.LoadBalancer)
		LoadBalancerStatus(&payload)
	}

	if e.Name == "plugins.GitPing" {
		payload := e.Payload.(plugins.GitPing)
		project := Project{}

		if err := db.Collection("projects").FindOne(bson.M{"repository": payload.Repository}, &project); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				log.Printf("Projects::FindOne::DocumentNotFoundError: repository: `%v`", payload.Repository)
			} else {
				log.Printf("Projects::FindOne::Error: %s", err.Error())
			}
			return err
		}

		project.Pinged = true

		if err := db.Collection("projects").Save(&project); err != nil {
			log.Printf("Projects::Save::Error: %v", err.Error())
			return err
		}

		return nil
	}

	if e.Name == "plugins.GitStatus" {
		payload := e.Payload.(plugins.GitStatus)
		project := Project{}
		feature := Feature{}

		if err := db.Collection("projects").FindOne(bson.M{"repository": payload.Repository}, &project); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				log.Printf("Projects::FindOne::DocumentNotFoundError: repository: `%v`", payload.Repository)
			} else {
				log.Printf("Projects::FindOne::Error: %s", err.Error())
			}
			return err
		}

		externalFlowStatus := ExternalFlowStatus{
			ProjectId:     project.Id,
			Hash:          payload.Hash,
			Context:       payload.Context,
			Message:       "",
			State:         StringToState(payload.State),
			OriginalState: payload.State,
		}

		if err := db.Collection("externalFlowStatuses").Save(&externalFlowStatus); err != nil {
			log.Printf("ExternalFlowStatuses::Save::Error: %v", err.Error())
			return err
		}

		if err := db.Collection("features").FindOne(bson.M{"hash": payload.Hash}, &feature); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				log.Printf("Features::FindOne::DocumentNotFoundError: hash: `%v`", payload.Hash)
			} else {
				log.Printf("Features::FindOne::Error: %s", err.Error())
			}
			return err
		}

		UpdateInProgessReleases(&feature)

		return nil
	}

	if e.Name == "plugins.GitCommit" {
		payload := e.Payload.(plugins.GitCommit)
		project := Project{}
		feature := Feature{}

		// Only build master branch
		if payload.Ref != "refs/heads/master" {
			return nil
		}

		if err := db.Collection("projects").FindOne(bson.M{"repository": payload.Repository}, &project); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				log.Printf("Projects::FindOne::DocumentNotFoundError: repository: `%v`", payload.Repository)
			} else {
				log.Printf("Projects::FindOne::Error: %s", err.Error())
			}
			return err
		}

		if !project.Pinged {
			project.Pinged = true

			if err := db.Collection("projects").Save(&project); err != nil {
				log.Printf("Projects::Save::Error: %v", err.Error())
				return err
			}
		}

		if err := db.Collection("features").FindOne(bson.M{"hash": payload.Hash}, &feature); err != nil {
			if _, ok := err.(*bongo.DocumentNotFoundError); ok {
				feature = Feature{
					ProjectId:  project.Id,
					Message:    payload.Message,
					User:       payload.User,
					Hash:       payload.Hash,
					ParentHash: payload.ParentHash,
				}
				if err := db.Collection("features").Save(&feature); err != nil {
					log.Printf("Features::Save::Error: %v", err.Error())
					return err
				}

				FeatureCreated(&feature, e)

				if project.ContinuousDelivery {
					user := User{}
					release := Release{}

					if err := db.Collection("users").FindOne(bson.M{"email": "codeflow"}, &user); err != nil {
						if _, ok := err.(*bongo.DocumentNotFoundError); ok {
							log.Printf("Users::FindOne::DocumentNotFoundError: email: `%v`", "codeflow")
						} else {
							log.Printf("Users::FindOne::Error: %s", err.Error())
						}
						return err
					}

					if err = NewRelease(feature, user, &release); err != nil {
						log.Printf("NewRelease::Error: %s", err.Error())
						return err
					}

					ReleaseCreated(&release)
				}
			} else {
				log.Printf("Features::FindOne::Error: %s", err.Error())
			}
			return err
		}

		log.Printf("Feature `%v:%v` already exists", project.Repository, payload.Hash)
	}

	if e.Name == "plugins.DockerBuild:status" {
		payload := e.Payload.(plugins.DockerBuild)
		DockerBuildStatus(&payload)
	}

	return nil
}
