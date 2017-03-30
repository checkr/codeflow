# Codeflow Okta
Codeflow uses Okta for SSO logins.

## Setup new Okta application for Codeflow
* Login to Okta with administrator rights.
* Go to Applications -> Add Application
* Click 'Create New App'

## Creating the Codeflow App
* Select SPA (Single Page App)
* Redirect URIs:
```
http://localhost:3000
http://localhost:3001/oauth2/callback/okta
https://codeflow.mydomain.com
https://codeflow-api.mydomain.com/oauth2/callback/okta
``` 

## Modify Codeflow settings for OKTA
* Using Kubernetes guide
  * Edit `codeflow/kubernetes/codeflow-config.env` and add:
```
CF_AUTH_HANDLER="okta"
CF_OKTA_ORG="myCompany"
```
  * Parse/Load the secrets
```
cat codeflow/kubernetes/codeflow-config.env | kubernetes-secret -n codeflow-config > codeflow-config.yaml
kubectl create -f codeflow-config.yaml
```

* Using docker-compose guide
  * Edit `server/configs/codeflow.dev.yml` under the plugins/codeflow section:
```
    auth:
      path: "/auth"
      handler: "okta"
      # This is important to set for Okta SSO
      okta_org: "MyCompanyOrg"
```
  * Re-launch Codeflow
```
make up
```