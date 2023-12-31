variable "PROJECT_ID" {
  type = string
}

variable "USERS_SECRET" {
  type = map(object({
    ZAIM = object({
      CONSUMER_KEY    = string
      CONSUMER_SECRET = string
      OAUTH_TOKEN     = string
      OAUTH_SECRET    = string
    })
    CSV_FOLDER = string
  }))
}

variable "RUN_HOST" {
  type = string
}

locals {
  region = "asia-northeast1"
}
