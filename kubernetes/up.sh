#!/bin/bash -e

echo This script helps configure codeflow for the initial deployment into a K8S environment.
echo REQUIRES: jq
echo
echo Additional settings can be configured prior to running this script by editing the files:
echo "* codeflow-services.yaml (optional annotations for your service like SSL)"
echo "* ../server/configs/codeflow.dev.yml (optional)"
echo 
if [ -z "$NONINTERACTIVE" ]; then
	read -p "Continue with the current settings? (y/n)" yn
	if [ "$yn" != "y" ]; then
		echo Aborting..
		exit 1
	fi
fi

# This requires jq to be installed
if [ ! -x "$(command -v jq)" ]; then
	echo up.sh requires jq to be installed.
	echo OSX: brew install jq
	echo Linux: apt-get
	exit 1
fi

# Sanity check, the config exists
if [ ! -e ../server/configs/codeflow.dev.yml ]; then
	echo file not found: ../server/configs/codeflow.dev.yml. Please cp codeflow.yml to codeflow.dev.yml and re-run.
	exit 1
fi

get_ingress_hostname () {
	ingress_hostname=$(kubectl get services --namespace=development-checkr-codeflow -ojson |jq -r ".items[] | select(.metadata.name==\"${1}\") | .status.loadBalancer.ingress[0].hostname")
}

wait_for_ingress_hostname () {
	get_ingress_hostname codeflow-api
	echo waiting for hostname for codeflow-api ...
	until [ -n "$ingress_hostname" ] && [ "$ingress_hostname" != "null" ]; do
		sleep 5
		get_ingress_hostname codeflow-api
	done

	get_ingress_hostname codeflow-dashboard
	echo waiting for hostname for codeflow-dashboard ...
	until [ -n "$ingress_hostname" ] && [ "$ingress_hostname" != "null" ]; do
		sleep 5
		get_ingress_hostname codeflow-dashboard
	done

}
 
# get_url 'servicename' 'scheme'
get_url () {
	url=${3}://$(kubectl get services --namespace=development-checkr-codeflow -ojson |jq -r ".items[] | select(.metadata.name==\"${2}\") | [ .status.loadBalancer.ingress[0].hostname, (.spec.ports[] | select(.name==\"${1}\") | .port |tostring) ] |join(\":\")")
}

get_dashboard_port () {
	port=$(kubectl get services --namespace=development-checkr-codeflow -ojson |jq ".items[] | select(.metadata.name==\"codeflow-dashboard\") | .spec.ports[] | select(.name==\"dashboard-port\") | .targetPort |tostring")
}

detect_ssl () {
	ssl_arn=$(kubectl get services --namespace=development-checkr-codeflow -ojson |jq -r ".items[] | select(.metadata.name==\"codeflow-dashboard\") |.metadata.annotations.\"service.beta.kubernetes.io\/aws-load-balancer-ssl-cert\"")
	if [ -n "$ssl_arn" ] && [ "$ssl_arn" != "null" ]; then
		echo Using TCP+SSL protocol..
		protocol=s
	fi
}

set +e
echo creating namespaces 
kubectl create namespace development-checkr-codeflow

echo creating mongodb and redis
kubectl create -f mongodb-service.yaml
kubectl create -f mongodb-deployment.yaml
kubectl create -f redis-service.yaml
kubectl create -f redis-deployment.yaml

echo creating codeflow services
kubectl create -f codeflow-services.yaml
set -e

echo configuring codeflow dashboard
envfile=react-configmap.yaml
cat << EOF > $envfile
# This ConfigMap is used to configure codeflow react service (dashboard).
kind: ConfigMap
apiVersion: v1
metadata:
  name: react-config
  namespace: development-checkr-codeflow
data:
EOF

detect_ssl

wait_for_ingress_hostname

get_url 'api-port' 'codeflow-api' "http${protocol}"
echo "  REACT_APP_API_ROOT: $url" >> $envfile

get_url 'webhooks-port' 'codeflow-api' "http${protocol}"
echo "  REACT_APP_WEBHOOKS_ROOT: $url" >> $envfile

get_url 'websockets-port' 'codeflow-api' "ws${protocol}"
echo "  REACT_APP_WS_ROOT: $url" >> $envfile

get_url 'dashboard-port' 'codeflow-dashboard' "http${protocol}"
echo "  REACT_APP_ROOT: $url"  >> $envfile

get_dashboard_port
echo "  REACT_APP_PORT: $port"  >> $envfile

kubectl apply -f $envfile --namespace=development-checkr-codeflow 

echo react-configmap generated and applied from file: $envfile
echo Services configured successfully..
echo
echo Dashboard URL:  $url

echo
echo configuring codeflow api
kubectl create configmap codeflow-config --from-file=../server/configs/codeflow.dev.yml --namespace=development-checkr-codeflow

echo running codeflow database migration job
kubectl create -f codeflow-migration-job.yaml

echo creating codeflow deployment
kubectl create -f codeflow-deployment.yaml
