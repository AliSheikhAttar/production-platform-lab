#!/usr/bin/env bash
set -euo pipefail

LAB_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
NATIVE_LINUX="${LAB_ROOT}/infrastructure/native-linux"

echo "TaskFlow Stage 1 — Native Linux Services"
echo "Lab root: ${LAB_ROOT}"
echo

if [[ "${EUID}" -ne 0 ]]; then
	echo "Run with sudo so systemd units and /etc/taskflow can be installed."
	exit 1
fi

install -d -m 0750 /etc/taskflow
install -m 0640 "${NATIVE_LINUX}/env/task-service.env.example" /etc/taskflow/task-service.env
install -m 0640 "${NATIVE_LINUX}/env/api-gateway.env.example" /etc/taskflow/api-gateway.env

echo "Building binaries..."
sudo -u asa bash -c "
	cd '${LAB_ROOT}/application/task-service' && go build -trimpath -o bin/task-service ./cmd/task-service
	cd '${LAB_ROOT}/application/api-gateway' && go build -trimpath -o bin/gateway ./cmd/gateway
"

echo "Setting up native PostgreSQL (port 5433)..."
sudo -u postgres psql -p 5433 -f "${NATIVE_LINUX}/postgres/001-setup.sql" || true
sudo -u postgres psql -p 5433 -d taskflow -f "${LAB_ROOT}/application/task-service/migrations/001_create_tasks.sql" || true

sed -i "s|/home/asa/asa/production-platform-lab|${LAB_ROOT}|g" \
	"${NATIVE_LINUX}/systemd/task-service.service" \
	"${NATIVE_LINUX}/systemd/api-gateway.service"

install -m 0644 "${NATIVE_LINUX}/systemd/task-service.service" /etc/systemd/system/task-service.service
install -m 0644 "${NATIVE_LINUX}/systemd/api-gateway.service" /etc/systemd/system/api-gateway.service

systemctl daemon-reload
systemctl enable task-service.service api-gateway.service
systemctl restart task-service.service api-gateway.service

echo
echo "Services enabled. Check status:"
echo "  systemctl status task-service api-gateway"
echo "  journalctl -u task-service -f"
echo "  journalctl -u api-gateway -f"
echo
echo "Test:"
echo "  curl -s http://127.0.0.1:8080/tasks"
echo
echo "Optional nginx reverse proxy:"
echo "  sudo cp ${NATIVE_LINUX}/nginx/taskflow.conf /etc/nginx/sites-available/taskflow"
echo "  sudo ln -sf /etc/nginx/sites-available/taskflow /etc/nginx/sites-enabled/taskflow"
echo "  sudo nginx -t && sudo systemctl reload nginx"
