BASE_URL="us-central1-docker.pkg.dev/pbaxter-infra/capstone-repo"


build:
	docker build . -t frontend-service:latest
	docker tag frontend-service:latest $(BASE_URL)/frontend-service:latest

deploy: build
	docker push $(BASE_URL)/frontend-service:latest
