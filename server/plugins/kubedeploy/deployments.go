package kubedeploy

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/pkg/api/errors"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	metav1 "k8s.io/client-go/pkg/apis/meta/v1"
	"k8s.io/client-go/pkg/util"
	"k8s.io/client-go/pkg/util/intstr"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/extemporalgenome/slug"
	"github.com/google/shlex"
	"github.com/spf13/viper"
)

func genNamespaceName(suggestedEnvironment string, projectSlug string) string {
	var environment string
	if viper.IsSet("plugins.kubedeploy.environment") {
		environment = viper.GetString("plugins.kubedeploy.environment")
	} else {
		environment = suggestedEnvironment
	}
	return fmt.Sprintf("%s-%s", environment, projectSlug)
}

func genDeploymentName(repoName string, serviceName string) string {
	return slug.Slug(repoName) + "-" + serviceName
}

func (x *KubeDeploy) sendDDResponse(e agent.Event, services []plugins.Service, state plugins.State, failureMessage string) {
	data := e.Payload.(plugins.DockerDeploy)
	data.Action = plugins.Status
	data.State = state
	data.Services = services
	data.StateMessage = failureMessage
	event := e.NewEvent(data, nil)
	x.events <- event
}

func (x *KubeDeploy) sendDDSuccessResponse(e agent.Event, services []plugins.Service) {
	x.sendDDResponse(e, services, plugins.Complete, "")
}

func (x *KubeDeploy) sendDDErrorResponse(e agent.Event, services []plugins.Service, failureMessage string) {
	x.sendDDResponse(e, services, plugins.Failed, failureMessage)
}

func (x *KubeDeploy) sendDDInProgress(e agent.Event, services []plugins.Service, message string) {
	x.sendDDResponse(e, services, plugins.Running, message)
}

func secretifyDockerCred() string {
	encodeMe := fmt.Sprintf("%s:%s",
		viper.GetString("plugins.docker_build.registry_username"),
		viper.GetString("plugins.docker_build.registry_password"))
	encodeResult := []byte(encodeMe)
	authField := base64.StdEncoding.EncodeToString(encodeResult)
	jsonFilled := fmt.Sprintf("{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"email\":\"%s\",\"auth\":\"%s\"}}",
		viper.GetString("plugins.docker_build.registry_host"),
		viper.GetString("plugins.docker_build.registry_username"),
		viper.GetString("plugins.docker_build.registry_password"),
		viper.GetString("plugins.docker_build.registry_user_email"),
		authField,
	)
	return jsonFilled
}

func (x *KubeDeploy) createDockerIOSecretIfNotExists(namespace string, coreInterface v1core.CoreV1Interface) error {
	// Load up the docker-io secrets for image pull if not exists
	_, dockerIOSecretErr := coreInterface.Secrets(namespace).Get("docker-io", metav1.GetOptions{})
	if dockerIOSecretErr != nil {
		if errors.IsNotFound(dockerIOSecretErr) {
			log.Printf("docker-io secret not found for %s, creating.", namespace)
			secretMap := map[string]string{
				".dockercfg": secretifyDockerCred(),
			}
			_, createDockerIOSecretErr := coreInterface.Secrets(namespace).Create(&v1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "docker-io",
					Namespace: namespace,
				},
				StringData: secretMap,
				Type:       v1.SecretTypeDockercfg,
			})
			if createDockerIOSecretErr != nil {
				log.Printf("Error '%s' creating docker-io secret for %s.", createDockerIOSecretErr, namespace)
				return createDockerIOSecretErr
			}
		} else {
			log.Printf("Error unhandled '%s' while attempting to lookup docker-io secret.", dockerIOSecretErr)
			return dockerIOSecretErr
		}
	}
	return nil
}

