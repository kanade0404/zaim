resource "google_cloud_scheduler_job" "scheduler" {
  name        = var.name
  description = var.description
  schedule    = var.schedule
  dynamic "pubsub_target" {
    for_each = var.pubsub_target == null ? [] : [1]
    content {
      topic_name = var.pubsub_target.topic_name
      data       = tostring(var.pubsub_target.data)
    }
  }
  time_zone = "Asia/Tokyo"
}

variable "name" {
  type = string
}

variable "schedule" {
  type = string
}

variable "description" {
  type    = string
  default = null
}

variable "pubsub_target" {
  type = object({
    topic_name = string
    data       = any
  })
  default = null
}
