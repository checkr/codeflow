package testdata

import (
	"fmt"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
)

// LoadBalancers

func TearDownLBTCP(t plugins.Type) agent.Event {
	// Start with destroy, then create.
	lbe := LBDataForTCP(plugins.Destroy, t)
	event := agent.NewEvent(lbe, nil)
	return event
}

func TearDownLBHTTPS(t plugins.Type) agent.Event {
	// Start with destroy, then create.
	lbe := LBDataForTCP(plugins.Destroy, t)
	lbe.ListenerPairs[0].Destination.Protocol = "HTTPS"
	lbe.Name = "nginx-test-lb-https"
	event := agent.NewEvent(lbe, nil)
	return event
}

func CreateLBTCP(t plugins.Type) agent.Event {
	lbe := LBDataForTCP(plugins.Create, t)
	event := agent.NewEvent(lbe, nil)
	return event
}

func UpdateLBTCP(t plugins.Type) agent.Event {
	lbe := UpdateLBDataForTCP(plugins.Update, t)
	event := agent.NewEvent(lbe, nil)
	return event
}

func CreateLBHTTPS(t plugins.Type) agent.Event {
	lbe := LBDataForTCP(plugins.Create, t)
	lbe.ListenerPairs[0].Destination.Protocol = "HTTPS"
	lbe.Name = "nginx-test-lb-https"
	event := agent.NewEvent(lbe, nil)
	return event
}

func UpdateLBHTTPS(t plugins.Type) agent.Event {
	lbe := UpdateLBDataForTCP(plugins.Update, t)
	lbe.ListenerPairs[0].Destination.Protocol = "HTTPS"
	lbe.ListenerPairs[1].Destination.Protocol = "HTTPS"
	lbe.Name = "nginx-test-lb-https"
	event := agent.NewEvent(lbe, nil)
	return event
}

func LBDataForTCP(action plugins.Action, t plugins.Type) plugins.LoadBalancer {
	project := plugins.Project{
		Slug:       "nginx-test-success",
		Repository: "checkr/nginx-test-success",
	}
	service := plugins.Service{
		Action:  action,
		Name:    "nginx",
		Command: "nginx -g 'daemon off;'",
		Listeners: []plugins.Listener{
			{
				Port:     80,
				Protocol: "TCP",
			},
		},
		State: plugins.Waiting,
		Spec: plugins.ServiceSpec{
			CpuRequest:                    "500m",
			CpuLimit:                      "1000m",
			MemoryRequest:                 "512Mi",
			MemoryLimit:                   "1Gi",
			TerminationGracePeriodSeconds: int64(600),
		},
		Replicas: 1,
	}
	lbe := plugins.LoadBalancer{
		Name:        "nginx-test-lb-asdf1234",
		Action:      action,
		Environment: "testing2",
		Type:        t,
		Project:     project,
		Service:     service,
		ListenerPairs: []plugins.ListenerPair{
			{
				Source:      plugins.Listener{Port: 443, Protocol: "TCP"},
				Destination: plugins.Listener{Port: 80, Protocol: "TCP"},
			},
		},
		Subdomain: "nginx-testing.checkrhq-dev.net",
	}
	return lbe
}

func UpdateLBDataForTCP(action plugins.Action, t plugins.Type) plugins.LoadBalancer {
	project := plugins.Project{
		Slug: "nginx-test-success",
	}
	service := plugins.Service{
		Action:  action,
		Name:    "nginx",
		Command: "nginx -g 'daemon off;'",
		Listeners: []plugins.Listener{
			{
				Port:     3000,
				Protocol: "TCP",
			},
			{
				Port:     3001,
				Protocol: "TCP",
			},
		},
		State: plugins.Waiting,
		Spec: plugins.ServiceSpec{
			CpuRequest:                    "500m",
			CpuLimit:                      "1000m",
			MemoryRequest:                 "512Mi",
			MemoryLimit:                   "1Gi",
			TerminationGracePeriodSeconds: int64(600),
		},
		Replicas: 1,
	}
	lbe := plugins.LoadBalancer{
		Name:        "nginx-test-lb-asdf1234",
		Action:      action,
		Environment: "testing2",
		Type:        t,
		Project:     project,
		Service:     service,
		ListenerPairs: []plugins.ListenerPair{
			{
				Source:      plugins.Listener{Port: 80, Protocol: "TCP"},
				Destination: plugins.Listener{Port: 3000, Protocol: "TCP"},
			},
			{
				Source:      plugins.Listener{Port: 443, Protocol: "TCP"},
				Destination: plugins.Listener{Port: 3001, Protocol: "TCP"},
			},
		},
	}
	return lbe
}

