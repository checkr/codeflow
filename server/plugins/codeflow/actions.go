package codeflow

import (
	"log"
	"strings"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2/bson"
)

// CurrentProject finds project by slug
func CurrentProject(r *rest.Request, project *Project) error {
	slug := r.PathParam("slug")

	if err := db.Collection("projects").FindOne(bson.M{"slug": slug}, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindOne::DocumentNotFoundError: slug: `%v`", slug)
		} else {
			log.Printf("Projects::FindOne::Error: %s", err.Error())
		}
		return err
	}

	return nil
}

// CurrentUser returns current user identifed by JWT token
func CurrentUser(r *rest.Request, user *User) error {
	if err := db.Collection("users").FindOne(bson.M{"username": r.Env["REMOTE_USER"]}, &user); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Users::FindOne::DocumentNotFoundError: username: `%v`", r.Env["REMOTE_USER"])
		} else {
			log.Printf("Users::FindOne::Error: %s", err.Error())
		}
		return err
	}

	return nil
}

func ProjectCreated(p *Project) error {
	wsMsg := plugins.WebsocketMsg{
		Channel: "projects",
		Payload: p,
	}
	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func ServiceCreated(s *Service) error {
	project := Project{}

	if err := db.Collection("projects").FindById(s.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", s.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	// TODO: Magic

	return nil
}

func ServiceUpdated(s *Service) error {
	project := Project{}

	if err := db.Collection("projects").FindById(s.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", s.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	// TODO: Magic

	return nil
}

func ServiceDeleted(s *Service) error {
	project := Project{}

	if err := db.Collection("projects").FindById(s.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", s.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	// TODO: Magic

	return nil
}

func ReleaseCreated(r *Release) error {
	project := Project{}

	if err := db.Collection("projects").FindById(r.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", r.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	for _, str := range project.Workflows {
		s := strings.Split(str, "/")
		t, n := s[0], s[1]
		flow := Flow{
			ReleaseId: r.Id,
			Type:      t,
			Name:      n,
			State:     plugins.Waiting,
		}

		if err := db.Collection("workflows").Save(&flow); err != nil {
			log.Printf("Workflows::Save::Error: %v", err.Error())
			return err
		}

		r.Workflow = append(r.Workflow, flow)
	}

	if err := CheckWorkflows(r); err != nil {
		log.Printf("CheckWorkflows::Error: %v", err.Error())
		return err
	}

	if err := ReleaseUpdated(r); err != nil {
		log.Printf("ReleaseUpdated::Error: %v", err.Error())
		return err
	}

	return nil
}

func CheckWorkflows(r *Release) error {
	var workflowStatus plugins.State = plugins.Complete

	for idx, _ := range r.Workflow {
		flow := &r.Workflow[idx]

		switch flow.Type {
		case "build":
			var build Build
			if err := db.Collection("builds").FindOne(bson.M{"featureHash": r.HeadFeature.Hash, "type": flow.Name}, &build); err != nil {
				if _, ok := err.(*bongo.DocumentNotFoundError); ok {
					log.Printf("Builds::FindOne::DocumentNotFoundError: hash: `%v`, type: %v", r.HeadFeature.Hash, flow.Name)
				} else {
					log.Printf("Builds::FindOne::Error: %s", err.Error())
				}
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
			break

		case "ci":
			var externalFlowStatus ExternalFlowStatus
			results := db.Collection("externalFlowStatuses").Find(bson.M{"hash": r.HeadFeature.Hash})
			results.Query.Sort("-$natural").Limit(1)
			hasNext := results.Next(&externalFlowStatus)
			if !hasNext {
				if results.Error != nil {
					log.Printf("ExternalFlowStatuses::FindOne::Error: %s", results.Error.Error())
				} else {
					log.Printf("ExternalFlowStatuses::FindOne::DocumentNotFound: hash: `%v`", r.HeadFeature.Hash)
				}
				flow.State = plugins.Waiting
			} else {
				switch externalFlowStatus.State {
				case plugins.Waiting:
					flow.State = plugins.Waiting
				case plugins.Complete:
					flow.State = plugins.Complete
				case plugins.Failed:
					flow.State = plugins.Failed
					r.State = plugins.Failed
				default:
					flow.State = plugins.Running
				}
			}
			break
		}

		if err := db.Collection("workflows").Save(flow); err != nil {
			log.Printf("Workflows::Save::Error: %v", err.Error())
			return err
		}

		if flow.State != plugins.Complete {
			workflowStatus = flow.State
		}
	}

	if err := db.Collection("releases").Save(r); err != nil {
		log.Printf("Releases::Save::Error: %v", err.Error())
		return err
	}

	if workflowStatus == plugins.Complete {
		if err := CreateDeploy(r); err != nil {
			log.Printf("CreateDeploy::Error: %v", err.Error())
			return err
		}
	}

	return nil
}

func ExtensionCreated(lb *LoadBalancer) error {
	if err := loadBalancer(lb, plugins.Create); err != nil {
		log.Printf("loadBalancer::Error: %v", err.Error())
		return err
	}
	return nil
}

func ExtensionUpdated(lb *LoadBalancer) error {
	if err := loadBalancer(lb, plugins.Update); err != nil {
		log.Printf("loadBalancer::Error: %v", err.Error())
		return err
	}
	return nil
}

func ExtensionDeleted(lb *LoadBalancer) error {
	if err := loadBalancer(lb, plugins.Destroy); err != nil {
		log.Printf("loadBalancer::Error: %v", err.Error())
		return err
	}
	return nil
}

func loadBalancer(lb *LoadBalancer, action plugins.Action) error {
	project := Project{}
	service := Service{}

	if err := db.Collection("projects").FindById(lb.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", lb.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	if err := db.Collection("services").FindById(lb.ServiceId, &service); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Services::FindById::DocumentNotFoundError: _id: `%v`", lb.ServiceId)
		} else {
			log.Printf("Services::FindById::Error: %s", err.Error())
		}
		return err
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

	cf.Events <- agent.NewEvent(loadBalancerEvent, nil)

	return nil
}

func FeatureCreated(f *Feature, e agent.Event) error {
	project := Project{}
	secrets := []Secret{}

	if err := db.Collection("projects").FindById(f.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", f.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	// TODO: make type dynamic
	build := Build{
		FeatureHash: f.Hash,
		Type:        "DockerImage",
		State:       plugins.Waiting,
	}

	if err := db.Collection("builds").FindOne(bson.M{"featureHash": f.Hash}, &build); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Builds::Save: hash: `%v`", f.Hash)
			if err := db.Collection("builds").Save(&build); err != nil {
				log.Printf("Builds::Save::Error: %v", err.Error())
				return err
			}
		} else {
			log.Printf("Builds::FindOne::Error: %s", err.Error())
			return err
		}
	}

	results := db.Collection("secrets").Find(bson.M{"projectId": project.Id, "type": plugins.Build, "deleted": false})
	secret := Secret{}
	for results.Next(&secret) {
		secrets = append(secrets, secret)
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
			Hash:       f.Hash,
			ParentHash: f.ParentHash,
			User:       f.User,
			Message:    f.Message,
		},
		Registry: plugins.DockerRegistry{
			Host:     viper.GetString("plugins.docker_build.registry_host"),
			Username: viper.GetString("plugins.docker_build.registry_username"),
			Password: viper.GetString("plugins.docker_build.registry_password"),
			Email:    viper.GetString("plugins.docker_build.registry_user_email"),
		},
		BuildArgs: buildArgs,
	}

	cf.Events <- e.NewEvent(dockerBuildEvent, nil)

	wsMsg := plugins.WebsocketMsg{
		Channel: "features",
		Payload: f,
	}

	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func DockerBuildStatus(br *plugins.DockerBuild) error {
	project := Project{}
	build := Build{}
	feature := Feature{}

	if err := db.Collection("projects").FindOne(bson.M{"repository": br.Project.Repository}, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindOne::DocumentNotFoundError: repository: `%v`", br.Project.Repository)
		} else {
			log.Printf("Projects::FindOne::Error: %s", err.Error())
		}
		return err
	}

	if err := db.Collection("builds").FindOne(bson.M{"featureHash": br.Feature.Hash}, &build); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Builds::FindOne::DocumentNotFoundError: featureHash: `%v`", br.Feature.Hash)
		} else {
			log.Printf("Builds::FindOne::Error: %s", err.Error())
		}
	}

	build.Image = br.Image
	build.State = br.State
	build.BuildLog = br.BuildLog
	build.BuildError = br.BuildError

	if err := db.Collection("builds").Save(&build); err != nil {
		log.Printf("Builds::Save::Error: %v", err.Error())
		return err
	}

	if err := db.Collection("features").FindOne(bson.M{"hash": br.Feature.Hash}, &feature); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Features::FindOne::DocumentNotFoundError: hash: `%v`", br.Feature.Hash)
		} else {
			log.Printf("Features::FindOne::Error: %s", err.Error())
		}
		return err
	}

	if err := UpdateInProgessReleases(&feature); err != nil {
		log.Printf("UpdateInProgessReleases::Error: %s", err.Error())
	}

	return nil
}