func (x *KubeDeploy) createNamespaceIfNotExists(namespace string, coreInterface v1core.CoreV1Interface) error {
	// Create namespace if it does not exist.
	_, nameGetErr := coreInterface.Namespaces().Get(namespace, metav1.GetOptions{})
	if nameGetErr != nil {
		if errors.IsNotFound(nameGetErr) {
			log.Printf("Namespace %s does not yet exist, creating.", namespace)
			namespaceParams := &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name: namespace,
				},
			}
			_, createNamespaceErr := coreInterface.Namespaces().Create(namespaceParams)
			if createNamespaceErr != nil {
				log.Printf("Error '%s' creating namespace %s", createNamespaceErr, namespace)
				return createNamespaceErr
			}
			log.Printf("Namespace created: %s", namespace)
		} else {
			log.Printf("Unhandled error occured looking up namespace %s: '%s'", namespace, nameGetErr)
			return nameGetErr
		}
	}
	return nil
}

func (x *KubeDeploy) doDeploy(e agent.Event) error {
	// Codeflow will load the kube config from a file, specified by KUBECONFIG environment variable
	kubeconfig := os.Getenv("KUBECONFIG")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Error getting cluster config.")
	}

	data := e.Payload.(plugins.DockerDeploy)
	x.sendDDInProgress(e, data.Services, "Deploy in-progress")
	namespace := genNamespaceName(data.Environment, data.Project.Slug)
	coreInterface := clientset.Core()

	successfulDeploys := 0
	timeout := e.Payload.(plugins.DockerDeploy).Timeout
	// Set default timeout to 600 seconds if not specified.
	if timeout == 0 {
		timeout = 600
	}
	curTime := 0

	var rollbackTargets []plugins.Service
	var allDeleteSuccess int
	totalDeploysRequested := 0
	var totalDeletesRequested int
	for _, service := range data.Services {
		switch service.Action {
		case plugins.Create:
			totalDeploysRequested++
		case plugins.Update:
			totalDeploysRequested++
		case plugins.Destroy:
			totalDeploysRequested++
			totalDeletesRequested++
		}
	}

	createNamespaceErr := x.createNamespaceIfNotExists(namespace, coreInterface)
	if createNamespaceErr != nil {
		x.sendDDErrorResponse(e, data.Services, createNamespaceErr.Error())
		return nil
	}

	createDockerIOSecretErr := x.createDockerIOSecretIfNotExists(namespace, coreInterface)
	if createDockerIOSecretErr != nil {
		x.sendDDErrorResponse(e, data.Services, createDockerIOSecretErr.Error())
		return nil
	}

	// Create secrets for this deploy
	var secretMap map[string]string
	secretMap = make(map[string]string)
	var myEnvVars []v1.EnvVar

	// This map is used in to create the secrets themselves
	for _, secret := range data.Secrets {
		secretMap[secret.Key] = secret.Value
	}

	secretParams := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v1.ObjectMeta{
			GenerateName: fmt.Sprintf("%v-", data.Project.Slug),
			Namespace:    namespace,
		},
		StringData: secretMap,
		Type:       v1.SecretTypeOpaque,
	}

	secretResult, secErr := coreInterface.Secrets(namespace).Create(secretParams)
	if secErr != nil {
		failMessage := fmt.Sprintf("Error '%s' creating secret %s", secErr, data.Project.Slug)
		x.sendDDErrorResponse(e, data.Services, failMessage)
		return nil
	}
	secretName := secretResult.Name
	log.Printf("Secrets created: %s", secretName)

	// This is for building the configuration to use the secrets from inside the deployment
	// as ENVs
	for _, secret := range data.Secrets {
		if secret.Type == plugins.Env {
			newEnv := v1.EnvVar{
				Name: secret.Key,
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: secretName,
						},
						Key: secret.Key,
					},
				},
			}
			myEnvVars = append(myEnvVars, newEnv)
		}
	}
	// as Files
	var volumeMounts []v1.VolumeMount
	var deployVolumes []v1.Volume
	var volumeSecretItems []v1.KeyToPath
	volumeMounts = append(volumeMounts, v1.VolumeMount{
		Name:      secretName,
		MountPath: "/etc/secrets",
		ReadOnly:  true,
	})
	for _, secret := range data.Secrets {
		if secret.Type == plugins.File {
			volumeSecretItems = append(volumeSecretItems, v1.KeyToPath{
				Path: secret.Key,
				Key:  secret.Key,
				Mode: util.Int32Ptr(256),
			})
		}
	}
	secretVolume := v1.SecretVolumeSource{
		SecretName:  secretName,
		Items:       volumeSecretItems,
		DefaultMode: util.Int32Ptr(256),
	}

	// Add the secrets
	deployVolumes = append(deployVolumes, v1.Volume{
		Name: secretName,
		VolumeSource: v1.VolumeSource{
			Secret: &secretVolume,
		},
	})

	// Do update/create of deployments and services
	depInterface := clientset.Extensions()

	// Validate we have some services to deploy
	if len(data.Services) == 0 {
		failMessage := fmt.Sprintf("ERROR: Zero services were found in the deploy message.")
		x.sendDDErrorResponse(e, data.Services, failMessage)
		return nil
	}

	// Codeflow docker building container requires docker socket.
	if data.Project.Slug == "checkr-codeflow" {
		deployVolumes = append(deployVolumes, v1.Volume{
			Name: "dockersocket",
			VolumeSource: v1.VolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/var/run/docker.sock",
				},
			},
		})
		volumeMounts = append(volumeMounts, v1.VolumeMount{
			Name:      "dockersocket",
			ReadOnly:  false,
			MountPath: "/var/run/docker.sock",
		})
	}

	for index, service := range data.Services {
		deploymentName := genDeploymentName(data.Project.Slug, service.Name)
		var deployPorts []v1.ContainerPort

		// ContainerPorts for the deployment
		for _, cPort := range service.Listeners {
			// Build the deployments containerports array
			newContainerPort := v1.ContainerPort{
				ContainerPort: cPort.Port,
				Protocol:      v1.Protocol(cPort.Protocol),
			}
			deployPorts = append(deployPorts, newContainerPort)
		}

		// Support ready and liveness probes
		var readyProbe v1.Probe
		var liveProbe v1.Probe
		var deployStrategy v1beta1.DeploymentStrategy
		if len(service.Listeners) >= 1 && service.Listeners[0].Protocol == "TCP" {
			// If the service is TCP, use a TCP Probe
			myPort := service.Listeners[0].Port
			readyProbe = v1.Probe{
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{IntVal: myPort},
					},
				},
			}
			liveProbe = v1.Probe{
				InitialDelaySeconds: 5,
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{IntVal: myPort},
					},
				},
			}
			deployStrategy = v1beta1.DeploymentStrategy{
				Type: v1beta1.RollingUpdateDeploymentStrategyType,
				RollingUpdate: &v1beta1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "30%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "60%",
					},
				},
			}
		} else {
			// If the service is non-TCP or has no ports use a simple exec probe
			runThis := []string{"/bin/true"}
			readyProbe = v1.Probe{
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: runThis,
					},
				},
			}
			liveProbe = v1.Probe{
				InitialDelaySeconds: 5,
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: runThis,
					},
				},
			}
			deployStrategy = v1beta1.DeploymentStrategy{
				Type: "Recreate",
			}
		}

		// Deployment
		replicas := int32(service.Replicas)
		if service.Action == plugins.Destroy {
			replicas = 0
		}

		// Command parsing into entrypoint vs. args
		commandArray, _ := shlex.Split(service.Command)
		commandEntryPoint := commandArray[0]
		var commandArgs []string
		if len(commandArray) > 1 {
			commandArgs = commandArray[1:]
		}

		var revisionHistoryLimit int32 = 10
		terminationGracePeriodSeconds := int64(600)
		deployParams := &v1beta1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "extensions/v1beta1",
			},
			ObjectMeta: v1.ObjectMeta{
				Name: deploymentName,
			},
			Spec: v1beta1.DeploymentSpec{
				Replicas:             &replicas,
				Strategy:             deployStrategy,
				RevisionHistoryLimit: &revisionHistoryLimit,
				Template: v1.PodTemplateSpec{
					ObjectMeta: v1.ObjectMeta{
						Name:   deploymentName,
						Labels: map[string]string{"app": deploymentName},
					},
					Spec: v1.PodSpec{
						TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
						ImagePullSecrets: []v1.LocalObjectReference{
							v1.LocalObjectReference{
								Name: "docker-io",
							},
						},
						Containers: []v1.Container{
							v1.Container{
								Name:    service.Name,
								Image:   data.Docker.Image,
								Command: []string{commandEntryPoint},
								Args:    commandArgs,
								Ports:   deployPorts,
								Resources: v1.ResourceRequirements{
									Limits: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("500m"),
										v1.ResourceMemory: resource.MustParse("1Gi"),
									},
									Requests: v1.ResourceList{
										v1.ResourceCPU:    resource.MustParse("300m"),
										v1.ResourceMemory: resource.MustParse("512Mi"),
									},
								},
								ImagePullPolicy: v1.PullIfNotPresent,
								Env:             myEnvVars,
								ReadinessProbe:  &readyProbe,
								LivenessProbe:   &liveProbe,
								VolumeMounts:    volumeMounts,
							},
						},
						Volumes:       deployVolumes,
						RestartPolicy: v1.RestartPolicyAlways,
						DNSPolicy:     v1.DNSClusterFirst,
					},
				},
			},
		}
		log.Printf("Getting list of deployments matching %s", deploymentName)
		_, err := depInterface.Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
		var myError error
		if err != nil {
			// Create deployment if it does not exist
			log.Printf("Existing deployment not found for %s. requested action: %s.", deploymentName, service.Action)
			// Sanity check that we were told to create this service or error out.
			_, myError = depInterface.Deployments(namespace).Create(deployParams)
			if myError != nil {
				// send failed status
				log.Printf("Failed to create service deployment %s, with error: %s", deploymentName, myError)
				data.Services[index].State = plugins.Failed
				data.Services[index].StateMessage = fmt.Sprintf("Error creating deployment: %s", myError)
				// continue so that rollback can happen for all services
				// shorten the timeout in this case so that we can rollback without waiting
				curTime = timeout
			} else {
				// Track rollbackTargets to exclude services that failed to create/updates from the coming rollback.
				rollbackTargets = append(rollbackTargets, service)
			}
		} else {
			// Deployment exists, update deployment with new configuration
			_, myError = depInterface.Deployments(namespace).Update(deployParams)
			if myError != nil {
				log.Printf("Failed to update service deployment %s, with error: %s", deploymentName, myError)
				data.Services[index].State = plugins.Failed
				data.Services[index].StateMessage = fmt.Sprintf("Failed to update deployment %s, with error: %s", deploymentName, myError)
				// continue so that rollback can happen for all services
				// shorten the timeout in this case so that we can rollback without waiting
				curTime = timeout
			} else {
				// Track rollbackTargets to exclude services that failed to create/update from the coming rollback.
				rollbackTargets = append(rollbackTargets, service)
			}
		}
	} // All service deployments initiated.

	log.Printf("Waiting %d seconds for deployment to succeed.", timeout)
	// Set all services initial state to Failed so that we know which have not succeeded.
	for i := range data.Services {
		data.Services[i].State = plugins.Failed
	}

	// Check status of all deployments till the succeed or timeout.
	for {
		for index, service := range data.Services {
			deploymentName := genDeploymentName(data.Project.Slug, service.Name)
			deployment, err := depInterface.Deployments(namespace).Get(deploymentName, metav1.GetOptions{})
			if err != nil {
				log.Printf("Error '%s' fetching deployment status for %s", err, deploymentName)
				continue
			}
			log.Printf("Waiting for %s; ObservedGeneration: %d, Generation: %d, UpdatedReplicas: %d, Replicas: %d, AvailableReplicas: %d, UnavailableReplicas: %d", deploymentName, deployment.Status.ObservedGeneration, deployment.ObjectMeta.Generation, deployment.Status.UpdatedReplicas, *deployment.Spec.Replicas, deployment.Status.AvailableReplicas, deployment.Status.UnavailableReplicas)

			if deployment.Status.ObservedGeneration >= deployment.ObjectMeta.Generation && deployment.Status.UpdatedReplicas == *deployment.Spec.Replicas && deployment.Status.AvailableReplicas >= deployment.Status.UpdatedReplicas && deployment.Status.UnavailableReplicas == 0 {
				// deployment success
				data.Services[index].State = plugins.Complete
				/// AH HA!! haha
				successfulDeploys = 0
				for _, d := range data.Services {
					if d.State == plugins.Complete {
						successfulDeploys++
					}
				}
				log.Printf("%s deploy: %d of %d successful.", deploymentName, successfulDeploys, totalDeploysRequested)
			}
			if successfulDeploys == totalDeploysRequested {
				// all success!
				log.Printf("All deployments successful.")
				// delete the deployments that were requested to be deleted
				//  if anything fails to delete, we need to initiate rollback also.
				allDeleteSuccess = 0
				for di, dService := range data.Services {
					if dService.Action != plugins.Destroy {
						continue
					}
					deleteName := genDeploymentName(data.Project.Slug, dService.Name)
					deleteError := depInterface.Deployments(namespace).Delete(deleteName, &v1.DeleteOptions{})
					if deleteError != nil {
						// If the deletion of a service fails, break out of the loop and continue toward rollback.
						errMessage := fmt.Sprintf("Error when deleting, aborting deploy: %s", deleteError)
						data.Services[di].State = plugins.Failed
						data.Services[di].StateMessage = errMessage
						break
					} else {
						data.Services[di].State = plugins.Complete
						allDeleteSuccess++
					}
				}
				if allDeleteSuccess == totalDeletesRequested {
					// send success message and return.
					x.sendDDSuccessResponse(e, data.Services)
					return nil
				}
			}
		}
		if curTime >= timeout {
			// timeout and get ready to rollback!
			log.Printf("Error, timeout reached waiting for all deployments to succeed.")
			break
		}
		time.Sleep(5 * time.Second)
		curTime += 5
	}

	// Rollback ALL services if anything went wrong.
	successfulRollbacks := 0
	if successfulDeploys != totalDeploysRequested || allDeleteSuccess != totalDeletesRequested {
		log.Printf("Rolling back all deployments for project %s", data.Project.Slug)
		for _, service := range rollbackTargets {
			deploymentName := genDeploymentName(data.Project.Slug, service.Name)
			switch service.Action {
			case plugins.Create:
				// Delete the service
				deleteError := depInterface.Deployments(namespace).Delete(deploymentName, &v1.DeleteOptions{})
				if deleteError != nil {
					log.Printf("Error during rollback when deleting: %s", deleteError)
				} else {
					successfulRollbacks++
				}
			case plugins.Destroy:
				// Rollback the service to the previous version
				deploymentRollback := &v1beta1.DeploymentRollback{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Deployment",
						APIVersion: "extensions/v1beta1",
					},
					Name: deploymentName,
					RollbackTo: v1beta1.RollbackConfig{
						Revision: 0,
					},
				}
				rollbackError := depInterface.Deployments(namespace).Rollback(deploymentRollback)
				if rollbackError != nil {
					log.Printf("Error '%s' rolling back deployment for %s", rollbackError, deploymentName)
				} else {
					successfulRollbacks++
				}
			case plugins.Update:
				// Rollback the service to the previous version
				deploymentRollback := &v1beta1.DeploymentRollback{
					TypeMeta: metav1.TypeMeta{
						Kind:       "Deployment",
						APIVersion: "extensions/v1beta1",
					},
					Name: deploymentName,
					RollbackTo: v1beta1.RollbackConfig{
						Revision: 0,
					},
				}
				rollbackError := depInterface.Deployments(namespace).Rollback(deploymentRollback)
				if rollbackError != nil {
					log.Printf("Error '%s' rolling back deployment for %s", rollbackError, deploymentName)
				} else {
					successfulRollbacks++
				}
			}
		}
		// Send fail message
		x.sendDDErrorResponse(e, data.Services, "Error: One or more deployments failed. Rollback initiated.")
	} // All rollbacks initiated.
	return nil
}
