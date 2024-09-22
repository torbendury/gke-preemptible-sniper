resource "google_container_cluster" "this" {
  name                        = "gke-preemptinble" # TODO typo correction later
  location                    = var.zone
  remove_default_node_pool    = true
  default_max_pods_per_node   = 32
  initial_node_count          = 1
  min_master_version          = "1.30"
  logging_service             = "logging.googleapis.com/kubernetes"
  monitoring_service          = "monitoring.googleapis.com/kubernetes"
  network                     = google_compute_network.this.self_link
  subnetwork                  = google_compute_subnetwork.this.self_link
  enable_shielded_nodes       = true
  enable_intranode_visibility = true

  release_channel {
    channel = "REGULAR"
  }

  cluster_autoscaling {
    enabled = false
  }

  binary_authorization {
    evaluation_mode = "DISABLED"
  }

  maintenance_policy {
    daily_maintenance_window {
      start_time = "08:00"
    }
  }

  master_authorized_networks_config {
    cidr_blocks {
      cidr_block   = "0.0.0.0/0"
      display_name = "Everyone"
    }
  }

  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  private_cluster_config {
    enable_private_endpoint = false
    enable_private_nodes    = true
    master_ipv4_cidr_block  = "172.19.7.240/28"
  }

  vertical_pod_autoscaling {
    enabled = true
  }

  addons_config {
    http_load_balancing {
      disabled = false
    }

    horizontal_pod_autoscaling {
      disabled = false
    }

    cloudrun_config {
      disabled = true
    }

    network_policy_config {
      disabled = false
    }

    gce_persistent_disk_csi_driver_config {
      enabled = true
    }

    dns_cache_config {
      enabled = true
    }

    gke_backup_agent_config {
      enabled = true
    }
  }

  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "172.22.192.0/18"
    services_ipv4_cidr_block = "172.22.0.0/22"
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }
}
