# Capstone Project

## Project Goals

1. To meet and fufil the criteria of the WGU Software Engineering Capstone performance assessment.
1. To demonstrate my skills in designing and implementing complex systems
1. Build a service that fixes a problem my peers and I have

## Code Base

I like to keep track of what the codebase looks like

![picture](code.png)

# Requirements

Make sure you have isntalled:
- make
- jq
- GNU sed (If on macos `bew install gnu-sed` and change references to sed -> gsed)
- docker
- [OpenTofu](https://opentofu.org/) _open source terraform_
- go

Optional
- [codevis](https://github.com/sloganking/codevis) - _make the pretty picture_
- [air](https://github.com/cosmtrek/air) - hot reload

# How to run

## Infrastructure

infrastructure is deployed via terraform. It also makes some assumptions about the GCPenvironment its being deployed into.
Those main assumptions that it has the preston-baxter.com hosted zone. You may need to change this to make this terraform template work for You

### Terraform Variables

Contents of `terraform.tfvars`

```toml
project_id = "pbaxter-infra"
project_region = "us-central1"
webhook_service_tag = "0.2.1"
frontend_service_tag = "0.2.1"
```

You only need to specify the `project_id` and `project_region`

Both of the serice tags will be updated automatically via `make deploy`

### Run Terraform

You can either cd into the `infra/` direcrotry and run `make deploy`
or
You can run `make deploy-tf` from the root directory

## Services

The webook service is located in the `service/` directory
The frontend service is located in the `ui/` directory

Both services get ran the same way. What works will one will work for the other

### Config

Sample config
``` yaml
jwt_secret: some_random_string_that_is_long_but_not_too_long
env: test
mongo:
  uri: "mongodb://localhost:27017"
  ent_db: capstoneDB
  ent_col: entities
  lock_db: capstoneDB
  lock_col: locks
app_settings:
  webhook_service_url: localhost:8080
  frontend_service_url: localhost:8080
vendors:
  pco:
    client_id: "test_client_id"
    client_secret: "test_secret"
    scopes:
      - 'people'
      - 'calendar'
      - 'services'
    auth_uri: "https://api.planningcenteronline.com/oauth/authorize"
    token_uri: "https://api.planningcenteronline.com/oauth/token"
    refresh_encode: json
  youtube:
    client_id: "test_client_id"
    client_secret: "test_secret"
    scopes:
      - "https://www.googleapis.com/auth/youtube"
      - "https://www.googleapis.com/auth/youtube.force-ssl"
      - "https://www.googleapis.com/auth/youtube.download"
      - "https://www.googleapis.com/auth/youtube.upload"
    auth_uri: "https://accounts.google.com/o/oauth2/v2/auth"
    token_uri: "https://oauth2.googleapis.com/token"
    refresh_encode: url
  test:
    client_id: "client_id"
    client_secret: "client_secret"
    scopes:
      - "scope 1"
      - "scope 2"
    auth_uri: "server/auth"
    token_uri: "server/token"
    refresh_encode: url
```

Config is expected at `/etc/capstone/config.yaml`

### Run service locally

Both services are configured with [air](https://github.com/cosmtrek/air). Air is a hot reload tool that speeds up the development process. It is particularly useful for working on the frontend service

To run locally, cd into the service directory and run `air`

### Make Docker Container

from the root directory

```bash
make build-service
make build-ui
```

### Run Docker Container Locally

```bash
docker run -p 8080:8080 -it webhook-service:latest
docker run -p 8080:8080 -it frontend-service:latest
```

# Deploy

Make suer versions are updated in `versions.json`

```bash
make deploy
```

Its that easy

NOTE: You may be asked to approve a change set.
