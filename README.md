# Production Platform Lab

A project-based DevOps and Site Reliability Engineering laboratory that builds a complete production-style application platform from the ground up.

The project begins with a manually deployed Linux service and gradually evolves into an automated, containerized, observable, highly available Kubernetes platform.

The goal is not simply to install DevOps tools. The goal is to understand how production systems are designed, deployed, operated, monitored, secured, troubleshot, and recovered from failure.

---

## Project Overview

Production Platform Lab simulates the infrastructure and operational environment of a real software company using local virtual machines.

The platform supports the complete software delivery lifecycle:

```text
Source Code
    ↓
Git Repository
    ↓
CI/CD Pipeline
    ↓
Container Registry
    ↓
Kubernetes Deployment
    ↓
Monitoring, Logging, Alerting, and Recovery
```

The final environment includes:

* Linux virtual machines
* Automated server configuration
* Containerized application services
* CI/CD pipelines
* Kubernetes orchestration
* Centralized monitoring and logging
* PostgreSQL high availability
* RabbitMQ or Kafka messaging
* Distributed storage
* SLI and SLO monitoring
* Incident-response procedures
* Backup and disaster recovery

---

## Objectives

This project is designed to develop practical skills in:

* Linux systems administration
* TCP/IP networking
* DNS, HTTP, HTTPS, and TLS
* Bash, Python, and Go automation
* Git-based collaboration
* Infrastructure automation with Ansible
* Docker and container internals
* GitLab CI/CD
* Kubernetes administration
* Prometheus and Grafana
* Centralized logging
* PostgreSQL operations and high availability
* RabbitMQ and Kafka
* Ceph distributed storage
* Site Reliability Engineering
* Incident management
* Technical documentation

---

## Final Architecture

```text
                         Developers
                              │
                              ▼
                     Git Repository
                              │
                              ▼
                       CI/CD Pipeline
                  ┌───────────┼───────────┐
                  │           │           │
                Tests       Build       Scan
                              │
                              ▼
                      Container Registry
                              │
                              ▼
                    Kubernetes Cluster
         ┌────────────────────┼────────────────────┐
         │                    │                    │
         ▼                    ▼                    ▼
    API Replicas       Background Workers      Frontend
         │                    │
         ├─────────────┬──────┘
         │             │
         ▼             ▼
 PostgreSQL HA    RabbitMQ / Kafka
         │
         ▼
 Distributed Persistent Storage
```

Observability services:

```text
Prometheus
Alertmanager
Grafana
Elasticsearch
Kibana
```

Infrastructure services:

```text
Linux Virtual Machines
KVM / QEMU
Ansible
Docker / containerd
Kubernetes
DNS
TLS
Nginx / HAProxy
Firewall Rules
Backup Systems
```

---

## Application Workload

The infrastructure hosts a small task-processing application.

The application is intentionally simple because the primary focus of the project is infrastructure and operations.

Application components:

* REST API
* Web interface
* PostgreSQL database
* Background worker
* RabbitMQ or Kafka messaging
* Health endpoints
* Prometheus metrics endpoint
* Structured application logs

Example functionality:

* Create a task
* Assign a task
* Update task status
* Process tasks asynchronously
* Store task history
* View task-processing results

---

## Project Stages

Each stage builds on the previous stage.

### Stage 1 — Native Linux Service

Deploy the application directly on a Linux virtual machine.

Topics:

* Linux users and permissions
* Process management
* systemd
* journald
* Nginx
* PostgreSQL
* SSH
* Firewall configuration
* Backup scripts

Tag:

```text
v1-linux-service
```

---

### Stage 2 — Networking and TLS

Separate the application, proxy, and database across multiple virtual machines.

Topics:

* IP addressing
* Subnetting
* Routing
* TCP and UDP
* DNS
* HTTP and HTTPS
* TLS certificates
* Reverse proxies
* Network troubleshooting

Tag:

```text
v2-networking-tls
```

---

### Stage 3 — Ansible Automation

Automate the complete Linux environment.

Topics:

* Inventories
* Playbooks
* Roles
* Variables
* Templates
* Handlers
* Ansible Vault
* Idempotence
* Configuration drift

Tag:

```text
v3-ansible
```

---

### Stage 4 — Containerization

Convert the application into a multi-container environment.

Topics:

* Dockerfiles
* Multi-stage builds
* Container networking
* Persistent volumes
* Health checks
* Non-root containers
* Resource limits
* Docker Compose

Tag:

```text
v4-docker
```

---

### Stage 5 — CI/CD

