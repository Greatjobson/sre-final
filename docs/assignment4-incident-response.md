# Assignment 4 Report: Incident Response Simulation

## Incident Summary

An incident was simulated in the `order-service` by introducing an invalid database hostname. This represents a realistic production configuration failure where a service cannot connect to its database after deployment.

The affected service was:

- `order-service`

The simulated faulty configuration was:

```env
MONGODB_URI=mongodb://wrong-mongo-host:27017
```

## Incident Scenario

The application normally runs with:

```env
MONGODB_URI=mongodb://mongo:27017
```

To simulate failure, the incident override file changes only the order service database connection:

```bash
docker compose -f docker-compose.yml -f docker-compose.incident.yml up -d order-service
```

Expected result:

- `order-service` fails to initialize or becomes unavailable
- `/orders` and order-related UI features stop working
- Other services continue running
- Prometheus shows degraded availability for `order-service`
- Logs show database connection or hostname resolution errors

## Impact Assessment

Customer impact:

- Users cannot create new orders.
- Users cannot view order details or order history.
- Product browsing, login, registration, and chat remain available.

Business impact:

- Checkout flow is blocked.
- Revenue-generating transactions are unavailable during the incident.

## Severity Classification

Severity: `High`

Reason:

- The incident affects a critical transactional service.
- The entire system is not down, but order creation is a core business function.

## Timeline of Events

| Time | Event |
| --- | --- |
| T+00 | System is healthy. All containers are running. |
| T+01 | Faulty `MONGODB_URI` is applied to `order-service`. |
| T+02 | `order-service` fails or becomes unavailable. |
| T+03 | Incident detected through failed order requests and monitoring. |
| T+04 | Logs are inspected with `docker compose logs order-service`. |
| T+05 | Root cause identified as incorrect database hostname. |
| T+06 | Configuration is restored to `mongodb://mongo:27017`. |
| T+07 | `order-service` is restarted. |
| T+08 | Service health and order endpoints are verified. |

## Detection

The incident can be detected through:

```bash
docker compose ps
docker compose logs order-service
curl http://localhost:8080/orders?user_id=test-user
```

Monitoring evidence:

- Prometheus target status for `order-service`
- Grafana `Service Availability` panel
- Grafana `Error Rate` panel
- Prometheus alert `ServiceDown`

## Analysis

The logs show that `order-service` cannot connect to MongoDB because the hostname is invalid. In Docker Compose networking, services resolve each other by service name. The correct hostname is `mongo`.

Observed log evidence:

```text
MongoDB: server selection error: context deadline exceeded
dial tcp: lookup wrong-mongo-host on 127.0.0.11:53: no such host
```

Root cause:

- Misconfigured `MONGODB_URI` environment variable for `order-service`.

Contributing factors:

- No pre-deployment validation for database connection strings.
- No readiness check blocking rollout of an unhealthy service.

## Mitigation

Restore the correct database configuration:

```bash
docker compose up -d order-service
```

If the incident override stack is active, return to the normal Compose configuration:

```bash
docker compose -f docker-compose.yml up -d order-service
```

Validate:

```bash
docker compose ps
docker compose logs --tail 50 order-service
curl http://localhost:8080/api/product
curl "http://localhost:8080/orders?user_id=test-user"
```

## Resolution Confirmation

The incident is resolved when:

- `order-service` is running.
- Prometheus target for `order-service` is `UP`.
- Grafana service availability panel returns to healthy state.
- Order endpoints respond through the gateway.
- Other services remain unaffected.

## Screenshots Required for PDF

Add screenshots of:

- System before incident: all containers running
- Prometheus targets before incident
- Grafana dashboard before incident
- Fault injected with `docker-compose.incident.yml`
- `order-service` logs showing database connection failure
- Prometheus or Grafana showing degraded service
- Restored configuration
- System after recovery
