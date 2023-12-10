resource "google_service_account" "account" {
  account_id = var.id
}

resource "google_project_iam_member" "iam" {
  for_each   = toset(var.roles)
  member     = "serviceAccount:${google_service_account.account.email}"
  role       = "roles/${each.value}"
  project    = var.project_id
  depends_on = [google_service_account.account]
}

variable "id" {
  type = string
}
variable "roles" {
  type = list(string)
}
variable "project_id" {
  type = string
}
output "email" {
  value = google_service_account.account.email
}
