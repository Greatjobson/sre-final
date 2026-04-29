output "public_ip" {
  description = "Public IP address of the provisioned VM."
  value       = google_compute_address.app.address
}

output "app_url" {
  description = "Application URL through the gateway."
  value       = "http://${google_compute_address.app.address}:8080"
}

output "http_80_url" {
  description = "HTTP URL if the gateway is published on port 80 in production."
  value       = "http://${google_compute_address.app.address}"
}

output "prometheus_url" {
  description = "Prometheus URL."
  value       = "http://${google_compute_address.app.address}:9090"
}

output "grafana_url" {
  description = "Grafana URL."
  value       = "http://${google_compute_address.app.address}:3000"
}

output "ssh_command" {
  description = "SSH command for VM access."
  value       = "ssh ${var.ssh_username}@${google_compute_address.app.address}"
}
