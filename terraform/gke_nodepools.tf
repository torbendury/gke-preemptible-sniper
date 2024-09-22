locals {
  nodepools = yamldecode(file("${path.module}/etc/nodepools.yaml"))
}

resource "google_container_node_pool" "this" {
  for_each = { for nodepool in local.nodepools : nodepool["name"] => nodepool }

  name               = each.value["name"]
  cluster            = google_container_cluster.this.name
  location           = var.zone
  project            = var.project_id
  initial_node_count = each.value["initial_node_count"]

  lifecycle {
    ignore_changes = [
      # initial node count is modified when scaling through the ui
      # changed value leads to destroying and recreating the node pool
      initial_node_count
    ]
  }

  autoscaling {
    min_node_count = each.value["autoscaling"]["min_node_count"]
    max_node_count = each.value["autoscaling"]["max_node_count"]
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }

  upgrade_settings {
    max_unavailable = each.value["upgrade_settings"]["max_unavailable"]
    max_surge       = each.value["upgrade_settings"]["max_surge"]
  }

  node_config {
    preemptible  = each.value["preemptible"]
    machine_type = each.value["machine_type"]
    image_type   = "COS_CONTAINERD"
    disk_size_gb = each.value["disk"]["size_gb"]
    disk_type    = each.value["disk"]["type"]
    tags         = ["gkenode"]

    kubelet_config {
      cpu_manager_policy = "none"
    }

    labels = each.value["labels"]
    metadata = {
      disable-legacy-endpoints = "true"
    }

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
      "https://www.googleapis.com/auth/devstorage.read_only",
    ]
  }
}
