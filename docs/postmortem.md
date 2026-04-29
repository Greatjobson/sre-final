# Postmortem: Order Service Database Configuration Failure

## Incident Overview

The `order-service` experienced a simulated outage caused by an invalid MongoDB connection string. The service was unable to connect to the database and could not process order-related requests.

## Customer Impact

Affected users could not:

- Create new orders
- View order history
- Check order status

Unaffected functionality:

- Product browsing
- User registration and login
- User profile pages
- Chat endpoint
- Monitoring stack

## Root Cause Analysis

The root cause was an incorrect database hostname in the `MONGODB_URI` environment variable:

```env
MONGODB_URI=mongodb://wrong-mongo-host:27017
```

In Docker Compose, service DNS names are based on service names. The correct MongoDB hostname is:

```env
MONGODB_URI=mongodb://mongo:27017
```

Because `order-service` depends on MongoDB during startup, the service failed to initialize correctly when the hostname was invalid.

## Detection and Response Evaluation

Detection methods:

- Container status through `docker compose ps`
- Logs through `docker compose logs order-service`
- Prometheus target status
- Grafana availability dashboard
- User-facing order endpoint failures

The response was effective because the failure was isolated to one service. Product browsing, authentication, user management, and chat remained available.

## Resolution Summary

The service was restored by reverting the faulty database connection string and restarting the service:

```bash
docker compose -f docker-compose.yml up -d order-service
```

After restart, the service reconnected to MongoDB and order endpoints became available again.

## Lessons Learned

- Environment variables are critical production dependencies.
- Service-specific failures should be isolated by microservice boundaries.
- Monitoring should include both service availability and request error rates.
- Logs must make configuration failures easy to identify.
- Incident simulation is useful for validating recovery procedures.

## Action Items

| Action Item | Owner | Priority |
| --- | --- | --- |
| Add service health checks for all backend services | Backend team | High |
| Add readiness checks for database connectivity | Backend team | High |
| Restrict SSH CIDR in Terraform to a known IP range | DevOps team | Medium |
| Add alert for unavailable Prometheus targets | DevOps team | High |
| Document rollback and service restart commands | DevOps team | Medium |
| Add CI validation for required environment variables | DevOps team | Medium |

## Prevention Plan

To reduce future risk:

1. Add Docker health checks for `/ping`.
2. Add startup validation for required environment variables.
3. Use secret/config management for production connection strings.
4. Add dashboards for service-specific availability.
5. Run incident simulations before final deployment.
