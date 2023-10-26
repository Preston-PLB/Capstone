BASE_URL="us-central1-docker.pkg.dev/pbaxter-infra/capstone-repo"


build:
	docker build . -t webhook-service:latest
	docker tag webhook-service:latest $(BASE_URL)/webhook-service:latest

deploy: build
	docker push $(BASE_URL)/webhook-service:latest
	
