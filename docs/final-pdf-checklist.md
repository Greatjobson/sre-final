# Final PDF Checklist

The assignment requires PDF submission with screenshots. Use this checklist before exporting.

## Source Code Evidence

- Screenshot of project structure with service folders:
  - `auth-service`
  - `product-service`
  - `order-service`
  - `user-service`
  - `chat-service`
  - `frontend`
  - `gateway`
  - `terraform`

## Docker Evidence

- `docker compose up -d --build`
- `docker compose ps`
- Browser opened at `http://localhost:8080`
- Product API response at `http://localhost:8080/api/product`
- Auth refresh response from `POST http://localhost:8080/auth/refresh`

## Monitoring Evidence

- Prometheus targets page: `http://localhost:9090/targets`
- Prometheus alerts or rules page
- Grafana dashboard list
- Grafana `Clothes Store Overview` dashboard

## Terraform Evidence

- `terraform init`
- `terraform plan`
- `terraform apply`
- Terraform outputs:
  - public IP
  - app URL
  - Prometheus URL
  - Grafana URL
- AWS EC2 instance page
- AWS security group inbound rules

## Incident Simulation Evidence

- Healthy system before incident
- Fault injection command:
  ```bash
  docker compose -f docker-compose.yml -f docker-compose.incident.yml up -d order-service
  ```
- Failed `order-service` logs
- Prometheus/Grafana showing failure
- Recovery command:
  ```bash
  docker compose -f docker-compose.yml up -d order-service
  ```
- Healthy system after recovery

## Documents to Include

- README/setup instructions
- Deployment guide
- Assignment 4 incident response report
- Postmortem analysis
- Assignment 5 Terraform report
