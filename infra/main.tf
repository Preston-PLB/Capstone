terraform {
  required_version = ">= 0.14"

  required_providers {
    google = ">= 3.3"
  }
}

variable "project_id" {
  description = "The GCP project ID where the infra will be built"
  type = string
}

variable "project_region" {
  description = "The GCP region where the infra will be built"
  type = string
}

provider "google" {
  project = var.project_id
}

resource "google_project_service" "run_api" {
  service = "run.googleapis.com"

  disable_on_destroy = true
}

resource "google_project_service" "artifact_api" {
  service = "artifactregistry.googleapis.com" 

  disable_on_destroy = true
}


resource "google_artifact_registry_repository" "capstone_repo" {
  location = var.project_region
  repository_id = "capstone-repo"
  description = "Images for capstone project"
  format = "DOCKER"
 
  docker_config {
    immutable_tags = false
  }

  depends_on = [ google_project_service.artifact_api ]
}

resource "google_cloud_run_service" "webhook_service_cr" {
  name = "webhook-service-cr" 
  location = var.project_region
  
  template {
    spec {
      containers {
        image = "${var.project_region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.capstone_repo.name}/webhook-service:latest"
      }
    }
  }

  traffic {
    percent = 100
    latest_revision = true
  }

  depends_on = [ google_project_service.run_api, google_artifact_registry_repository.capstone_repo ]
}

resource "google_cloud_run_service_iam_member" "webhook_service_run_all_users" {
  service = google_cloud_run_service.webhook_service_cr.name
  location = var.project_region
  role = "roles/run.invoker"
  member = "allUsers"
}

resource "google_cloud_run_service" "frontend_service_cr" {
  name = "frontend-service-cr" 
  location = var.project_region
  
  template {
    spec {
      containers {
        image = "${var.project_region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.capstone_repo.name}/frontend-service:latest"
      }
    }
  }

  traffic {
    percent = 100
    latest_revision = true
  }

  depends_on = [ google_project_service.run_api, google_artifact_registry_repository.capstone_repo ]
}

resource "google_cloud_run_service_iam_member" "frontend_service_run_all_users" {
  service = google_cloud_run_service.frontend_service_cr.name
  location = var.project_region
  role = "roles/run.invoker"
  member = "allUsers"
}
