# Clothes Store

A premium, modern online store platform built with Go. The project now includes a microservice-ready deployment topology with Nginx gateway routing across dedicated domain services.

## Team
- Zhumagali Beibarys

## Features

### Customer Experience
- **Modern Shop**: Advanced filtering (category, gender, color, size) and sorting.
- **Product Details**: High-quality imagery, size selection, and stock status.
- **Cart & Wishlist**: Persistent client-side shopping cart and wishlist management.
- **Checkout**: Seamless checkout flow with address management and order confirmation.
- **User Accounts**: Registration, login, and order history tracking.

### Admin Dashboard
- **Analytics**: Key performance indicators (Total Sales, Orders, Users).
- **Product Management**: Complete CRUD with local image uploads and advanced validation.
- **Order Management**: Track and update order statuses.
- **User Management**: Overview of registered users.

## Tech Stack

- **Backend API**: Go (Gin Web Framework)
- **Frontend Web**: Go-rendered templates + static assets (served by a separate web service)
- **Gateway / Reverse Proxy**: Nginx
- **Database**: MongoDB for application data + PostgreSQL container provisioned for assignment requirements/migration
- **Authentication**: JWT (JSON Web Tokens) with Secure Cookies
- **Frontend**: Semantic HTML5, Vanilla CSS (Modern CSS variables), JavaScript (ES6+)
- **Icons**: Lucide Icons

## Microservices Topology (Docker Compose)

- `gateway` (Nginx): single entry point and route dispatcher.
- `auth-service`: `/auth/*`
- `product-service`: `/api/product*`
- `order-service`: `/orders*`
- `user-service`: `/api/users*`
- `chat-service`: `/chat*`
- `frontend`: web UI service behind gateway.

## Database Performance

- **Multi-stage aggregation**: Analytics uses MongoDB pipelines (`$facet`, `$group`, `$lookup`, `$sort`) to compute totals, revenue trends, and top products without loading every order into memory.
- **Compound indexes**: `orders` uses `{ userId: 1, createdAt: -1 }` for user history and recent sorting; `order_items` uses `{ orderId: 1, productId: 1 }` to accelerate joins and product sales grouping.
- **Reduced transfer**: Aggregations return compact summaries and only a small window of recent orders.

## Project Structure

```bash
├── auth-service/
│   ├── cmd/server/main.go
│   └── internal/{domain,usecase,ports,adapters}
├── product-service/
│   ├── cmd/server/main.go
│   └── internal/{domain,usecase,ports,adapters}
├── order-service/
│   ├── cmd/server/main.go
│   └── internal/{domain,usecase,ports,adapters}
├── user-service/
│   ├── cmd/server/main.go
│   └── internal/{domain,usecase,ports,adapters}
├── chat-service/
│   ├── cmd/server/main.go
│   └── internal/{domain,usecase,ports,adapters}
├── backend/
│   └── Dockerfile       # Builds all backend service binaries
├── frontend/
│   └── Dockerfile       # Frontend container build
├── gateway/
│   └── nginx.conf       # API gateway routing
├── grafana/
│   ├── dashboards/
│   └── provisioning/
├── terraform/
│   ├── main.tf
│   ├── variables.tf
│   ├── outputs.tf
│   └── terraform.tfvars
└── docs/
    ├── deployment-guide.md
    ├── assignment4-incident-response.md
    ├── postmortem.md
    └── assignment5-terraform-report.md
```

## Setup & Installation

### Prerequisites
- Docker Desktop or Docker Engine with Docker Compose plugin
- Go 1.24+ for local tests
- Terraform 1.6+ for Assignment 5 infrastructure provisioning

### 1. Environment Configuration
Create a `.env` file in the root directory (you can copy from `.env.example`):
```env
JWT_SECRET=replace_with_strong_random_secret
ADMIN_EMAIL=admin@example.com
POSTGRES_DB=clothes_store
POSTGRES_USER=clothes
POSTGRES_PASSWORD=replace_with_strong_postgres_password
GRAFANA_ADMIN_PASSWORD=replace_with_strong_grafana_password
```

### 2. Run Tests
```bash
go test ./...
```

### 3. Run with Docker Compose
```bash
docker compose up -d --build
```
- Gateway (single entrypoint): http://localhost
- Frontend (direct access): http://localhost:8081
- Auth API via gateway: http://localhost/auth
- Product API via gateway: http://localhost/api/product
- Order API via gateway: http://localhost/orders
- User API via gateway: http://localhost/api/users
- Chat API via gateway: http://localhost/chat/messages
- MongoDB: localhost:27017
- PostgreSQL: localhost:5432
- Prometheus: http://localhost:9090
- Grafana: http://localhost:3000
- Grafana login: `admin` / `admin`

Docker build files:
- `frontend/Dockerfile` for web service
- `backend/Dockerfile` for backend service binaries