Build an automated software delivery pipeline.

Pipeline stages:

```text
Lint
  ↓
Unit Test
  ↓
Integration Test
  ↓
Build Container Image
  ↓
Security Scan
  ↓
Push Image
  ↓
Deploy to Staging
  ↓
Smoke Test
  ↓
Manual Approval
  ↓
Deploy to Production
```

Topics:

* GitLab Runner
* Pipeline jobs
* Artifacts
* Caching
* Container registry
* Protected variables
* Release versioning
* Rollbacks

Tag:

```text
v5-cicd
```

---

### Stage 6 — Kubernetes

Move the platform to a multi-node Kubernetes cluster.

Topics:

* Pods
* Deployments
* StatefulSets
* Services
* Ingress
* ConfigMaps
* Secrets
* PersistentVolumes
* Resource requests and limits
* Health probes
* Horizontal autoscaling
* RBAC
* NetworkPolicy
* Helm
* Cluster upgrades
* etcd backups

Tag:

```text
v6-kubernetes
```

---

### Stage 7 — Observability

Add monitoring, logging, dashboards, and alerting.

Topics:

* Prometheus
* PromQL
* Grafana
* Alertmanager
* Node Exporter
* Application metrics
* PostgreSQL metrics
* Kubernetes metrics
* Elasticsearch
* Kibana
* Structured logs
* Correlation IDs

Tag:

```text
v7-observability
```

---

### Stage 8 — SRE Practices

Add reliability objectives and incident-management processes.

Topics:

* SLI
* SLO
* SLA
* Error budgets
* Burn-rate alerts
* Incident severity
* On-call procedures
* Runbooks
* Blameless postmortems
* Toil reduction

Example SLO:

```text
99.9% of valid HTTP requests must complete successfully
within 500 milliseconds over a rolling 30-day period.
```

Tag:

```text
v8-sre
```

---

### Stage 9 — PostgreSQL High Availability

Build and test a highly available database environment.

Topics:

* Streaming replication
* WAL
* Primary and standby servers
* Synchronous and asynchronous replication
* Connection pooling
* Automated backups
* Point-in-time recovery
* Failover
* RPO and RTO

Tag:

```text
v9-database-ha
```

---

### Stage 10 — Messaging

Integrate asynchronous task processing.

RabbitMQ topics:

* Exchanges
* Queues
* Bindings
* Acknowledgements
* Retries
* Dead-letter queues
* Quorum queues

Kafka topics:

* Topics
* Partitions
* Consumer groups
* Offsets
* Replication
* Retention
* Consumer lag
* Rebalancing

Tag:

```text
v10-messaging
```

---

### Stage 11 — Distributed Storage

Add a Ceph storage laboratory.

Topics:

* MON
* MGR
* OSD
* RADOS
* CRUSH
* Pools
* Placement Groups
* RBD
* CephFS
* Kubernetes CSI
* Disk-failure recovery

Tag:

```text
v11-ceph
```

---

## Repository Structure

```text
production-platform-lab/
├── application/
│   ├── api/
│   ├── worker/
│   ├── frontend/
│   └── tests/
│
├── infrastructure/
│   ├── virtual-machines/
│   ├── cloud-init/
│   ├── networking/
│   └── ansible/
│
├── containers/
│   ├── api.Dockerfile
│   ├── worker.Dockerfile
│   └── docker-compose.yml
│
├── ci/
│   └── gitlab-ci.yml
│
├── kubernetes/
│   ├── base/
│   ├── development/
│   ├── staging/
│   └── production/
│
├── helm/
│   └── taskflow/
│
├── database/
│   ├── migrations/
│   ├── replication/
│   ├── backup/
│   └── restore/
│
├── messaging/
│   ├── rabbitmq/
│   └── kafka/
│
├── observability/
│   ├── prometheus/
│   ├── alertmanager/
│   ├── grafana/
│   └── elastic/
│
├── storage/
│   └── ceph/
│
├── security/
│   ├── rbac/
│   ├── network-policies/
│   └── policies/
│
├── architecture/
├── diagrams/
├── runbooks/
├── postmortems/
├── disaster-recovery/
├── screenshots/
├── .env.example
├── LICENSE
└── README.md
```

---

## Local Laboratory

The complete environment can run on local virtual machines.

Example topology:

```text
Host Machine
│
├── VM 1: GitLab and CI Runner
├── VM 2: Load Balancer
├── VM 3: Kubernetes Control Plane
├── VM 4: Kubernetes Worker
├── VM 5: Kubernetes Worker
└── VM 6: Storage or Monitoring
```

