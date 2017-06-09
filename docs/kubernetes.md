# Bootstrapping Codeflow as a Kubernetes Deployment
This guide assumes you have a kubernetes cluster running in AWS.  The kubernetes config files in codeflow/kubernetes could be adapted to any kubernetes cluster.

## Notes
Once you have bootstrapped using this guide you will be able to use Codeflow to deploy itself and change it's own settings.  That means for this initial bootstrap you don't need to configure any settings beyond the stock settings unless you want to.

## Service settings (optional) 
* Edit `kubernetes/codeflow-services.yaml` to customize the service types to your needs.

## Codeflow settings
* Copy the stock codeflow.yml config into a new file for local modification.
```
cp server/configs/codeflow.yml server/configs/codeflow.dev.yml
```
* Edit `server/configs/codeflow.dev.yml` (optional)

## Create Kubernetes resources

### Automatic bootstrap
To assist with configuring the codeflow URLs in the react-config we've provided a helper script that uses jq to populate the react-configmap.
* __Requires__: jq
```
cd kubernetes
./up.sh
```

### Manual bootstrap
```
cd kubernetes
# MongoDB
kubectl create -f mongodb-deployment.yaml
kubectl create -f mongodb-service.yaml

# Redis
kubectl create -f redis-deployment.yaml
kubectl create -f redis-service.yaml

# Codeflow Services (you will need the URLs to configure react below)
kubectl create -f codeflow-services.yaml

# First edit the react-configmap.yaml to match the codeflow service URLs
# Then create the configmap
kubectl create -f react-configmap.yaml

# Load the codeflow configmap
kubectl create configmap codeflow-config --from-file=../server/configs/codeflow.dev.yml --namespace=development-checkr-codeflow

# Codeflow
kubectl create -f codeflow-deployment.yaml
```