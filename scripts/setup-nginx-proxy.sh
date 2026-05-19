#!/usr/bin/env bash
set -euo pipefail

DOMAIN_NAME="${DOMAIN_NAME:-greatjobson.me}"
UPSTREAM_HOST="${UPSTREAM_HOST:-127.0.0.1}"
UPSTREAM_PORT="${UPSTREAM_PORT:-30080}"
SITE_CONF="/etc/nginx/sites-available/clothes-store"

if ! command -v nginx >/dev/null 2>&1; then
  apt-get update
  apt-get install -y nginx
fi

cat > "$SITE_CONF" <<EOF
server {
    listen 80 default_server;
    server_name ${DOMAIN_NAME} _;

    location / {
        proxy_pass http://${UPSTREAM_HOST}:${UPSTREAM_PORT};
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}
EOF

ln -sf "$SITE_CONF" /etc/nginx/sites-enabled/clothes-store
rm -f /etc/nginx/sites-enabled/default

nginx -t
systemctl restart nginx
