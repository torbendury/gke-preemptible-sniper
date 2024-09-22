resource "google_compute_network" "this" {
  project                 = var.project_id
  name                    = "gke-preemptible-sniper"
  auto_create_subnetworks = false
  routing_mode            = "GLOBAL"
}

resource "google_compute_subnetwork" "this" {
  name                     = "gke-nodes"
  project                  = var.project_id
  region                   = var.region
  network                  = google_compute_network.this.self_link
  ip_cidr_range            = "172.19.6.0/24"
  private_ip_google_access = true
}

resource "google_service_networking_connection" "peering" {
  network                 = google_compute_network.this.self_link
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.service-allocation.name]
}

resource "google_compute_global_address" "service-allocation" {
  project       = var.project_id
  name          = "service-allocation"
  purpose       = "VPC_PEERING"
  address       = "172.22.16.0"
  prefix_length = 22
  address_type  = "INTERNAL"
  network       = google_compute_network.this.self_link
}
