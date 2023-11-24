BASE_URL=us-central1-docker.pkg.dev/pbaxter-infra/capstone-repo
FRONTEND_VERSION=$(shell jq -rc ".frontend_version" versions.json)
WEBHOOK_VERSION=$(shell jq -rc ".webhook_version" versions.json)

deploy: deploy-ui deploy-service deploy-tf

deploy-tf:
	cd infra; make deploy

deploy-ui: build-ui
	docker push $(BASE_URL)/frontend-service:latest
	docker push $(BASE_URL)/frontend-service:$(FRONTEND_VERSION)

deploy-service: build-service
	docker push $(BASE_URL)/webhook-service:latest
	docker push $(BASE_URL)/webhook-service:$(WEBHOOK_VERSION)

build: build-ui build-service

build-ui:
	docker build -f ./docker/ui.dockerfile . -t frontend-service:latest
	docker build -f ./docker/ui.dockerfile . -t frontend-service:$(FRONTEND_VERSION)
	docker tag frontend-service:latest $(BASE_URL)/frontend-service:latest
	docker tag frontend-service:$(FRONTEND_VERSION) $(BASE_URL)/frontend-service:$(FRONTEND_VERSION)

build-service:
	docker build -f ./docker/service.dockerfile . -t webhook-service:latest
	docker build -f ./docker/service.dockerfile . -t webhook-service:$(WEBHOOK_VERSION)
	docker tag webhook-service:latest $(BASE_URL)/webhook-service:latest
	docker tag webhook-service:$(WEBHOOK_VERSION) $(BASE_URL)/webhook-service:$(WEBHOOK_VERSION)

image: SHELL := /bin/bash
image:
	[[ -d "/tmp/capstone" ]] || mkdir /tmp/capstone
	cp -R infra/ /tmp/capstone/
	cp -R service/ /tmp/capstone/
	cp -R ui/ /tmp/capstone/
	rm -rf /tmp/capstone/ui/templates/*_templ.go
	codevis -i /tmp/capstone --whitelist-extension go,hcl,tf,templ -o ./out.png
	rm -rf /tmp/capstone/*
