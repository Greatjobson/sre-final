variable "project_name" {
  description = "Name prefix for created resources."
  type        = string
  default     = "clothes-store"
}

variable "project_id" {
  description = "Google Cloud project ID."
  type        = string
}

variable "gcp_region" {
  description = "Google Cloud region where resources will be created."
  type        = string
  default     = "europe-west3"
}

variable "gcp_zone" {
  description = "Google Cloud zone where the VM will be created."
  type        = string
  default     = "europe-west3-c"
}

variable "machine_type" {
  description = "Compute Engine machine type."
  type        = string
  default     = "e2-medium"
}

variable "os_image" {
  description = "Boot image for the VM."
  type        = string
  default     = "ubuntu-os-cloud/ubuntu-2204-lts"
}

variable "disk_size_gb" {
  description = "Boot disk size in GB."
  type        = number
  default     = 20
}

variable "disk_type" {
  description = "Boot disk type, for example pd-standard or pd-ssd."
  type        = string
  default     = "pd-standard"
}

variable "network" {
  description = "VPC network name."
  type        = string
  default     = "default"
}

variable "network_tag" {
  description = "Network tag applied to the VM and used by firewall rules."
  type        = string
  default     = "clothes-store-app"
}

variable "ssh_username" {
  description = "SSH username for VM login."
  type        = string
  default     = "greatti"
}

variable "ssh_public_key" {
  description = "SSH public key content for VM access."
  type        = string
}

variable "ssh_cidr" {
  description = "CIDR block allowed to connect to SSH."
  type        = string
  default     = "0.0.0.0/0"
}

variable "app_cidr" {
  description = "CIDR block allowed to access app ports."
  type        = string
  default     = "0.0.0.0/0"
}

variable "app_ports" {
  description = "Application ports exposed by firewall."
  type        = list(number)
  default     = [80, 443, 8080, 3000, 9090]
}

variable "repository_url" {
  description = "Git repository URL to clone on the VM. Leave empty to provision only infrastructure."
  type        = string
  default     = ""
}

variable "repository_ref" {
  description = "Git branch, tag, or commit to deploy."
  type        = string
  default     = "main"
}

variable "compose_file" {
  description = "Docker Compose file to run on the VM."
  type        = string
  default     = "docker-compose.yml"
}
