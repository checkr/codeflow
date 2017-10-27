.PHONY:	up server dashboard build

up:
	docker-compose run server go install
	docker-compose run go run main.go --config ./configs/codeflow.dev.yml migrate up
	docker-compose up -d redis mongo
	docker-compose up server dashboard

build:
	docker-compose build server
	docker-compose build dashboard

destroy:
	docker-compose stop
	docker-compose rm -f

test:
	cd dashboard && npm i && npm run lint