// Deploys

func TearDownPreviousDeploys(ag agent.Agent) {
	ag.Events <- TeardownPreviousDeploy("nginx-test-success")
	ag.Events <- TeardownPreviousDeploy("nginx-test-failure")
	ag.Events <- TeardownPreviousDeploy("checkr-codeflow")
}

func TeardownPreviousDeploy(name string) agent.Event {
	deploy := DeployData("nginx-test-success", plugins.Destroy)
	event := agent.NewEvent(deploy, nil)
	return event
}

func CreateSuccessDeploy() agent.Event {
	deploy := DeployData("nginx-test-success", plugins.Create)
	event := agent.NewEvent(deploy, nil)
	return event
}

func CreateSuccessDeployRenamed() agent.Event {
	deploy := DeployDataRenamed("nginx-test-success", plugins.Create)
	event := agent.NewEvent(deploy, nil)
	return event
}

func CreateDockerSocketDeploy() agent.Event {
	deploy := DeployData("checkr-codeflow", plugins.Create)
	deploy.Environment = "testing2"
	event := agent.NewEvent(deploy, nil)
	return event
}

func CreateSuccessMixedActionDeploy() agent.Event {
	actions := []plugins.Action{plugins.Update, plugins.Update, plugins.Destroy, plugins.Create}
	deploy := DeployDataMixedActions("nginx-test-success", actions)
	return agent.NewEvent(deploy, nil)
}

func CreateFailDeploy() agent.Event {
	deploy := DeployDataFail("nginx-test-failure", plugins.Create)
	// Set an Image that's invalid so we can test failure
	deploy.Docker.Image = "checkr/deploy-test:INVALID"
	event := agent.NewEvent(deploy, nil)
	return event
}

func CreateFailDeployCommand() agent.Event {
	deploy := DeployDataFail("nginx-test-failure", plugins.Create)
	// Set an Image that's invalid so we can test failure
	event := agent.NewEvent(deploy, nil)
	return event
}

func DeleteSuccessDeploy() agent.Event {
	deploy := DeployData("nginx-test-success", plugins.Destroy)
	event := agent.NewEvent(deploy, nil)
	return event
}

func DeleteFailedDeploy() agent.Event {
	deploy := DeployData("nginx-test-DOESNTEXIST", plugins.Destroy)
	event := agent.NewEvent(deploy, nil)
	return event
}

func DeployDataMixedActions(name string, actions []plugins.Action) plugins.DockerDeploy {
	project := plugins.Project{
		Slug: name,
	}

	headFeature := plugins.Feature{
		Message:    "test1",
		User:       "jeremy@checkr.com",
		Hash:       "112",
		ParentHash: "112",
	}

	tailFeature := plugins.Feature{
		Message:    "test2",
		User:       "jeremy@checkr.com",
		Hash:       "456",
		ParentHash: "456",
	}

	release := plugins.Release{
		HeadFeature: headFeature,
		TailFeature: tailFeature,
	}

	listener := plugins.Listener{
		Port:     80,
		Protocol: "TCP",
	}

	var serviceArray []plugins.Service
	for i, action := range actions {
		serviceArray = append(serviceArray, plugins.Service{
			Action:    action,
			Name:      fmt.Sprintf("nginx%d", i),
			Command:   "nginx -g 'daemon off;'",
			Listeners: []plugins.Listener{listener},
			State:     plugins.Waiting,
			Spec: plugins.ServiceSpec{
				CpuRequest:                    "500m",
				CpuLimit:                      "1000m",
				MemoryRequest:                 "512Mi",
				MemoryLimit:                   "1Gi",
				TerminationGracePeriodSeconds: int64(600),
			},
			Replicas: 1,
		})
	}

	docker := plugins.Docker{
		Image: "checkr/deploy-test:latest",
	}

	kubeDeploy := plugins.DockerDeploy{
		Action:      plugins.Create,
		Docker:      docker,
		Environment: "testing2",
		Project:     project,
		Timeout:     60,
		Release:     release,
		Services:    serviceArray,
		Secrets: []plugins.Secret{
			{
				Key:   "MY_SECRET_KEY",
				Value: "MY_SECRET_VALUE",
				Type:  plugins.Env,
			},
		},
	}
	return kubeDeploy
}

