# Final PDF Checklist

The assignment requires PDF submission with screenshots. Use this checklist before exporting.

## Source Code Evidence

- Screenshot of project structure with service folders:
  - `auth-service`
  - `product-service`
  - `order-service`
  - `user-service`
  - `chat-service`
  - `notification-service`
  - `frontend`
  - `gateway`
  - `terraform`
  - `k8s`
  - `ansible`

## Docker Evidence

- `docker compose up -d --build`
- `docker compose ps`
- Browser opened at `http://localhost`
- Product API response at `http://localhost/api/product`
- Notification counters at `http://localhost/notifications/events`
- Notification logs:
  ```bash
  docker compose logs --tail 50 notification-service
  ```

## Monitoring Evidence

- Prometheus targets page: `http://localhost:9090/targets`
- Prometheus alerts or rules page
- Grafana dashboard list
- Grafana `Clothes Store Overview` dashboard

## Docker Swarm Evidence

- `docker swarm init`
- `docker stack deploy -c docker-stack.yml clothes`
- `docker node ls`
- `docker stack services clothes`
- `docker stack ps clothes`
- `docker service logs clothes_notification-service --tail 50`

## Kubernetes Evidence

- `kubectl apply -f k8s/namespace.yaml`
- `kubectl apply -f k8s/app-stack.yaml`
- `kubectl apply -f k8s/monitoring.yaml`
- `kubectl get pods -n clothes-store`
- `kubectl get svc -n clothes-store`
- `kubectl get hpa -n clothes-store`
- `kubectl logs -n clothes-store deploy/notification-service`

## Ansible Evidence

- `ansible-playbook -i ansible/inventory.ini ansible/playbook.yml`
- Successful Ansible play recap
- Deployment status task output with `docker compose ps`

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
- Docker Swarm guide
- Assignment 4 incident response report
- Postmortem analysis
- Assignment 5 Terraform report
- Assignment 6 automation/capacity report
