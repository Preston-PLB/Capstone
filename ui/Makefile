BASE_URL="us-central1-docker.pkg.dev/pbaxter-infra/capstone-repo"

local-build:
	rm -f templates/*.go; templ generate -path ./templates
	npx tailwindcss -i static/index.css -o dist/output.css
	GOEXPERIMENT=loopvar go build -o ./tmp/main .

local-run: local-build
	./tmp/main

infra-clean:
	cd docker; docker compose down

infra:
	cd docker; docker compose up -d --remove-orphans


build:
	docker build -f ../docker/ui.dockerfile . -t frontend-service:latest
	docker tag frontend-service:latest $(BASE_URL)/frontend-service:latest

deploy: build
	docker push $(BASE_URL)/frontend-service:latest