Possible virtualization platforms:

* KVM/QEMU
* Libvirt
* Proxmox
* VirtualBox
* VMware Workstation

Local virtual networks simulate:

* Cloud private networks
* Cloud instances
* Security groups
* Block storage
* Internal DNS
* Public and private services
* Load balancers
* Node failures

---

## Hardware Recommendations

Minimum practical environment:

```text
CPU: 8 logical cores
RAM: 16 GB
Storage: 150 GB SSD
```

Preferred environment:

```text
CPU: 12 or more logical cores
RAM: 32 GB
Storage: 250 GB or more SSD
```

Heavy components do not need to run simultaneously.

For example:

* Run GitLab only while testing CI/CD.
* Run Elasticsearch only while working on centralized logging.
* Run Ceph in a separate laboratory stage.
* Shut down unused virtual machines.

---

## Failure Engineering

The platform is intentionally tested under failure conditions.

Example scenarios:

* Application crash
* Incorrect systemd service
* Full disk
* Memory exhaustion
* DNS failure
* TLS certificate failure
* Closed firewall port
* Unreachable database
* Broken container image
* Failed CI/CD deployment
* Kubernetes node failure
* CrashLoopBackOff
* Broken readiness probe
* NetworkPolicy failure
* PostgreSQL primary failure
* Replica lag
* RabbitMQ queue backlog
* Kafka consumer lag
* Ceph OSD failure
* Monitoring outage

Each incident should produce:

```text
Incident title
Severity
Detection method
User impact
Timeline
Root cause
Mitigation
Resolution
Recovery time
Corrective actions
Runbook updates
```

---

## Documentation

The repository includes operational documentation for:

* Installation
* Configuration
* Deployment
* Upgrades
* Backup
* Restore
* Rollback
* Troubleshooting
* Incident response
* Disaster recovery
* Security
* Capacity planning

Documentation is treated as part of the product rather than an optional addition.

---

## Security Principles

The project follows these practices:

* No secrets committed to Git
* Example configuration files only
* Restricted Linux users
* SSH key authentication
* TLS encryption
* Non-root containers
* Minimal container images
* Kubernetes RBAC
* Network policies
* Resource limits
* Protected CI/CD variables
* Secret rotation procedures
* Backup encryption where applicable

Files containing sensitive values should be excluded through `.gitignore`.

Example files:

```text
.env.example
vault.example.yml
secrets.example.yaml
```

---

## Definition of Done

A project stage is complete only when it includes:

* Working source code
* Repeatable installation
* Automated deployment
* Architecture diagram
* Health checks
* Metrics
* Logs
* Backup procedure
* Restore test
* Security considerations
* Failure simulations
* Troubleshooting runbook
* Clear documentation

Installing a tool alone does not complete a stage.

---

## Final Demonstration

The final demonstration should show the full lifecycle of a software release:

1. Push a code change.
2. Run automated tests.
3. Build a container image.
4. Scan the image.
5. Publish it to the registry.
6. Deploy it to Kubernetes.
7. Generate application traffic.
8. View service metrics in Grafana.
9. Search request logs in Kibana.
10. Introduce a faulty release.
11. Detect the failure through an alert.
12. Investigate the incident.
13. Roll back the deployment.
14. Confirm service recovery.
15. Publish an incident postmortem.

---

## Skills Demonstrated

This project demonstrates practical ability in:

* Infrastructure engineering
* DevOps
* Platform engineering
* Site Reliability Engineering
* Linux administration
* Network troubleshooting
* Infrastructure automation
* CI/CD
* Kubernetes
* Observability
* Database operations
* Messaging systems
* Distributed storage
* Incident response
* Disaster recovery

---

## Project Status

```text
[ ] Stage 1 — Native Linux Service
[ ] Stage 2 — Networking and TLS
[ ] Stage 3 — Ansible Automation
[ ] Stage 4 — Containerization
[ ] Stage 5 — CI/CD
[ ] Stage 6 — Kubernetes
[ ] Stage 7 — Observability
[ ] Stage 8 — SRE Practices
[ ] Stage 9 — PostgreSQL High Availability
[ ] Stage 10 — Messaging
[ ] Stage 11 — Distributed Storage
[ ] Final Production Simulation
```

---

## Disclaimer

This project is a learning and portfolio environment. It simulates production practices but should not be treated as a complete production platform without additional security review, capacity planning, testing, and organizational controls.

---

## License

This project is available under the MIT License.
