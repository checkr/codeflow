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
* Finish (all set, ready to login with your Okta login)


