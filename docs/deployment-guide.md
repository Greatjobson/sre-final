# Deployment Guide

## 1. Prerequisites

- Docker Desktop or Docker Engine with Docker Compose plugin
- Go 1.24+ for local builds and tests
- Terraform 1.6+ for infrastructure provisioning
- AWS account and an existing EC2 key pair for the Terraform scenario

## 2. Local Container Deployment

Build and start all services:

```bash
docker compose up -d --build
```

Validate containers:

```bash
docker compose ps
```

Main URLs:

- Application gateway: `http://localhost:8080`
- Direct frontend: `http://localhost:8081`
- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`
- Grafana credentials: `admin` / `admin`

## 3. Service Routes

The Nginx gateway routes traffic to independent backend services:

| Gateway Path | Service |
| --- | --- |
| `/auth` | `auth-service` |
| `/api/product` | `product-service` |
| `/orders` | `order-service` |
| `/api/users` | `user-service` |
| `/api/analytics` | `order-service` |
| `/chat` | `chat-service` |
| `/` | `frontend` |

## 4. Health Checks

```bash
curl http://localhost:8080/ping
curl http://localhost:8080/api/product
curl -X POST http://localhost:8080/auth/refresh
curl "http://localhost:8080/chat/messages?user_id=u1&peer_id=u2"
```

## 5. Monitoring

Prometheus is configured in `prometheus.yml` and scrapes:

- `frontend:8081`
- `auth-service:8082`
- `product-service:8083`
- `order-service:8084`
- `user-service:8085`
- `chat-service:8086`

Grafana is provisioned automatically with:

- Prometheus datasource
- `Clothes Store Overview` dashboard

## 6. Infrastructure Deployment with Terraform

Update `terraform/terraform.tfvars`:

```hcl
key_name       = "your-existing-aws-key-pair"
ssh_cidr       = "your-public-ip/32"
repository_url = "https://github.com/your-org/your-repo.git"
```

Run:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

After apply, Terraform outputs:

- public IP
- application URL
- Prometheus URL
- Grafana URL
- SSH command

## 7. Screenshot Checklist

Add these screenshots to the final PDF:

- `docker compose ps`
- application home page through `http://localhost:8080`
- product API response through gateway
- Prometheus targets page
- Grafana `Clothes Store Overview` dashboard
- Terraform `plan` or `apply` output
- Terraform outputs with public IP
