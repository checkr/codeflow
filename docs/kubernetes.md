# Running Codeflow as a Kubernetes Deployment
For the purposes of this document we will be running a kubernetes cluster using [Minikube](https://github.com/kubernetes/minikube).  You may also adapt these instructions to run Codeflow on __any__ kubernetes cluster.
## Codeflow Services Architecture
Codeflow has 4 main services that we will need to setup in Kubernetes.  See below.
```
                       +------------------------+                   +------------------------------+
          +------------+   Client Web Browser   +------+            | Github Webhook Notifications |
          |            |                        |      |            |                              |
          |            +--------+---------------+      |            +------------+-----------------+
          |                     |                      |                         |
          |                     |                      |                         |
+---------v---------+  +--------v------+  +------------v-------+  +--------------v---+
| Dashboard Service |  |  API Service  |  | Websockets Service |  | Webhooks Service |
| React SPA         |  |               |  |                    |  |                  |
|                   |  |               |  |                    |  |                  |
+-------------------+  +----+---+------+  +------------+-------+  +------+-----------+
                            |   |                      |                 |
     +--------------+       |   |                      |                 |
     |              |       |   |              +-------v-----------------v--+
     |  MongoDB     <-------+   |              |                            |
     |              |           +-------------->           Redis            |
     +--------------+                          |                            |
                                               +----------------------------+
```

## Setting up Codeflow Services for access via your preferred method of Load Balancing
* If using Minikube you can skip this step.
* Edit `codeflow/kubernetes/codeflow-services.yaml` to customize the service type to your needs.

## Basic Required Settings

Edit the deployment spec for Codeflow in `codeflow/kubernetes/codeflow-deployment.yaml`.

* __Service URLs__
 * For Minikube our services will be listening on the NodePort of our VM.  For other ingress(s) you must enter the proper URLs for each service.

 * Get the IP of Minikube
```
$ minikube ip
192.168.99.100
```
 * Use this IP address to setup the 4 service URLs for each Container.
```
        - name: REACT_APP_API_ROOT
          value: "http://192.168.99.100:31001"
        - name: REACT_APP_ROOT
          value: "http://192.168.99.100:31004"
        - name: REACT_APP_WEBHOOKS_ROOT
          value: "http://192.168.99.100:31002"
        - name: REACT_APP_WS_ROOT
          value: "ws://192.168.99.100:31003"
```

* Configure [Okta](okta.md) and allow the Minikube Url.

* __Kubernetes credentials__ 
 * These must be named to match the filenames in codeflow-deployment.yaml.  Eg.
```
# Importing credentials from Minikube example
# create a temp directory
mkdir tempdir && cd tempdir
# Copy/Rename the files to the expected names
cp $HOME/.minikube/ca.crt ./ca.pem
cp $HOME/.minikube/apiserver.crt ./admin.pem
cp $HOME/.minikube/apiserver.key ./admin-key.pem
cp $HOME/.kube/config ./kubeconfig
# Load these as secrets into kubernetes
kubectl create secret generic codeflow-kubernetes-secrets --from-file=./ca.pem --from-file=./admin.pem --from-file=./admin-key.pem --from-file=./kubeconfig
```

* __Docker Hub credentials__
```
        - name: CF_PLUGINS_DOCKER_BUILD_REGISTRY_USER_EMAIL
          value: "na@example.com"
        - name: CF_PLUGINS_DOCKER_BUILD_REGISTRY_USERNAME
          value: "naregistry"
        - name: CF_PLUGINS_DOCKER_BUILD_REGISTRY_PASSWORD
          value: "naregistry"
```
* __OKTA Organization Name__
```
       - name: CF_PLUGINS_CODEFLOW_AUTH_OKTA_ORG
          value: "my company org"
```
* __JWT Token__
```
       - name: CF_PLUGINS_CODEFLOW_JWT_SECRET_KEY
          value: "randomizedsecret"
```
* For an in-depth explanation of __all__ Codeflow settings see [Configuration Settings](settings.md) 

# Create the Kubernetes resources.
Then create the kubernetes resources.
```
kubectl create -f codeflow-services.yaml
kubectl create -f codeflow-deployment.yaml
```