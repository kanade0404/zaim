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
  source     = "./modules/service_account"
  id         = "zaim-func"
  roles      = ["secretmanager.secretAccessor", "storage.objectUser"]
  project_id = var.PROJECT_ID
}


module "pubsub" {
  source        = "./modules/pubsub"
  name          = "zaim-trigger"
  subscriptions = [{ name : "zaim-func-trigger" }]
}

module "scheduler" {
  source = "./modules/scheduler"
  name   = "zaim"
  pubsub_target = {
    topic_name = module.pubsub.topic_id
    data = base64encode(jsonencode({
      "users" : keys(var.USERS_SECRET),
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
  location      = "asia-northeast1"
  repository_id = "zaim-api"
  format        = "DOCKER"
}