func UpdateInProgessReleases(f *Feature) error {
	release := Release{}

	results := db.Collection("releases").Find(bson.M{"featureId": f.Id, "state": plugins.Waiting})
	for results.Next(&release) {
		if err := CheckWorkflows(&release); err != nil {
			log.Printf("CheckWorkflows::Error: %s, releaseId: %v"+err.Error(), release.Id)
		}
		if err := ReleaseUpdated(&release); err != nil {
			log.Printf("ReleaseUpdated::Error: %s, releaseId: %v"+err.Error(), release.Id)
		}
	}

	return nil
}

func NewRelease(f Feature, u User, release *Release) error {
	secrets := []Secret{}
	secret := Secret{}
	services := []Service{}
	service := Service{}

	results := db.Collection("secrets").Find(bson.M{"projectId": f.ProjectId, "type": bson.M{"$in": []plugins.Type{plugins.Env, plugins.File}}, "deleted": false})
	for results.Next(&secret) {
		secrets = append(secrets, secret)
	}

	results = db.Collection("services").Find(bson.M{"projectId": f.ProjectId, "state": bson.M{"$in": []plugins.State{plugins.Waiting, plugins.Running}}})
	for results.Next(&service) {
		services = append(services, service)
	}

	release.ProjectId = f.ProjectId
	release.HeadFeatureId = f.Id
	release.HeadFeature = f
	release.UserId = u.Id
	release.User = u
	release.State = plugins.Waiting
	release.Secrets = secrets
	release.Services = services

	currentRelease := Release{}
	if err := GetCurrentRelease(f.ProjectId, &currentRelease); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("GetCurrentRelease::DocumentNotFound: %v", err.Error())
			release.TailFeatureId = f.Id
			release.TailFeature = f
		} else {
			log.Printf("GetCurrentRelease::Error: %s", err.Error())
			return err
		}
	} else {
		release.TailFeatureId = currentRelease.HeadFeature.Id
		release.TailFeature = currentRelease.HeadFeature
	}

	if err := db.Collection("releases").Save(release); err != nil {
		log.Printf("Releases::Save::Error: %v", err.Error())
		return err
	}

	return nil
}

