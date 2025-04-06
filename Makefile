ifneq (,$(wildcard .env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

GIT_SHA := $(shell git rev-parse --short HEAD)
GIT_REPO := "github.com/thoughtgears/cloudflare-tunnels-poc"

.PHONY: dev lint spec build push deploy

dev:
	go mod tidy
	go run main.go

lint:
	golangci-lint run --timeout 5m
	hadolint Dockerfile

spec:
	rm -rf docs
	swag init

build: lint spec
	docker build --platform linux/amd64 --build-arg SRC_PATH=$(GIT_REPO) -t $(DOCKER_REPO)/$(SERVICE_NAME) .
	docker tag $(DOCKER_REPO)/$(SERVICE_NAME):latest $(DOCKER_REPO)/$(SERVICE_NAME):$(GIT_SHA)

push:
	docker push $(DOCKER_REPO)/$(SERVICE_NAME):latest
	docker push $(DOCKER_REPO)/$(SERVICE_NAME):$(GIT_SHA)

deploy:
	gcloud run deploy $(SERVICE_NAME) \
		--image $(DOCKER_REPO)/$(SERVICE_NAME):$(GIT_SHA) \
		--platform managed \
		--region $(GCP_REGION) \
		--allow-unauthenticated \
		--project $(GCP_PROJECT_ID) \
		--set-env-vars GIT_SHA=$(GIT_SHA) \
		--concurrency 20 \
		--cpu 1 \
		--memory 128Mi