.PHONY:	up build destroy

up:
	docker-compose run --rm circuit go run main.go --config ./configs/circuit.dev.yml migrate up
	docker-compose up -d redis mongo

build:
	docker-compose build server

destroy:
	docker-compose stop
	docker-compose rm -f