func CreateDeploy(r *Release) error {
	project := Project{}
	build := Build{}

	if err := db.Collection("projects").FindById(r.ProjectId, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindById::DocumentNotFoundError: _id: `%v`", r.ProjectId)
		} else {
			log.Printf("Projects::FindById::Error: %s", err.Error())
		}
		return err
	}

	if err := db.Collection("builds").FindOne(bson.M{"featureHash": r.HeadFeature.Hash, "state": plugins.Complete}, &build); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Builds::FindOne::DocumentNotFoundError: featureHash: `%v`, state: %v", r.HeadFeature.Hash, plugins.Complete)
		} else {
			log.Printf("Builds::FindOne::Error: %s", err.Error())
		}
		return err
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
	for _, service := range r.Services {
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
				Host:     viper.GetString("plugins.docker_build.registry_host"),
				Username: viper.GetString("plugins.docker_build.registry_username"),
				Password: viper.GetString("plugins.docker_build.registry_password"),
				Email:    viper.GetString("plugins.docker_build.registry_user_email"),
			},
		},
		Services:    services,
		Secrets:     secrets,
		Environment: "development",
	}

	cf.Events <- agent.NewEvent(dockerDeployEvent, nil)

	return nil
}

func LoadBalancerStatus(lb *plugins.LoadBalancer) error {
	project := Project{}
	extension := LoadBalancer{}

	if err := db.Collection("projects").FindOne(bson.M{"slug": lb.Project.Slug}, &project); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Projects::FindOne::DocumentNotFoundError: slug: `%v`", lb.Project.Slug)
		} else {
			log.Printf("Projects::FindOne::Error: %s", err.Error())
		}
		return err
	}

	if err := db.Collection("extensions").FindOne(bson.M{"projectId": project.Id, "name": lb.Name}, &extension); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Extensions::FindOne::DocumentNotFoundError: projectId: `%v`, name: `%v`", project.Id, lb.Name)
		} else {
			log.Printf("Extensions::FindOne::Error: %s", err.Error())
		}
		return err
	}

	extension.DNSName = lb.DNSName
	extension.State = lb.State
	extension.StateMessage = lb.StateMessage

	if err := db.Collection("extensions").Save(&extension); err != nil {
		log.Printf("Extension::Save::Error: %v", err.Error())
		return err
	}

	return nil
}

