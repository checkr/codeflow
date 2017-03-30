# Codeflow

Extendable deployment pipeline

# Local Development
Install [Docker Compose](https://docs.docker.com/compose/install/)

Copy `server/codeflow.yml` to `server/codeflow.dev.yml`

To start the server and all dependencies run:
```
$ make up
```
Check `Makefile` to see what's happening under the hood.

Local dashboard will be started on [http://localhost:3000](http://localhost:3000)

`dashboard` and `server` will reload on any file change :boom: Happy coding!!!

