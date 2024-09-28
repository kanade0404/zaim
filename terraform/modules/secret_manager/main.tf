resource "google_secret_manager_secret" "secret" {
  secret_id = var.id
  replication {
    auto {}
  }
  lifecycle {
    prevent_destroy = true
  }
}
resource "google_secret_manager_secret_version" "secret-version" {
  secret      = google_secret_manager_secret.secret.id
  secret_data = var.data
  depends_on  = [google_secret_manager_secret.secret]
}
variable "id" {
  type = string
}
variable "data" {
  type      = string
  sensitive = true
}
output "name" {
  value = google_secret_manager_secret.secret.name
}
output "secret_id" {
  value = google_secret_manager_secret.secret.secret_id
}
output "version" {
  value = google_secret_manager_secret_version.secret-version.version
}