func DockerDeployStatus(e *plugins.DockerDeploy) error {
	release := Release{}

	if err := db.Collection("releases").FindById(bson.ObjectIdHex(e.Release.Id), &release); err != nil {
		if _, ok := err.(*bongo.DocumentNotFoundError); ok {
			log.Printf("Releases::FindById::DocumentNotFoundError: _id: `%v`", e.Release.Id)
		} else {
			log.Printf("Releases::FindById::Error: %s", err.Error())
		}
		return err
	}

	release.State = e.State

	if err := db.Collection("releases").Save(&release); err != nil {
		log.Printf("Releases::Save::Error: %v", err.Error())
		return err
	}

	if release.State == plugins.Complete {
		if err := ReleasePromoted(&release); err != nil {
			log.Printf("ReleasePromoted::Error: %v", err.Error())
		}
	}

	if err := FeatureUpdated(&release); err != nil {
		log.Printf("FeatureUpdated::Error: %v", err.Error())
	}

	if err := ReleaseUpdated(&release); err != nil {
		log.Printf("ReleaseUpdated::Error: %v", err.Error())
	}

	return nil
}

func ReleaseUpdated(r *Release) error {
	wsMsg := plugins.WebsocketMsg{
		Channel: "releases",
		Payload: r,
	}

	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func FeatureUpdated(r *Release) error {
	wsMsg := plugins.WebsocketMsg{
		Channel: "features",
		Payload: r.HeadFeature,
	}

	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func ReleasePromoted(r *Release) error {
	wsMsg := plugins.WebsocketMsg{
		Channel: "releases/promote",
		Payload: r,
	}

	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func BookmarksUpdated(u *User) error {
	bookmarks := []Bookmark{}
	bookmark := Bookmark{}

	results := db.Collection("bookmarks").Find(bson.M{"userId": u.Id})
	for results.Next(&bookmark) {
		bookmarks = append(bookmarks, bookmark)
	}

	wsMsg := plugins.WebsocketMsg{
		Channel: "bookmarks/" + u.Id.Hex(),
		Payload: bookmarks,
	}

	cf.Events <- agent.NewEvent(wsMsg, nil)

	return nil
}

func GetCurrentRelease(projectId bson.ObjectId, release *Release) error {
	results := db.Collection("releases").Find(bson.M{"projectId": projectId, "state": plugins.Complete})
	results.Query.Sort("-$natural").Limit(1)

	hasNext := results.Next(release)

	if !hasNext {
		// There could have been an error fetching the next one, which would set the Error property on the resultset
		if results.Error != nil {
			return results.Error
		} else {
			return &bongo.DocumentNotFoundError{}
		}
	}

	return nil
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

func StringToState(s string) plugins.State {
	switch strings.ToLower(s) {
	case "pending":
		return plugins.Running
	case "success":
		return plugins.Complete
	case "failed":
		return plugins.Failed
	}
	return plugins.Waiting
}

func CollectStats(save bool, stats *Statistics) error {
	// Project stats
	if count, err := db.Collection("projects").Collection().Find(bson.M{}).Count(); err != nil {
		return err
	} else {
		stats.Projects = count
	}

	// Feature stats
	if count, err := db.Collection("projects").Collection().Find(bson.M{}).Count(); err != nil {
		return err
	} else {
		stats.Features = count
	}

	// Release stats
	if count, err := db.Collection("releases").Collection().Find(bson.M{}).Count(); err != nil {
		return err
	} else {
		stats.Releases = count
	}

	// User stats
	if count, err := db.Collection("users").Collection().Find(bson.M{}).Count(); err != nil {
		return err
	} else {
		stats.Users = count
	}

	return nil
}
