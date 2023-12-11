resource "google_storage_bucket" "bucket" {
  location                    = var.location
  name                        = var.name
  uniform_bucket_level_access = var.uniform_bucket_level_access
}

variable "location" {
  type = string
}

variable "name" {
  type = string
}

variable "uniform_bucket_level_access" {
  type    = bool
  default = true
}

output "name" {
  value = google_storage_bucket.bucket.name
}
