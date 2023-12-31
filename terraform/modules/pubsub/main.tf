resource "google_pubsub_topic" "topic" {
  name = var.name
}

resource "google_pubsub_subscription" "subscription" {
  for_each = { for sub in var.subscriptions : sub.name => sub }
  name     = each.key
  topic    = google_pubsub_topic.topic.id
  dynamic "dead_letter_policy" {
    for_each = each.value.dead_letter_topic == null ? [] : [1]
    content {
      dead_letter_topic     = each.value.dead_letter_topic.id
      max_delivery_attempts = each.value.dead_letter_topic.max_delivery_attempts
    }
  }
  dynamic "push_config" {
    for_each = each.value.push == null ? [] : [1]
    content {
      push_endpoint = each.value.push.endpoint
      oidc_token {
        service_account_email = each.value.push.service_account_email
      }
    }
  }
  depends_on = [google_pubsub_topic.topic]
}

variable "name" {
  type = string
}

variable "subscriptions" {
  type = list(object({
    name = string
    dead_letter_topic = optional(object({
      id                    = string
      max_delivery_attempts = optional(number)
    }))
    push = optional(object({
      endpoint              = string
      service_account_email = string
    }))
  }))
}

output "topic_id" {
  value = google_pubsub_topic.topic.id
}
