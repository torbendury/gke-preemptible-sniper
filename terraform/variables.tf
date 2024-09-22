variable "project_id" {
  description = "The project ID to deploy resources"
  type        = string
  default     = ""
}

variable "region" {
  description = "location for regional deployments"
  type        = string
  default     = "europe-west1"
}

variable "zone" {
  description = "location for zonal deployments"
  type        = string
  default     = "europe-west1-b"
}