func DeployDataRenamed(name string, action plugins.Action) plugins.DockerDeploy {
	project := plugins.Project{
		Slug: name,
	}

	headFeature := plugins.Feature{
		Message:    "test1",
		User:       "jeremy@checkr.com",
		Hash:       "112",
		ParentHash: "112",
	}

	tailFeature := plugins.Feature{
		Message:    "test2",
		User:       "jeremy@checkr.com",
		Hash:       "456",
		ParentHash: "456",
	}

	release := plugins.Release{
		HeadFeature: headFeature,
		TailFeature: tailFeature,
	}

	listener := plugins.Listener{
		Port:     80,
		Protocol: "TCP",
	}

	var serviceArray []plugins.Service

	// Two web services
	for i := 0; i < 2; i++ {
		serviceArray = append(serviceArray, plugins.Service{
			Action:    action,
			Name:      fmt.Sprintf("newguy%d", i),
			Command:   "nginx -g 'daemon off;'",
			Listeners: []plugins.Listener{listener},
			State:     plugins.Waiting,
			Spec: plugins.ServiceSpec{
				CpuRequest:                    "500m",
				CpuLimit:                      "1000m",
				MemoryRequest:                 "512Mi",
				MemoryLimit:                   "1Gi",
				TerminationGracePeriodSeconds: int64(600),
			},
			Replicas: 1,
		})
	}
	// One worker
	serviceArray = append(serviceArray, plugins.Service{
		Action:  action,
		Name:    "worker",
		Command: "/bin/sh -c 'while(/bin/true); do sleep 1; echo waiting forever...; done'",
		State:   plugins.Waiting,
		Spec: plugins.ServiceSpec{
			CpuRequest:                    "500m",
			CpuLimit:                      "1000m",
			MemoryRequest:                 "512Mi",
			MemoryLimit:                   "1Gi",
			TerminationGracePeriodSeconds: int64(600),
		},
		Replicas: 1,
	})

	docker := plugins.Docker{
		Image: "checkr/deploy-test:latest",
	}

	kubeDeploy := plugins.DockerDeploy{
		Action:      action,
		Docker:      docker,
		Environment: "testing2",
		Project:     project,
		Timeout:     60,
		Release:     release,
		Services:    serviceArray,
		Secrets: []plugins.Secret{
			{
				Key:   "MY_SECRET_KEY",
				Value: "MY_SECRET_VALUE",
				Type:  plugins.Env,
			},
		},
	}
	return kubeDeploy
}

