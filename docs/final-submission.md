# Containerized Microservices System with Terraform and Incident Response

## Project Title

Design and Deployment of a Containerized Microservices System with Terraform-Based Infrastructure Provisioning and Incident Response Simulation

## System Overview

The project implements an online clothes store using a containerized microservices architecture. The system includes a web frontend, API gateway, independent backend services, database containers, Prometheus monitoring, Grafana dashboards, Terraform infrastructure provisioning, and an incident response simulation.

## Architecture

Core services:

- `frontend`: web user interface
- `gateway`: Nginx reverse proxy and routing layer
- `auth-service`: authentication and JWT token handling
- `product-service`: product retrieval and product management
- `order-service`: order creation, order status, and analytics
- `user-service`: user administration APIs
- `chat-service`: user chat endpoint
- `mongo`: current application database
- `postgres`: provisioned PostgreSQL database container
- `prometheus`: metrics collection
- `grafana`: dashboard visualization

## Functional Requirements Coverage

| Requirement | Implementation |
| --- | --- |
| Web interface | `frontend` service at `http://localhost:8080` through gateway |
| Authentication and authorization | `auth-service`, JWT, cookies |
| Product display | `product-service`, `/api/product` |
| Transactional operations | `order-service`, `/orders` |
| Independent backend services | Separate containers and service folders |
| Metrics collection | `/metrics` endpoints and Prometheus |
| Failure logging/detection | Docker logs, Prometheus targets, Grafana dashboard, alerts |

## Docker Deployment

Run:

```bash
docker compose up -d --build
```

Validate:

```bash
docker compose ps
curl http://localhost:8080/ping
curl http://localhost:8080/api/product
```

## Monitoring

Prometheus scrapes all application services. Grafana is provisioned with a Prometheus datasource and the `Clothes Store Overview` dashboard.

Monitoring URLs:

- Prometheus: `http://localhost:9090`
- Grafana: `http://localhost:3000`

## Terraform Infrastructure

Terraform files are located in `terraform/`.

Run:

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Terraform provisions an AWS EC2 virtual machine, security group rules, public IP, Docker installation, and optional application deployment.

## Incident Simulation

The incident simulates a database configuration failure in `order-service`.

Start incident:

```bash
docker compose -f docker-compose.yml -f docker-compose.incident.yml up -d order-service
```

Restore service:

```bash
docker compose -f docker-compose.yml up -d order-service
```

Observed root cause:

```text
dial tcp: lookup wrong-mongo-host on 127.0.0.11:53: no such host
```

## Reports

Detailed supporting documents:

- `docs/deployment-guide.md`
- `docs/assignment4-incident-response.md`
- `docs/postmortem.md`
- `docs/assignment5-terraform-report.md`
- `docs/final-pdf-checklist.md`

## Screenshot Placeholders

Insert screenshots into the final PDF:

1. Running containers
2. Application home page
3. Product API response
4. Prometheus targets
5. Grafana dashboard
6. Terraform init/plan/apply
7. EC2 public IP output
8. Incident before/failure/recovery states