## Monitoring

- Prometheus scrapes every application service through `/metrics`.
- Prometheus also scrapes host metrics from `node-exporter` (`node-exporter:9100`) for system-level telemetry.
- Grafana is automatically provisioned with the Prometheus datasource.
- Grafana includes the `Clothes Store Overview` dashboard.
- Alerts are defined in `alerts.yml`, including high latency, high error rate, service down, high CPU usage, and restart loop warning.

## SRE Automation & Capacity Planning

Operational scripts are available in `scripts/`:

- `scripts/predeploy-check.sh` — validates compose config, required services, health checks, and env keys before deployment.
- `scripts/log-check.sh` — scans recent `docker compose logs` for DB failures and restart-loop patterns.
- `scripts/load-sim.sh` — concurrent request load simulation for a selected endpoint.
- `scripts/capacity-metrics.sh` — snapshot of CPU, memory, request rate, error rate, and restart frequency from Prometheus.
- `scripts/capacity-run.sh` — end-to-end capacity run (pre-load metrics, load scenarios, post-load metrics).

Example:

```bash
./scripts/predeploy-check.sh
./scripts/capacity-run.sh
./scripts/log-check.sh 15m
```

## Terraform Infrastructure

Terraform files are in `terraform/`.

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Before applying on Google Cloud, run `gcloud auth application-default login` and update `terraform/terraform.tfvars` with your `project_id`, zone, SSH public key, and optional repository URL.

## Incident Simulation

The order-service incident is simulated with an invalid database hostname:

```bash
docker compose -f docker-compose.yml -f docker-compose.incident.yml up -d order-service
```

Restore normal configuration:

```bash
docker compose -f docker-compose.yml up -d order-service
```

Reports and PDF source material are in `docs/`.

## API Documentation

### Authentication
- **POST** `/auth/register`
  - Request (JSON):
    ```json
    { "fullName": "John Doe", "email": "john@example.com", "password": "secret123" }
    ```
  - Response `201`:
    ```json
    { "message": "user registered" }
    ```
- **POST** `/auth/login`
  - Request (JSON):
    ```json
    { "email": "john@example.com", "password": "secret123" }
    ```
  - Response `200`:
    ```json
    { "token": "jwt-token" }
    ```
- **GET** `/auth/logout` (browser redirect)

**Auth for API**: send `Authorization: Bearer <token>` or cookie `auth_token`.

### Products
- **GET** `/api/product`
  - Response `200`:
    ```json
    [{ "id": "p1", "name": "Sneakers", "price": 120, "sizes": ["41","42"], "colors": ["black"] }]
    ```
- **GET** `/api/product/:id`
  - Response `200`: product object
- **POST** `/api/product` (admin, multipart/form-data)
  - Example:
    ```bash
    curl -X POST http://localhost:8080/api/product \
      -H "Authorization: Bearer <token>" \
      -F "name=Sneakers" -F "price=120" -F "category=shoes" \
      -F "gender=unisex" -F "sizes=41,42" -F "colors=black" \
      -F "stock=41:5,42:3" -F "image=@./sneakers.jpg"
    ```
  - Response `201`: product object
- **PUT** `/api/product/:id` (admin, JSON body = product)
- **DELETE** `/api/product/:id` (admin) → `204`

### Orders
- **GET** `/orders?user_id={userId}`
  - Response `200`: array of orders
- **POST** `/orders`
  - Request (JSON):
    ```json
    {
      "user_id": "u1",
      "payment_method": "card",
      "delivery_method": "courier",
      "delivery_address": "Almaty, Abay 10",
      "comment": "leave at door",
      "items": [
        {
          "product_id": "p1",
          "product_name": "Sneakers",
          "selected_size": "42",
          "selected_color": "black",
          "quantity": 1,
          "unit_price": 120
        }
      ]
    }
    ```
  - Response `201`: order object
- **GET** `/orders/:id`
  - Response `200`: order object
- **PATCH** `/orders/:id/status`
  - Request (JSON):
    ```json
    { "status": "completed" }
    ```
  - Response `200`:
    ```json
    { "id": "orderId", "status": "completed" }
    ```

### Analytics (admin)
- **GET** `/api/analytics/stats` → dashboard stats
- **GET** `/api/analytics/top-products` → top product sales
- **GET** `/api/analytics/revenue?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`
  - Response `200`:
    ```json
    [{ "date": "2026-02-01", "revenue": 520, "orders": 4 }]
    ```
- **GET** `/api/analytics/orders-status` → `{ "pending": 2, "completed": 5 }`

## Code Quality
- **Clean Architecture**: Separation of concerns between layers.
- **Optimized Assets**: Localized assets for faster loading and reliability.
- **Sanitized**: Codebase is free of redundant comments and junk files.
