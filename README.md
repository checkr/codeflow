# Codeflow

Extendable deployment pipeline

# Local Development
Install [Docker Compose](https://docs.docker.com/compose/install/)

Copy `server/codeflow.yml` to `server/codeflow.dev.yml`

You should set `jwt_secret_key` and `okta_org` before starting the server!
```yaml
codeflow:
  jwt_secret_key:
  allowed_origins:
  auth:
    okta_org:
```

To start the server and all dependencies run:
```
$ make
```

Local dashboard will be started on [http://localhost:3000](http://localhost:3000)

`dashboard` and `server` will reload on any file change :boom: Happy coding!!!

