# Assignment 5 Report: Terraform Infrastructure Provisioning

## Objective

The objective of Assignment 5 is to provision reproducible infrastructure for the containerized microservices system using Terraform. The infrastructure code creates a cloud virtual machine, configures access rules, installs Docker, and optionally starts the application with Docker Compose.

## Infrastructure Provider

The implementation uses Google Cloud Platform (Compute Engine) as the cloud VM provider.

Terraform files:

- `terraform/main.tf`
- `terraform/variables.tf`
- `terraform/outputs.tf`
- `terraform/terraform.tfvars`
- `terraform/user_data.sh.tftpl`

## Provisioned Resources

The Terraform configuration provisions:

- One Ubuntu 22.04 Compute Engine instance
- Two VPC firewall rules (SSH + application ports)
- One static external IP (`google_compute_address`)
- Docker Engine and Docker Compose plugin through startup script
- Optional application deployment from a Git repository

## Network Access Rules

The firewall rules allow:

| Port | Purpose |
| --- | --- |
| `22` | SSH administration |
| `80` | HTTP application access |
| `8080` | Local gateway mapping used by Docker Compose |
| `3000` | Grafana dashboard access |
| `9090` | Prometheus access |

SSH access is controlled by `ssh_cidr`. Application access is controlled by `app_cidr`.

## Reproducibility

The infrastructure is reproducible through the standard Terraform workflow:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

The `terraform.tfvars` file contains environment-specific values such as:

- GCP project ID, region, and zone
- machine type and disk settings
- SSH username and public key
- allowed CIDR ranges
- Git repository URL
- Git branch or tag

## Outputs

Terraform outputs:

- `public_ip`: static public IP address of the Compute Engine instance
- `app_url`: application URL through the gateway
- `http_80_url`: HTTP URL when gateway is exposed on port 80
- `prometheus_url`: Prometheus monitoring URL
- `grafana_url`: Grafana monitoring URL
- `ssh_command`: SSH command for instance access

## Deployment Flow

1. Terraform provisions the Compute Engine VM and static IP.
2. The user data script installs Docker and Docker Compose.
3. If `repository_url` is set, the VM clones the repository.
4. Docker Compose starts the microservices stack.
5. The application becomes reachable through the public IP.

## Evidence to Include in PDF

Add screenshots of:

- `terraform init`
- `terraform plan`
- `terraform apply`
- Terraform outputs showing public IP
- GCP Compute Engine instance running
- VPC firewall rules
- Application running through the public IP
- Prometheus and Grafana pages through the public IP
