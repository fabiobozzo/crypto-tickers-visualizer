.DEFAULT_GOAL := help

export GOPATH := $(shell go env GOPATH)

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

create-images: 
	docker-compose build --parallel

start: 
	docker-compose up -d --build

stop: 
	docker-compose down -v
