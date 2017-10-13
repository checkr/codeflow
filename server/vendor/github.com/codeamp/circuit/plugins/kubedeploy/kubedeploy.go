package kubedeploy

import (
	log "github.com/codeamp/logger"
	"github.com/codeamp/transistor"
)

type KubeDeploy struct {
	events chan transistor.Event
}

func init() {
	transistor.RegisterPlugin("kubedeploy", func() transistor.Plugin {
		return &KubeDeploy{}
	})
}

func (x *KubeDeploy) Description() string {
	return "Deploy to Kubernetes"
}

func (x *KubeDeploy) SampleConfig() string {
	return ` `
}

func (x *KubeDeploy) Start(e chan transistor.Event) error {
	x.events = e
	log.Println("Started Kubedeploy")
	return nil
}

func (x *KubeDeploy) Stop() {
	log.Println("Stopping Kubedeploy")
}

func (x *KubeDeploy) Subscribe() []string {
	return []string{
		"plugins.DockerDeploy:create",
		"plugins.DockerDeploy:update",
		"plugins.DockerDeploy:destroy",
		"plugins.LoadBalancer:create",
		"plugins.LoadBalancer:update",
		"plugins.LoadBalancer:destroy",
	}
}

func (x *KubeDeploy) Process(e transistor.Event) error {
	log.Printf("Process KubeDeploy event: %s:%s", e.Name, e.ID)

	switch e.Name {
	case "plugins.DockerDeploy:create":
		x.doDeploy(e)
	case "plugins.DockerDeploy:update":
		x.doDeploy(e)
	case "plugins.DockerDeploy:destroy":
		x.doDeploy(e)
	case "plugins.LoadBalancer:create":
		x.doLoadBalancer(e)
	case "plugins.LoadBalancer:update":
		x.doLoadBalancer(e)
	case "plugins.LoadBalancer:destroy":
		x.doDeleteLoadBalancer(e)
	}

	log.Printf("Processed KubeDeploy event: %s:%s", e.Name, e.ID)

	return nil
}
