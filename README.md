# THIS PROJECT IS NOT MAINTAINED ANYMORE! All further development has been moved to CodeAmp https://github.com/codeamp

# Codeflow [![CircleCI](https://circleci.com/gh/checkr/codeflow/tree/master.svg?style=svg)](https://circleci.com/gh/checkr/codeflow/tree/master)

Extendable deployment pipeline

# Local Development with Docker
Install [Docker Compose](https://docs.docker.com/compose/install/)

### Create DEV configs
```
$ cp server/configs/codeflow.yml server/configs/codeflow.dev.yml
$ cp dashboard/.env dashboard/.env.development
$ cd dashboard/ && npm install && cd ../
```

### To start the server and all dependencies run:
```
$ make up
```

Check `Makefile` to see what's happening under the hood.

Local dashboard will be started on [http://localhost:3000](http://localhost:3000)

`dashboard` and `server` will reload on any file change :boom: Happy coding!!!

### Hosted docs
[https://codeflow.checkr.com/](https://codeflow.checkr.com/)

`master` branch continuously deployed with Codeflow!

## Slack
* [Signup for codeflow-team on Slack](http://codeflow-slack.checkr.com/)

### Screenshots
![](/docs/images/codeflow1.png)
![](/docs/images/codeflow2.png)

