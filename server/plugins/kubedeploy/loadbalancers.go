package kubedeploy

import (
	"fmt"
	"log"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"

	"strings"

	"github.com/checkr/codeflow/server/agent"
	"github.com/checkr/codeflow/server/plugins"
	"github.com/spf13/viper"
)

func (x *KubeDeploy) sendLBResponse(e agent.Event, state plugins.State, failureMessage string, dns string, lbpayload plugins.LoadBalancer) {
	payload := e.Payload.(plugins.LoadBalancer)
	payload.Action = plugins.Status
	payload.Service = lbpayload.Service
	payload.StateMessage = failureMessage
	payload.DNS = dns
	payload.Subdomain = lbpayload.Subdomain
	payload.State = state
	event := e.NewEvent(payload, nil)
	x.events <- event
}

func (x *KubeDeploy) doDeleteLoadBalancer(e agent.Event) error {
	// Codeflow will load the kube config from a file, specified by CF_PLUGINS_KUBEDEPLOY_KUBECONFIG environment variable
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("plugins.kubedeploy.kubeconfig"))

	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Error getting cluster config.")
	}

	payload := e.Payload.(plugins.LoadBalancer)
	coreInterface := clientset.Core()

	namespace := genNamespaceName(payload.Environment, payload.Project.Slug)

	_, svcGetErr := coreInterface.Services(namespace).Get(payload.Name, meta_v1.GetOptions{})
	if svcGetErr == nil {
		// Service was found, ready to delete
		svcDeleteErr := coreInterface.Services(namespace).Delete(payload.Name, &meta_v1.DeleteOptions{})
		if svcDeleteErr != nil {
			failMessage := fmt.Sprintf("Error '%s' deleting service %s", svcDeleteErr, payload.Name)
			log.Printf("ERROR managing loadbalancer %s: %s", payload.Service.Name, failMessage)
			x.sendLBResponse(e, plugins.Failed, failMessage, "", payload)
			return nil
		}
		x.sendLBResponse(e, plugins.Deleted, "", "", payload)
	} else {
		// Send failure message that we couldn't find the service to delete
		failMessage := fmt.Sprintf("Error finding %s service: '%s'", payload.Name, svcGetErr)
		log.Printf("ERROR managing loadbalancer %s: %s", payload.Service.Name, failMessage)
		x.sendLBResponse(e, plugins.Failed, failMessage, "", payload)
	}
	return nil
}

