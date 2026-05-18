.PHONY: predeploy-check docker-build k8s-deploy health-check deploy-team3

predeploy-check:
	bash scripts/predeploy-check.sh

docker-build:
	docker compose build

k8s-deploy:
	kubectl apply -f k8s/namespace.yaml
	kubectl apply -f k8s/app-stack.yaml
	kubectl apply -f k8s/monitoring.yaml

health-check:
	kubectl get pods -n clothes-store
	kubectl get hpa -n clothes-store
	kubectl rollout status deployment/gateway -n clothes-store --timeout=180s
	curl -fsS http://localhost:30080/health

deploy-team3:
	bash scripts/deploy-team3.sh team-3
