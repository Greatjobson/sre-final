# Assignment 6 Report: Automation in SRE and Capacity Planning

## 1. Objective

This assignment extends the microservices system with automation mechanisms and capacity-planning practices to improve reliability, reduce manual intervention, and prepare the platform for higher load.

## 2. Implemented Automation Mechanisms

### 2.1 Automated Deployment

- Multi-container deployment is automated with `docker compose up -d --build`.
- Infrastructure provisioning is managed with Terraform in `terraform/`.
- Runtime configuration is standardized through environment variables (`.env` and service env blocks).

### 2.2 Health Checks and Self-Healing

- Services expose HTTP health endpoints (`/health`).
- `docker-compose.yml` contains per-service health checks.
- `restart: unless-stopped` is configured to improve recovery behavior.

### 2.3 Monitoring-Based Alerting

- Prometheus scrapes service metrics and host metrics (`node-exporter`).
- Alert rules in `alerts.yml` include:
  - `HighCPUUsageCritical`
  - `ServiceDown`
  - `HighErrorRateCritical`
  - `HighLatencyWarning`
  - `ServiceRestartLoopWarning`

### 2.4 Log-Based Troubleshooting Automation

- Container logs are centralized via Docker.
- `scripts/log-check.sh` automates log scanning for:
  - Database connection failures (host resolution, dial, selection errors)
  - Restart-loop indicators (restarts, unhealthy/exited patterns)

### 2.5 Configuration Validation

- `scripts/predeploy-check.sh` performs pre-deployment checks:
  - compose syntax validation (`docker compose config -q`)
  - required service presence
  - health-check endpoint sanity checks
  - invalid Mongo hostname guard (`wrong-mongo-host`)
  - required env key checks (`JWT_SECRET`, `ADMIN_EMAIL`) when `.env` is present

## 3. Capacity Planning

### 3.1 Metrics Collection

Capacity-related metrics are collected from Prometheus and exposed through `scripts/capacity-metrics.sh`:

1. CPU usage (%)
2. Memory utilization (%)
3. Request rate (req/sec)
4. Error rate (%)
5. Restart frequency (changes in `process_start_time_seconds`)

### 3.2 Load Simulation

Single endpoint load simulation:

```bash
REQUESTS=1000 CONCURRENCY=50 BASE_URL=http://localhost ./scripts/load-sim.sh /api/product
```

Full run scenario:

```bash
BASE_URL=http://localhost REQUESTS=600 CONCURRENCY=30 ./scripts/capacity-run.sh
```

### 3.3 Observation Template

Use the following table after running load tests and screenshots:

| Scenario | Concurrency | Requests | Avg req/sec | Error rate | CPU usage | Memory usage | Restart frequency |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Baseline | - | - | [fill] | [fill] | [fill] | [fill] | [fill] |
| Homepage load (`/`) | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] |
| Product API load (`/api/product`) | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] |
| Auth refresh load (`/auth/refresh`) | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] | [fill] |

### 3.4 Capacity Analysis Guidelines

Evaluate:

1. Maximum sustainable request rate before errors rise.
2. Resource consumption trends by service and host.
3. Failure thresholds (latency/error/restart behavior).
4. Resource-heavy services (typically order and API-heavy paths).

## 4. Scaling Strategy

### 4.1 Horizontal Scaling

- Scale stateless services (especially `order-service`) with multiple instances.
- Add load balancing at gateway or orchestrator level.

### 4.2 Vertical Scaling

- Increase VM/container CPU and memory limits when sustained CPU or memory pressure is detected.
- Apply infrastructure changes through Terraform.

### 4.3 Database Optimization

- Optimize high-frequency queries and indexes.
- Use connection tuning and workload-aware resource sizing.

### 4.4 Auto-Scaling Considerations

- Use metric-based thresholds (CPU, latency, error rate) for scaling triggers.
- For advanced policy automation, migrate runtime orchestration to Kubernetes/HPA when required.

## 5. Evidence Checklist

Add screenshots for:

1. `predeploy-check.sh` successful output
2. Prometheus Targets (including `node-exporter`)
3. Active alerts page/rules
4. Load simulation execution output
5. Capacity metrics before and after load
6. Grafana dashboard under load
7. `log-check.sh` output (normal and/or incident pattern)
