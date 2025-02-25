terraform {
  backend "gcs" {
    bucket = "tfstate-kanade0404-070dc2e4-a61e-e22d-2010-d15e23acf81d"
    prefix = "zaim"
  }
  required_version = "1.9.6"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "6.22.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.7.1"
    }
  }
}

provider "google" {
  project = var.PROJECT_ID
  region  = local.region
  zone    = "${local.region}-a"
  default_labels = {
    service = "zaim"
  }
}

import {
  id = "${local.region}/zaim-api"
  to = google_cloud_run_v2_service.app
}

resource "random_uuid" "uuid" {}

module "zaim-consumer-key" {
  for_each = var.USERS_SECRET
  source   = "./modules/secret_manager"
  id       = "zaim-consumer-key-${each.key}"
  data     = each.value.ZAIM.CONSUMER_KEY
}

module "zaim-consumer-secret" {
  for_each = var.USERS_SECRET
  source   = "./modules/secret_manager"
  id       = "zaim-consumer-secret-${each.key}"
  data     = each.value.ZAIM.CONSUMER_SECRET
}

module "zaim-oauth-token" {
  for_each = var.USERS_SECRET
  source   = "./modules/secret_manager"
  id       = "zaim-oauth-token-${each.key}"
  data     = each.value.ZAIM.OAUTH_TOKEN
}

module "zaim-oauth-secret" {
  for_each = var.USERS_SECRET
  source   = "./modules/secret_manager"
  id       = "zaim-oauth-secret-${each.key}"
  data     = each.value.ZAIM.OAUTH_SECRET
}

module "zaim-csv-folder" {
  for_each = var.USERS_SECRET
  source   = "./modules/secret_manager"
  id       = "csv-folder-${each.key}"
  data     = each.value.CSV_FOLDER
}

module "zaim-func" {
  source        = "./modules/service_account"
  id            = "zaim-func"
  project_roles = ["secretmanager.secretAccessor", "storage.objectUser"]
  project_id    = var.PROJECT_ID
}

module "database-url" {
  source = "./modules/secret_manager"
  id     = "database-url"
  data   = var.DATABASE_URL
}

module "zaim-file" {
  source   = "./modules/gcs"
  location = "ASIA-NORTHEAST1"
  name     = "zaim-${var.PROJECT_ID}-${random_uuid.uuid.id}"
}

resource "google_artifact_registry_repository" "repo" {
  location      = local.region
  repository_id = "zaim-api"
  format        = "DOCKER"
}

resource "google_cloud_run_v2_service" "app" {
  location = local.region
  name     = "zaim-api"
  ingress = "INGRESS_TRAFFIC_ALL"
  template {
    service_account = module.zaim-func.email
    containers {
      name = "app-1"
      image = "asia-northeast1-docker.pkg.dev/${var.PROJECT_ID}/${google_artifact_registry_repository.repo.name}/app:latest"
      ports {
        container_port = 8888
        name = "http1"
      }
      dynamic "env" {
        for_each = { "HOST" : var.RUN_HOST, "PROJECT_ID" : var.PROJECT_ID, "BUCKET_NAME" : module.zaim-file.name, "ENV": "prd" }
        content {
          name  = env.key
          value = env.value
        }
      }
      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret = module.database-url.secret_id
            version = module.database-url.version
          }
        }
      }
    }
  }
  traffic {
    percent         = 100
    type = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
  }
  lifecycle {
    ignore_changes = [
      client,
      client_version,
    ]
  }
  depends_on = [module.zaim-func, google_artifact_registry_repository.repo, module.database-url]
}
