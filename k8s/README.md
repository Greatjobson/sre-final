# Kubernetes Deployment

This folder demonstrates the Kubernetes orchestration requirement for the SRE project.

## Apply

```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/app-stack.yaml
kubectl apply -f k8s/monitoring.yaml
```

## Verify

```bash
kubectl get pods -n clothes-store
kubectl get svc -n clothes-store
kubectl get hpa -n clothes-store
kubectl logs -n clothes-store deploy/notification-service
```

Default NodePort endpoints:

- Gateway: `http://localhost:30080`
- Prometheus: `http://localhost:30090`
- Grafana: `http://localhost:30300`

The application images use the local Docker Compose image names, so on Docker Desktop Kubernetes build them first:

```bash
docker compose build
```
