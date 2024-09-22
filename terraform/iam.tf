resource "google_service_account" "gke_preemptible_sniper" {
  account_id   = "gke-preemptible-sniper"
  display_name = "GKE preemptible sniper"
}

resource "google_service_account_iam_member" "this" {
  service_account_id = google_service_account.gke_preemptible_sniper.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "serviceAccount:${var.project_id}.svc.id.goog[gke-preemptible-sniper/gke-preemptible-sniper]"
}

resource "google_project_iam_member" "this" {
  project = var.project_id
  role    = "roles/editor"
  member  = "serviceAccount:${google_service_account.gke_preemptible_sniper.email}"
}

resource "google_iam_workload_identity_pool" "this" {
  workload_identity_pool_id = "example"
}
