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

resource "google_project_service" "serverless_vpc_api" {
  service = "vpcaccess.googleapis.com" 

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

resource "google_cloud_run_v2_service" "webhook_service_cr" {
  name = "webhook-service-cr" 
  location = var.project_region
  launch_stage = "BETA"
  
  template {
      containers {
        image = "${var.project_region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.capstone_repo.name}/webhook-service:latest"
      }
      vpc_access{
        network_interfaces {
          network = "main"
        subnetwork = "application"
        }
        egress = "ALL_TRAFFIC"
      }
  }
  depends_on = [ google_project_service.run_api, google_artifact_registry_repository.capstone_repo ]
}


resource "google_cloud_run_v2_service_iam_member" "webhook_service_run_all_users" {
  project = var.project_id
  name = google_cloud_run_v2_service.webhook_service_cr.name
  location = var.project_region
  role = "roles/run.invoker"
  member = "allUsers"
}

resource "google_cloud_run_v2_service" "frontend_service_cr" {
  name = "frontend-service-cr" 
  location = var.project_region
  launch_stage = "BETA"
  
  template {
    containers {
      image = "${var.project_region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.capstone_repo.name}/frontend-service:latest"
    }
    vpc_access{
      network_interfaces {
        network = "main"
        subnetwork = "application"
      }
      egress = "ALL_TRAFFIC"
    }
  }

  depends_on = [ google_project_service.run_api, google_artifact_registry_repository.capstone_repo ]
}


resource "google_cloud_run_v2_service_iam_member" "frontend_service_run_all_users" {
  project = var.project_id
  name = google_cloud_run_v2_service.frontend_service_cr.name
  location = var.project_region
  role = "roles/run.invoker"
  member = "allUsers"
}

data "google_dns_managed_zone" "preston_baxter_zone" {
  name = "pbaxter-main-zone"
}

resource "google_dns_record_set" "webhook_cname" {
  name         = "webhook.${data.google_dns_managed_zone.preston_baxter_zone.dns_name}"
  managed_zone = data.google_dns_managed_zone.preston_baxter_zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas = [
    "ghs.googlehosted.com."
  ]

  depends_on = [ google_cloud_run_v2_service.webhook_service_cr ]
}

resource "google_dns_record_set" "frontend_cname" {
  name         = "frontend.${data.google_dns_managed_zone.preston_baxter_zone.dns_name}"
  managed_zone = data.google_dns_managed_zone.preston_baxter_zone.name
  type         = "CNAME"
  ttl          = 300
  rrdatas = [
    "ghs.googlehosted.com."
  ]

  depends_on = [ google_cloud_run_v2_service.frontend_service_cr ]
}

resource "google_cloud_run_domain_mapping" "frontend_cname_mapping" {
  location = "us-central1"
  name     = trimsuffix("frontend.${data.google_dns_managed_zone.preston_baxter_zone.dns_name}", ".")

  metadata {
    namespace = var.project_id 
  }

  spec {
    route_name = google_cloud_run_v2_service.frontend_service_cr.name
  }
}

resource "google_cloud_run_domain_mapping" "webhook_cname_mapping" {
  location = "us-central1"
  name     = trimsuffix("webhook.${data.google_dns_managed_zone.preston_baxter_zone.dns_name}", ".")

  metadata {
    namespace = var.project_id 
  }
    

  spec {
    route_name = google_cloud_run_v2_service.webhook_service_cr.name
  }
}
