terraform {
  backend "gcs" {
    bucket = "tfstate-kanade0404-070dc2e4-a61e-e22d-2010-d15e23acf81d"
    prefix = "zaim"
  }
  required_version = "1.6.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "5.7.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.5.1"
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

module "pubsub-user" {
  source        = "./modules/service_account"
  id            = "zaim-pubsub"
  project_roles = ["run.invoker"]
  service_roles = ["iam.serviceAccountTokenCreator"]
  project_id    = var.PROJECT_ID
}


module "pubsub" {
  source        = "./modules/pubsub"
  name          = "zaim-trigger"
  subscriptions = [{ name : "zaim-func-trigger", push : { endpoint : google_cloud_run_service.app.status[0].url, service_account_email : module.pubsub-user.email } }]
  depends_on    = [module.pubsub-user]
}

module "scheduler" {
  source = "./modules/scheduler"
  name   = "zaim"
  pubsub_target = {
    topic_name = module.pubsub.topic_id
    data = base64encode(jsonencode({
      "users" : keys(var.USERS_SECRET),
      "dry_run" : false,
    }))
  }
  schedule = "0 0 2 * *"
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

resource "google_cloud_run_service" "app" {
  location = local.region
  name     = "zaim-api"
  template {
    spec {
      service_account_name = module.zaim-func.email
      containers {
        image = "asia-northeast1-docker.pkg.dev/${var.PROJECT_ID}/${google_artifact_registry_repository.repo.name}/app"
        ports {
          container_port = 8888
          name           = "http1"
        }
        dynamic "env" {
          for_each = { "HOST" : var.RUN_HOST }
          content {
            name  = env.key
            value = env.value
          }
        }
      }
    }
  }
  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "internal"
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
  autogenerate_revision_name = true
  lifecycle {
    ignore_changes = [
      metadata[0].annotations["run.googleapis.com/operation-id"],
      metadata[0].annotations["client.knative.dev/user-image"],
      metadata[0].annotations["run.googleapis.com/client-name"],
      metadata[0].annotations["run.googleapis.com/client-version"],
      metadata[0].annotations["serving.knative.dev/creator"],
      metadata[0].annotations["serving.knative.dev/lastModifier"],
      metadata[0].annotations["run.googleapis.com/ingress-status"],
      metadata[0].annotations["run.googleapis.com/launch-stage"],
      metadata[0].labels["cloud.googleapis.com/location"],
      status[0].latest_created_revision_name,
      status[0].latest_ready_revision_name,
      status[0].observed_generation,
      template[0].metadata[0].annotations["client.knative.dev/user-image"],
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["run.googleapis.com/client-version"],
      template[0].metadata[0].annotations["run.googleapis.com/sandbox"],
    ]
  }
  depends_on = [module.zaim-func, google_artifact_registry_repository.repo]
}
