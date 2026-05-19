#!/usr/bin/env bash
set -euo pipefail

BRANCH="${1:-team-3}"
APP_DIR="${APP_DIR:-/opt/clothes-store}"

cd "$APP_DIR"

echo "[deploy] sync repository branch: ${BRANCH}"
git fetch origin "$BRANCH"
git checkout "$BRANCH"
git pull --ff-only origin "$BRANCH"

echo "[deploy] run pre-deploy validation"
make predeploy-check

echo "[deploy] build docker images"
make docker-build

echo "[deploy] apply kubernetes manifests"
make k8s-deploy

echo "[deploy] configure nginx proxy for domain"
sudo bash scripts/setup-nginx-proxy.sh

echo "[deploy] run health checks"
make health-check

echo "[deploy] completed successfully"
