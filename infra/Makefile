FRONTEND_VERSION=$(shell jq -rc ".frontend_version" ../versions.json)
WEBHOOK_VERSION=$(shell jq -rc ".webhook_version" ../versions.json)

deploy: SHELL := /bin/bash
deploy:
	sed -i -Ee "s/(webhook_service_tag = \").*(\")/\1$(WEBHOOK_VERSION)\2/g" terraform.tfvars
	sed -i -Ee "s/(frontend_service_tag = \").*(\")/\1$(FRONTEND_VERSION)\2/g" terraform.tfvars
	tofu apply
