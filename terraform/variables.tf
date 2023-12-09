variable "PROJECT_ID" {
  type = string
}

variable "GOOGLE_CREDENTIALS" {
  type = any
}


variable "USERS_SECRET" {
  type = map(object({
    ZAIM = object({
      CONSUMER_KEY    = string
      CONSUMER_SECRET = string
      OAUTH_TOKEN = string
      OAUTH_SECRET = string
    })
    CSV_FOLDER = string
  }))
}

locals {
  region = "asia-northeast1"
}