// Make changes to kubernetes services (aka load balancers)
func (x *KubeDeploy) doLoadBalancer(e agent.Event) error {
	payload := e.Payload.(plugins.LoadBalancer)
	// Codeflow will load the kube config from a file, specified by CF_PLUGINS_KUBEDEPLOY_KUBECONFIG environment variable
	config, err := clientcmd.BuildConfigFromFlags("", viper.GetString("plugins.kubedeploy.kubeconfig"))

	if err != nil {
		failMessage := fmt.Sprintf("ERROR: %s; you must set the environment variable CF_PLUGINS_KUBEDEPLOY_KUBECONFIG=/path/to/kubeconfig", err.Error())
		log.Println(failMessage)
		x.sendLBResponse(e, plugins.Failed, failMessage, "", payload)
		return nil
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		failMessage := fmt.Sprintf("ERROR: %s; setting NewForConfig in doLoadBalancer", err.Error())
		log.Println(failMessage)
		x.sendLBResponse(e, plugins.Failed, failMessage, "", payload)
		return nil
	}

	coreInterface := clientset.Core()
	deploymentName := genDeploymentName(payload.Project.Slug, payload.Service.Name)

	var serviceType v1.ServiceType
	var servicePorts []v1.ServicePort
	serviceAnnotations := make(map[string]string)

	namespace := genNamespaceName(payload.Environment, payload.Project.Slug)
	createNamespaceErr := x.createNamespaceIfNotExists(namespace, coreInterface)
	if createNamespaceErr != nil {
		x.sendLBResponse(e, plugins.Failed, createNamespaceErr.Error(), "", payload)
		return nil
	}

	// Begin create
	switch payload.Type {
	case plugins.Internal:
		serviceType = v1.ServiceTypeClusterIP
	case plugins.External:
		serviceType = v1.ServiceTypeLoadBalancer
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"] = "true"
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-timeout"] = "300"
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"] = "true"
		if viper.IsSet("plugins.kubedeploy.access_log_s3_bucket") {
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-emit-interval"] = "5"
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-enabled"] = "true"
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-s3-bucket-name"] = viper.GetString("plugins.kubedeploy.access_log_s3_bucket")
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-s3-bucket-prefix"] = fmt.Sprintf("%s/%s", payload.Project.Slug, payload.Service.Name)
		}
	case plugins.Office:
		serviceType = v1.ServiceTypeLoadBalancer
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-internal"] = "0.0.0.0/0"
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-enabled"] = "true"
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-connection-draining-timeout"] = "300"
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled"] = "true"
		if viper.IsSet("plugins.kubedeploy.access_log_s3_bucket") {
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-emit-interval"] = "5"
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-enabled"] = "true"
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-s3-bucket-name"] = viper.GetString("plugins.kubedeploy.access_log_s3_bucket")
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-access-log-s3-bucket-prefix"] = fmt.Sprintf("%s/%s", payload.Project.Slug, payload.Service.Name)
		}
	}
	var sslPorts []string
	for _, p := range payload.ListenerPairs {
		var realProto string
		switch p.Destination.Protocol {
		case "HTTPS":
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "http"
			realProto = "TCP"
		case "SSL":
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "tcp"
			realProto = "TCP"
		case "HTTP":
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "http"
			realProto = "TCP"
		case "TCP":
			serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-backend-protocol"] = "tcp"
			realProto = "TCP"
		case "UDP":
			realProto = "UDP"
		}
		convPort := intstr.IntOrString{
			IntVal: p.Destination.Port,
		}
		newPort := v1.ServicePort{
			Name:       fmt.Sprintf("%d-%s", p.Source.Port, strings.ToLower(realProto)),
			Port:       p.Source.Port,
			TargetPort: convPort,
			Protocol:   v1.Protocol(realProto),
		}
		if p.Destination.Protocol == "HTTPS" || p.Destination.Protocol == "SSL" {
			sslPorts = append(sslPorts, fmt.Sprintf("%d", p.Source.Port))
		}
		servicePorts = append(servicePorts, newPort)
	}
	if len(sslPorts) > 0 {
		sslPortsCombined := strings.Join(sslPorts, ",")
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-ssl-ports"] = sslPortsCombined
		serviceAnnotations["service.beta.kubernetes.io/aws-load-balancer-ssl-cert"] = viper.GetString("plugins.kubedeploy.ssl_cert_arn")
	}
	serviceSpec := v1.ServiceSpec{
		Selector: map[string]string{"app": deploymentName},
		Type:     serviceType,
		Ports:    servicePorts,
	}
	serviceParams := v1.Service{
		TypeMeta: meta_v1.TypeMeta{
			Kind:       "service",
			APIVersion: "v1",
		},
		ObjectMeta: meta_v1.ObjectMeta{
			Name:        payload.Name,
			Annotations: serviceAnnotations,
		},
		Spec: serviceSpec,
	}

	// Implement service update-or-create semantics.
	service := coreInterface.Services(namespace)
	svc, err := service.Get(payload.Name, meta_v1.GetOptions{})
	switch {
	case err == nil:
		serviceParams.ObjectMeta.ResourceVersion = svc.ObjectMeta.ResourceVersion
		serviceParams.Spec.ClusterIP = svc.Spec.ClusterIP
		_, err = service.Update(&serviceParams)
		if err != nil {
			x.sendLBResponse(e, plugins.Failed, fmt.Sprintf("Error: failed to update service: %s", err.Error()), "", payload)
			return nil
		}
		log.Printf("Service updated: %s", payload.Name)
	case errors.IsNotFound(err):
		_, err = service.Create(&serviceParams)
		if err != nil {
			x.sendLBResponse(e, plugins.Failed, fmt.Sprintf("Error: failed to create service: %s", err.Error()), "", payload)
			return nil
		}
		log.Printf("Service created: %s", payload.Name)
	default:
		x.sendLBResponse(e, plugins.Failed, fmt.Sprintf("Unexpected error: %s", err.Error()), "", payload)
		return nil
	}

	// If ELB grab the DNS name for the response
	var ELBDNS string
	if payload.Type == plugins.External || payload.Type == plugins.Office {
		log.Printf("Waiting for ELB address for %s", payload.Name)
		// Timeout waiting for ELB DNS name after 600 seconds
		timeout := 600
		for {
			elbResult, elbErr := coreInterface.Services(namespace).Get(payload.Name, meta_v1.GetOptions{})
			if elbErr != nil {
				log.Printf("Error '%s' describing service %s", elbErr, payload.Name)
			} else {
				ingressList := elbResult.Status.LoadBalancer.Ingress
				if len(ingressList) > 0 {
					ELBDNS = ingressList[0].Hostname
					break
				}
				if timeout <= 0 {
					failMessage := fmt.Sprintf("Error: timeout waiting for ELB DNS name for: %s", payload.Name)
					x.sendLBResponse(e, plugins.Failed, failMessage, "", payload)
					return nil
				}
			}
			time.Sleep(time.Second * 5)
			timeout -= 5
		}
	} else {
		ELBDNS = fmt.Sprintf("%s.%s", payload.Name, genNamespaceName(payload.Environment, payload.Project.Slug))
	}
	x.sendLBResponse(e, plugins.Complete, "", ELBDNS, payload)

	return nil
}
