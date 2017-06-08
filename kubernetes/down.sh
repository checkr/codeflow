#!/bin/bash

read -p "Do you want to delete all codeflow services? " yn

if [ "$yn" != "y" ]; then
	echo aborting..
	exit 0
fi

set -x

kubectl delete -f mongodb-service.yaml
kubectl delete -f mongodb-deployment.yaml
kubectl delete -f redis-service.yaml
kubectl delete -f redis-deployment.yaml
kubectl delete -f codeflow-services.yaml
kubectl delete configmap codeflow-config --namespace=development-checkr-codeflow
kubectl delete configmap react-config --namespace=development-checkr-codeflow
kubectl delete -f codeflow-deployment.yaml
kubectl delete -f codeflow-migration-job.yaml
kubectl delete namespace development-checkr-codeflow
