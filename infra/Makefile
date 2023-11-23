FRONTEND_VERSION=$(shell jq -rc ".frontend_version" ../versions.json)
WEBHOOK_VERSION=$(shell jq -rc ".webhook_version" ../versions.json)

replace: SHELL := /bin/bash
replace:
	sed -i -Ee "s/(webhook_service_tag = \").*(\")/\1$(WEBHOOK_VERSION)\2/g" terraform.tfvars
	sed -i -Ee "s/(frontend_service_tag = \").*(\")/\1$(FRONTEND_VERSION)\2/g" terraform.tfvars

deploy: SHELL := /bin/bash
deploy: replace
	tofu apply

deploy-yes: SHELL := /bin/bash
deploy-yes: replace
	tofu apply -auto-approve
