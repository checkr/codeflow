.PHONY:	up server dashboard build

up:
	docker-compose run --rm server /bin/sh -c 'cd server/ && go run main.go --config ./configs/codeflow.dev.yml migrate up'
	docker-compose up -d redis mongo
	docker-compose up server dashboard

build:
	docker-compose build server
	docker-compose build dashboard

dashboard:
	cd ./dashboard && npm run start

server:
	cd ./server && reflex -c reflex.conf
