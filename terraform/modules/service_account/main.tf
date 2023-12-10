resource "google_service_account" "account" {
  account_id = var.id
}

resource "google_project_iam_member" "iam" {
  for_each   = toset(var.project_roles)
  member     = "serviceAccount:${google_service_account.account.email}"
  role       = "roles/${each.value}"
  project    = var.project_id
  depends_on = [google_service_account.account]
}

resource "google_service_account_iam_member" "service_iam" {
  for_each           = toset(var.service_roles)
  member             = "serviceAccount:${google_service_account.account.email}"
  role               = "roles/${each.value}"
  service_account_id = google_service_account.account.name
}

variable "id" {
  type = string
}
variable "project_roles" {
  type    = list(string)
  default = []
}
variable "service_roles" {
  type    = list(string)
  default = []
}
variable "project_id" {
  type = string
}
output "email" {
  value = google_service_account.account.email
}
