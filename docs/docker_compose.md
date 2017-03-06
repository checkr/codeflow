# Running Codeflow with Docker Compose

The docker-compose method of starting Codeflow is to make it easy to develop Codeflow as well as try it out locally.

## Checkout the source

`git clone https://github.com/checkr/codeflow.git`

## Configure basic settings

* Configure [Okta](okta.md)

```
cp server/config/codeflow.yml server/config/codeflow.dev.yml
```

Edit the codeflow.dev.yml and set the following:

```
jwt_secret_key: "some random value"
okta_org: "your org name"
```

* Configure [Additional Settings](settings.md)

## Build and Start Codeflow
```
docker-compose up
```