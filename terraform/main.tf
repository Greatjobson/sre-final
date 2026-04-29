terraform {
  required_version = ">= 1.6.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

provider "google" {
  project = var.project_id
  region  = var.gcp_region
  zone    = var.gcp_zone
}

resource "google_compute_address" "app" {
  name   = "${var.project_name}-ip"
  region = var.gcp_region
}

resource "google_compute_firewall" "ssh" {
  name    = "${var.project_name}-allow-ssh"
  network = var.network

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = [var.ssh_cidr]
  target_tags   = [var.network_tag]
}

resource "google_compute_firewall" "app" {
  name    = "${var.project_name}-allow-app"
  network = var.network

  allow {
    protocol = "tcp"
    ports    = [for port in var.app_ports : tostring(port)]
  }

  source_ranges = [var.app_cidr]
  target_tags   = [var.network_tag]
}

resource "google_compute_instance" "app" {
  name         = "${var.project_name}-vm"
  machine_type = var.machine_type
  zone         = var.gcp_zone
  tags         = [var.network_tag]

  boot_disk {
    initialize_params {
      image = var.os_image
      size  = var.disk_size_gb
      type  = var.disk_type
    }
  }

  network_interface {
    network = var.network

    access_config {
      nat_ip = google_compute_address.app.address
    }
  }

  metadata = {
    ssh-keys = "${var.ssh_username}:${var.ssh_public_key}"
  }

  metadata_startup_script = templatefile("${path.module}/user_data.sh.tftpl", {
    repository_url = var.repository_url
    repository_ref = var.repository_ref
    compose_file   = var.compose_file
  })

  labels = {
    project = var.project_name
  }
}