func DeployData(name string, action plugins.Action) plugins.DockerDeploy {
	project := plugins.Project{
		Slug: name,
	}

	headFeature := plugins.Feature{
		Message:    "test1",
		User:       "jeremy@checkr.com",
		Hash:       "112",
		ParentHash: "112",
	}

	tailFeature := plugins.Feature{
		Message:    "test2",
		User:       "jeremy@checkr.com",
		Hash:       "456",
		ParentHash: "456",
	}

	release := plugins.Release{
		HeadFeature: headFeature,
		TailFeature: tailFeature,
	}

	listener := plugins.Listener{
		Port:     80,
		Protocol: "TCP",
	}

	var serviceArray []plugins.Service

	// Two web services
	for i := 0; i < 2; i++ {
		serviceArray = append(serviceArray, plugins.Service{
			Action:    action,
			Name:      fmt.Sprintf("nginx%d", i),
			Command:   "nginx -g 'daemon off;'",
			Listeners: []plugins.Listener{listener},
			State:     plugins.Waiting,
			Spec: plugins.ServiceSpec{
				CpuRequest:                    "500m",
				CpuLimit:                      "1000m",
				MemoryRequest:                 "512Mi",
				MemoryLimit:                   "1Gi",
				TerminationGracePeriodSeconds: int64(600),
			},
			Replicas: 1,
		})
	}
	// One worker
	serviceArray = append(serviceArray, plugins.Service{
		Action:  action,
		Name:    "worker",
		Command: "/bin/sh -c 'while(/bin/true); do sleep 1; echo waiting forever...; done'",
		State:   plugins.Waiting,
		Spec: plugins.ServiceSpec{
			CpuRequest:                    "500m",
			CpuLimit:                      "1000m",
			MemoryRequest:                 "512Mi",
			MemoryLimit:                   "1Gi",
			TerminationGracePeriodSeconds: int64(600),
		},
		Replicas: 1,
	})

	docker := plugins.Docker{
		Image: "checkr/deploy-test:latest",
	}

	kubeDeploy := plugins.DockerDeploy{
		Action:      action,
		Docker:      docker,
		Environment: "testing2",
		Project:     project,
		Timeout:     60,
		Release:     release,
		Services:    serviceArray,
		Secrets: []plugins.Secret{
			{
				Key:   "MY_SECRET_KEY",
				Value: "MY_SECRET_VALUE",
				Type:  plugins.Env,
			},
		},
	}
	return kubeDeploy
}

func DeployDataFail(name string, action plugins.Action) plugins.DockerDeploy {
	project := plugins.Project{
		Slug: name,
	}

	headFeature := plugins.Feature{
		Message:    "test1",
		User:       "jeremy@checkr.com",
		Hash:       "112",
		ParentHash: "112",
	}

	tailFeature := plugins.Feature{
		Message:    "test2",
		User:       "jeremy@checkr.com",
		Hash:       "456",
		ParentHash: "456",
	}

	release := plugins.Release{
		HeadFeature: headFeature,
		TailFeature: tailFeature,
	}

	listener := plugins.Listener{
		Port:     80,
		Protocol: "TCP",
	}

	var serviceArray []plugins.Service

	// Two web services
	for i := 0; i < 2; i++ {
		serviceArray = append(serviceArray, plugins.Service{
			Action:    action,
			Name:      fmt.Sprintf("nginx%d", i),
			Command:   "boom",
			Listeners: []plugins.Listener{listener},
			State:     plugins.Waiting,
			Spec: plugins.ServiceSpec{
				CpuRequest:                    "500m",
				CpuLimit:                      "1000m",
				MemoryRequest:                 "512Mi",
				MemoryLimit:                   "1Gi",
				TerminationGracePeriodSeconds: int64(600),
			},
			Replicas: 1,
		})
	}
	// One worker
	serviceArray = append(serviceArray, plugins.Service{
		Action:  action,
		Name:    "worker",
		Command: "boom",
		State:   plugins.Waiting,
		Spec: plugins.ServiceSpec{
			CpuRequest:                    "500m",
			CpuLimit:                      "1000m",
			MemoryRequest:                 "512Mi",
			MemoryLimit:                   "1Gi",
			TerminationGracePeriodSeconds: int64(600),
		},
		Replicas: 1,
	})

	docker := plugins.Docker{
		Image: "checkr/deploy-test:latest",
	}

	kubeDeploy := plugins.DockerDeploy{
		Action:      action,
		Docker:      docker,
		Environment: "testing2",
		Project:     project,
		Timeout:     60,
		Release:     release,
		Services:    serviceArray,
		Secrets: []plugins.Secret{
			{
				Key:   "MY_SECRET_KEY",
				Value: "MY_SECRET_VALUE",
				Type:  plugins.Env,
			},
		},
	}
	return kubeDeploy
}
