# Running Codeflow with Docker Compose

The docker-compose method of starting Codeflow is to make it easy to develop Codeflow as well as try it out locally.

## Checkout the source

`git clone https://github.com/checkr/codeflow.git`

## Configure basic settings

```
cp server/configs/codeflow.yml server/configs/codeflow.dev.yml
cp dashboard/.env dashboard/.env.development
```

Edit the codeflow.dev.yml and set the following:

* Optional: Configure [Additional Settings](settings.md)

## Build and Start Codeflow
```
# (re)build codeflow docker images
make build
# start up codeflow
make up
```
