# Production Platform Lab

A project-based **DevOps, Platform Engineering, Site Reliability Engineering, and Distributed Systems laboratory** that builds a production-style microservice platform from the ground up.

The project begins with a small application running as native Linux services and gradually evolves into an automated, containerized, observable, highly available Kubernetes platform.

The objective is **not simply to install DevOps tools**.

The objective is to understand how production systems are:

* designed
* implemented
* deployed
* automated
* networked
* secured
* monitored
* scaled
* troubleshot
* recovered from failure

The project intentionally introduces complexity gradually so that every major infrastructure component exists to solve a problem encountered in an earlier stage.

---

# Project Goals

This laboratory is designed to develop practical experience in:

* Linux systems engineering
* TCP/IP networking
* DNS
* HTTP and HTTPS
* TLS
* Go
* Python
* Bash
* Git
* systemd
* Nginx / HAProxy
* PostgreSQL
* Infrastructure as Code
* Ansible
* containers and OCI
* Docker / containerd
* GitLab CI/CD
* Kubernetes
* Helm
* gRPC
* Protocol Buffers
* RabbitMQ
* Kafka
* Prometheus
* Grafana
* centralized logging
* OpenTelemetry
* PostgreSQL high availability
* distributed storage with Ceph
* SLI / SLO engineering
* incident response
* disaster recovery
* failure engineering
* distributed-system design

---

# Application

The infrastructure hosts a distributed task-processing system called **TaskFlow**.

TaskFlow is intentionally small from a product perspective because the main focus of the project is infrastructure and distributed-system behavior.

Users can:

* create tasks
* assign tasks
* update task status
* process tasks asynchronously
* view task-processing results
* query task history
* receive task notifications
* view basic analytics

The application evolves from a single service into a polyglot microservice architecture.

---

# Application Architecture

The final application consists of several independently deployable services.

```text
                              Clients
                                 │
                              HTTPS/REST
                                 │
                                 ▼
                        ┌─────────────────┐
                        │   API Gateway   │
                        │       Go        │
                        └────────┬────────┘
                                 │
                        Internal gRPC
                   ┌─────────────┴─────────────┐
                   │                           │
                   ▼                           ▼
          ┌─────────────────┐         ┌─────────────────┐
          │  Task Service   │         │  User Service   │
          │       Go        │         │       Go        │
          └────────┬────────┘         └────────┬────────┘
                   │                           │
                   ▼                           ▼
                Task DB                     User DB
                   │
                   │ command
                   ▼
              ┌──────────┐
              │ RabbitMQ │
              └────┬─────┘
                   │
                   ▼
         ┌─────────────────────┐
         │ Processing Service  │
         │       Python        │
         └─────────┬───────────┘
                   │
                   │ domain events
                   ▼
              ┌──────────┐
              │  Kafka   │
              └────┬─────┘
                   │
       ┌───────────┼──────────────┐
       │           │              │
       ▼           ▼              ▼
┌────────────┐ ┌───────────┐ ┌──────────────┐
│ Analytics  │ │   Audit   │ │ Notification │
│   Python   │ │    Go     │ │ Go / Python  │
└────────────┘ └───────────┘ └──────────────┘
```

---

# Services

## API Gateway — Go

The public entry point into TaskFlow.

Responsibilities:

* expose REST/HTTPS APIs
* route requests to internal services
* authentication and authorization enforcement
* request validation
* rate limiting
* correlation IDs
* trace propagation
* external API versioning

The gateway communicates with internal services primarily using **gRPC**.

---

## Task Service — Go

Owns the core task domain.

Responsibilities:

* create tasks
* update tasks
* retrieve task state
* manage task lifecycle
* persist task metadata
* publish task-processing commands
* publish task domain events

The Task Service owns its own PostgreSQL data.

---

## User Service — Go

Owns user and identity information.

Responsibilities:

* user accounts
* user metadata
* identity information
* authorization-related data
* ownership relationships

Other services access user information through **gRPC**, rather than reading the User Service database directly.

---

## Processing Service — Python

Consumes asynchronous work from RabbitMQ.

Responsibilities:

* execute background tasks
* process long-running operations
* retry failed operations
* handle transient failures
* emit processing results
* publish completion/failure events

Running this service in Python also allows the laboratory to compare operational characteristics of Go and Python workloads.

---

## Analytics Service — Python

Consumes events from Kafka and produces derived information.

Responsibilities:

* task-processing statistics
* aggregates
* historical analysis
* processing latency reports
* failure-rate analysis
* derived read models

---

## Audit Service — Go

Consumes Kafka events and maintains an immutable operational history.

Example events:

```text
TaskCreated
TaskAssigned
TaskStarted
TaskCompleted
TaskFailed
UserCreated
```

The Audit Service provides a permanent history of important application events.

---

## Notification Service — Go or Python

Consumes relevant Kafka events.

Example responsibilities:

* notify when a task completes
* notify when processing fails
* notify when retry limits are exceeded

---

# Communication Models

The project intentionally uses multiple communication mechanisms because each solves a different distributed-system problem.

## REST / HTTPS

Used for communication between external clients and the API Gateway.

```text
Client
   │
   │ HTTPS / REST
   ▼
API Gateway
```

Topics explored:

* HTTP
* HTTPS
* TLS
* reverse proxies
* request routing
* rate limiting
* API versioning

---

## gRPC

Used for synchronous internal service-to-service communication.

```text
API Gateway
     │
     ├──── gRPC ───► Task Service
     │
     └──── gRPC ───► User Service
```

Topics explored:

* Protocol Buffers
* `.proto` contracts
* code generation
* unary RPCs
* streaming RPCs
* deadlines
* cancellation propagation
* metadata
* status codes
* schema evolution
* service discovery
* load balancing

---

## RabbitMQ

Used for **commands and work distribution**.

```text
Task Service
     │
     │ ProcessTask
     ▼
  RabbitMQ
     │
     ▼
Processing Service
```

Topics explored:

* exchanges
* queues
* bindings
* acknowledgements
* retries
* dead-letter queues
* quorum queues
* prefetch
* consumer crashes
* poison messages
* duplicate delivery
* idempotent consumers
* backpressure

---

## Kafka

Used for **durable domain event streaming**.

```text
Task Service / Worker
        │
        │ TaskCompleted
        ▼
       Kafka
        │
   ┌────┼─────────────┐
   │    │             │
   ▼    ▼             ▼
Audit Analytics Notification
```

Topics explored:

* topics
* partitions
* offsets
* consumer groups
* replication
* retention
* event replay
* ordering
* consumer lag
* rebalancing

---

# Service Data Ownership

Each service owns its own data.

```text
Task Service
     │
     ▼
 task_db


User Service
     │
     ▼
 user_db


Analytics Service
     │
     ▼
 analytics_db
```

Services must not directly modify another service's tables.

Communication occurs through:

* gRPC
* REST where appropriate
* RabbitMQ commands
* Kafka events

Initially, multiple logical databases may run inside the same PostgreSQL cluster.

```text
PostgreSQL Cluster
│
├── task_db
├── user_db
└── analytics_db
```

Logical service ownership is kept separate even when physical infrastructure is shared.

---

# Technology Stack

## Application

```text
Go
Python
Protocol Buffers
gRPC
REST
```

## Data

```text
PostgreSQL
Redis (optional later-stage cache)
```

## Messaging

```text
RabbitMQ
Kafka
```

## Infrastructure

```text
Linux
KVM / QEMU
Libvirt
Ansible
Docker
containerd
Kubernetes
Helm
```

## Networking

```text
TCP/IP
DNS
HTTP
HTTPS
TLS
Nginx
HAProxy
Kubernetes Services
Ingress
NetworkPolicy
```

## Observability

```text
Prometheus
Alertmanager
Grafana
OpenTelemetry
Structured Logging
Centralized Log Storage
```

## Storage

```text
Ceph
RADOS
RBD
CephFS
Kubernetes CSI
```

## CI/CD

```text
Git
GitLab
GitLab Runner
Container Registry
```

---

# Final Platform Architecture

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
                                   │
          ┌────────────────────────┼────────────────────────┐
          │                        │                        │
          ▼                        ▼                        ▼
    API Gateway              Task/User Services       Python Workers
    Go replicas                Go replicas              autoscaled
          │                        │                        │
          └─────────────┬──────────┴───────────┬────────────┘
                        │                      │
                        ▼                      ▼
                PostgreSQL HA             RabbitMQ
                        │                      │
                        │                      ▼
                        │                 Worker Jobs
                        │
                        ▼
                       Kafka
                        │
             ┌──────────┼───────────┐
             │          │           │
             ▼          ▼           ▼
         Analytics    Audit    Notifications
             │
             ▼
      Persistent Storage
             │
             ▼
            Ceph
```

---

# Observability Architecture

Every application service emits:

* metrics
* structured logs
* distributed traces

```text
             Go / Python Services
                     │
          ┌──────────┼───────────┐
          │          │           │
          ▼          ▼           ▼
       Metrics      Logs       Traces
          │          │           │
          ▼          ▼           ▼
     Prometheus   Log Store   OpenTelemetry
          │          │           │
          └──────────┼───────────┘
                     ▼
                   Grafana
```

Important identifiers propagate across service boundaries:

```text
request_id
trace_id
task_id
user_id
```

This allows a single operation to be traced across:

```text
Client
  ↓
Gateway
  ↓
Task Service
  ↓
RabbitMQ
  ↓
Python Worker
  ↓
Kafka
  ↓
Analytics / Audit
```

---

# PostgreSQL High Availability

The final database architecture introduces replication and failover.

```text
                      Applications
                           │
                           ▼
                       PgBouncer
                           │
                           ▼
                      ┌─────────┐
                      │ Primary │
                      └────┬────┘
                           │
                    WAL Replication
                ┌──────────┴──────────┐
                │                     │
                ▼                     ▼
           ┌─────────┐           ┌─────────┐
           │Replica 1│           │Replica 2│
           └─────────┘           └─────────┘
```

Topics:

* WAL
* LSN
* streaming replication
* primary/standby architecture
* synchronous replication
* asynchronous replication
* replication lag
* failover
* connection pooling
* PITR
* RPO
* RTO
* split-brain risks

---

# Project Stages

Each stage solves problems exposed by the previous stage.

---

## Stage 0 — Application Foundation

Build the workload that the rest of the infrastructure will operate.

### Goals

* establish the Go and Python application structure
* build the initial Task Service
* implement health endpoints
* implement structured logging
* implement graceful shutdown
* introduce PostgreSQL persistence
* introduce gRPC
* introduce the first Python worker
* establish service boundaries
* establish testing

### Concepts

* Go networking
* `net/http`
* Go concurrency
* context propagation
* graceful shutdown
* process signals
* Python workers
* service contracts
* Protocol Buffers
* gRPC
* persistence
* stateless services
* service ownership
* synchronous vs asynchronous communication

Tag:

```text
v0-application-foundation
```

---

## Stage 1 — Native Linux Services

Deploy the application directly onto Linux virtual machines.

### Goals

Understand how applications run without containers or orchestration.

### Topics

* Linux users and groups
* permissions
* process lifecycle
* signals
* systemd
* service units
* journald
* Nginx
* PostgreSQL
* SSH
* firewall configuration
* environment files
* privilege separation
* backup scripts
* restore procedures

Tag:

```text
v1-linux-service
```

---

## Stage 2 — Networking and TLS

Separate services across multiple virtual machines.

### Goals

Introduce real network boundaries and partial failures.

### Topics

* IP addressing
* subnetting
* routing
* ARP
* TCP
* UDP
* DNS
* HTTP
* HTTPS
* TLS
* certificate management
* reverse proxies
* load balancing
* gRPC over HTTP/2
* firewall rules
* network troubleshooting
* tcpdump
* ss
* dig
* traceroute

Tag:

```text
v2-networking-tls
```

---

## Stage 3 — Ansible Automation

Automate the manually built Linux environment.

### Topics

* inventories
* playbooks
* roles
* variables
* templates
* handlers
* Ansible Vault
* idempotence
* configuration drift
* service deployment
* secret management

Tag:

```text
v3-ansible
```

---

## Stage 4 — Containerization

Package all application services as OCI containers.

### Topics

* Dockerfiles
* multi-stage builds
* Go static binaries
* Python runtime images
* OCI images
* layers
* container networking
* namespaces
* cgroups
* persistent volumes
* health checks
* non-root containers
* resource limits
* Docker Compose
* image security
* image-size comparison

Tag:

```text
v4-containers
```

---

## Stage 5 — CI/CD

Build an automated software delivery pipeline.

```text
Lint
  ↓
Unit Tests
  ↓
Integration Tests
  ↓
gRPC Contract Tests
  ↓
Build Container Images
  ↓
Security Scan
  ↓
Push Images
  ↓
Deploy to Staging
  ↓
Smoke Tests
  ↓
Manual Approval
  ↓
Deploy to Production
```

### Topics

* GitLab Runner
* pipeline jobs
* artifacts
* caching
* container registry
* protected variables
* semantic versioning
* deployment strategies
* rollback
* image scanning
* supply-chain security

Tag:

```text
v5-cicd
```

---

## Stage 6 — Kubernetes

Move the platform onto a multi-node Kubernetes cluster.

### Topics

* Pods
* Deployments
* StatefulSets
* Services
* Ingress
* ConfigMaps
* Secrets
* PersistentVolumes
* resource requests
* resource limits
* liveness probes
* readiness probes
* startup probes
* Horizontal Pod Autoscaler
* RBAC
* ServiceAccounts
* NetworkPolicy
* Pod disruption
* rolling updates
* Helm
* cluster upgrades
* etcd backups
* Kubernetes service discovery
* gRPC load balancing

Tag:

```text
v6-kubernetes
```

---

## Stage 7 — Observability

Build metrics, logs, traces, dashboards, and alerts.

### Topics

* Prometheus
* PromQL
* Grafana
* Alertmanager
* Node Exporter
* application metrics
* Go metrics
* Python metrics
* PostgreSQL metrics
* RabbitMQ metrics
* Kafka metrics
* Kubernetes metrics
* structured logs
* centralized logging
* correlation IDs
* OpenTelemetry
* distributed tracing
* RED method
* USE method

Tag:

```text
v7-observability
```

---

## Stage 8 — SRE Practices

Define and operate explicit reliability objectives.

### Topics

* SLI
* SLO
* SLA
* error budgets
* burn-rate alerts
* incident severity
* on-call procedures
* runbooks
* postmortems
* toil reduction
* capacity planning

Example SLO:

```text
99.9% of valid API requests must complete successfully
within 500 milliseconds over a rolling 30-day period.
```

Tag:

```text
v8-sre
```

---

## Stage 9 — PostgreSQL High Availability

Build and test a highly available relational database environment.

### Topics

* WAL
* streaming replication
* primary and standby servers
* synchronous replication
* asynchronous replication
* connection pooling
* automated backups
* PITR
* failover
* replica promotion
* replication lag
* RPO
* RTO
* failure detection

Tag:

```text
v9-database-ha
```

---

## Stage 10 — Distributed Messaging

Deepen the RabbitMQ and Kafka architecture.

### RabbitMQ

* exchanges
* queues
* bindings
* acknowledgements
* retries
* dead-letter queues
* quorum queues
* prefetch
* publisher confirms
* idempotency
* duplicate delivery
* poison messages

### Kafka

* topics
* partitions
* consumer groups
* offsets
* replication
* retention
* ordering
* consumer lag
* rebalancing
* replay
* delivery semantics

Tag:

```text
v10-messaging
```

---

## Stage 11 — Distributed Storage

Add Ceph as the distributed storage layer.

### Topics

* MON
* MGR
* OSD
* RADOS
* CRUSH
* pools
* placement groups
* replication
* failure domains
* RBD
* CephFS
* Kubernetes CSI
* rebalancing
* disk-failure recovery

Tag:

```text
v11-ceph
```

---

# Failure Engineering

The platform is intentionally broken throughout the project.

Failure scenarios include:

### Application

* process crash
* panic
* Python worker crash
* memory exhaustion
* CPU saturation
* bad configuration
* invalid deployment

### Network

* DNS failure
* packet loss
* high latency
* blocked port
* TLS failure
* unreachable dependency
* gRPC deadline exceeded
* connection reset

### Kubernetes

* Pod crash
* CrashLoopBackOff
* failed readiness probe
* node failure
* broken Service
* bad NetworkPolicy
* resource exhaustion

### PostgreSQL

* primary failure
* replication lag
* full disk
* broken connection pool
* failed backup
* failover failure

### RabbitMQ

* consumer failure
* queue backlog
* poison message
* lost acknowledgement
* duplicate delivery

### Kafka

* broker failure
* consumer lag
* rebalance
* unavailable partition
* duplicate processing

### Ceph

* OSD failure
* degraded placement groups
* node loss
* recovery/rebalancing

### Observability

* Prometheus unavailable
* missing metrics
* broken alerts
* centralized logging outage

---

# Incident Documentation

Each significant failure should produce an incident report containing:

```text
Incident title
Severity
Detection method
User impact
Timeline
Root cause
Contributing factors
Mitigation
Resolution
Recovery time
Corrective actions
Runbook updates
```

The goal is not only to repair failures, but to improve the system after each incident.

---

# Repository Structure

```text
production-platform-lab/
│
├── application/
│   │
│   ├── gateway/
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── go.mod
│   │
│   ├── task-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── go.mod
│   │
│   ├── user-service/
│   │   ├── cmd/
│   │   ├── internal/
│   │   └── go.mod
│   │
│   ├── processing-service/
│   │   ├── src/
│   │   ├── tests/
│   │   └── requirements.txt
│   │
│   ├── analytics-service/
│   │   └── ...
│   │
│   ├── audit-service/
│   │   └── ...
│   │
│   └── notification-service/
│       └── ...
│
├── proto/
│   ├── task/
│   └── user/
│
├── infrastructure/
│   ├── virtual-machines/
│   ├── cloud-init/
│   ├── networking/
│   └── ansible/
│
├── containers/
│   ├── gateway.Dockerfile
│   ├── task-service.Dockerfile
│   ├── processing-service.Dockerfile
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
│   ├── logging/
│   └── opentelemetry/
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
│
├── .env.example
├── .gitignore
├── LICENSE
└── README.md
```

---

# Local Laboratory

The complete environment is designed to run locally using Linux virtual machines.

Example topology:

```text
Host Machine
│
├── VM 1: GitLab / CI Runner / Registry
├── VM 2: Load Balancer
├── VM 3: Kubernetes Control Plane
├── VM 4: Kubernetes Worker
├── VM 5: Kubernetes Worker
└── VM 6: Monitoring / Storage
```

Possible virtualization platforms:

* KVM/QEMU
* Libvirt
* Proxmox
* VirtualBox
* VMware Workstation

The virtual environment simulates:

* cloud private networks
* VM instances
* security groups
* block storage
* internal DNS
* load balancers
* private services
* public services
* node failures
* network failures

Heavy components do not need to run simultaneously.

For example:

* run GitLab only while working on CI/CD
* run heavy centralized logging only during observability labs
* run Ceph as a dedicated stage
* shut down unused virtual machines

---

# Security Principles

The project follows the following security practices:

* no secrets committed to Git
* SSH key authentication
* least-privilege Linux users
* service-specific users
* TLS encryption
* non-root containers
* minimal container images
* Kubernetes RBAC
* Kubernetes NetworkPolicies
* resource limits
* protected CI/CD variables
* secret rotation procedures
* backup encryption where appropriate
* isolated service credentials
* independent database credentials
* dependency and image scanning

Example configuration committed to Git:

```text
.env.example
vault.example.yml
secrets.example.yaml
```

Real secrets remain excluded.

---

# Definition of Done

A project stage is considered complete only when it contains:

* working source code
* repeatable installation
* automated deployment where applicable
* architecture documentation
* health checks
* metrics
* logs
* traces where applicable
* backup procedure
* successful restore test
* security considerations
* failure simulations
* troubleshooting runbook
* clear documentation

Installing a tool alone does **not** complete a stage.

---

# Final Production Simulation

The final demonstration should exercise the complete lifecycle.

```text
1. Push a source-code change.

2. Run automated tests.

3. Run contract and integration tests.

4. Build Go and Python container images.

5. Scan the images.

6. Publish images to the registry.

7. Deploy to staging.

8. Run smoke tests.

9. Promote the release.

10. Deploy to Kubernetes.

11. Generate application traffic.

12. Observe metrics in Grafana.

13. Trace a request across multiple services.

14. Inspect distributed logs.

15. Observe RabbitMQ and Kafka activity.

16. Introduce a faulty release.

17. Detect the problem using alerts.

18. Investigate the incident.

19. Roll back the deployment.

20. Confirm service recovery.

21. Simulate infrastructure failure.

22. Verify failover/recovery.

23. Publish an incident postmortem.
```

---

# Distributed Systems Concepts Demonstrated

The laboratory intentionally exposes important distributed-system problems including:

* partial failure
* unreliable networks
* timeouts
* retries
* exponential backoff
* jitter
* duplicate requests
* idempotency
* eventual consistency
* replication
* consistency trade-offs
* message ordering
* consumer coordination
* queue backpressure
* event replay
* failure detection
* service discovery
* load balancing
* distributed tracing
* graceful degradation
* horizontal scaling
* leader/follower database architecture
* stateful vs stateless workloads
* data ownership
* service boundaries
* fault isolation

---

# Skills Demonstrated

This project demonstrates practical ability in:

* Software Engineering
* Go
* Python
* Linux Engineering
* DevOps
* Platform Engineering
* Site Reliability Engineering
* Distributed Systems
* TCP/IP Networking
* Infrastructure Automation
* Infrastructure as Code
* CI/CD
* Containers
* Kubernetes
* gRPC
* Event-Driven Architecture
* RabbitMQ
* Kafka
* PostgreSQL
* Database High Availability
* Observability
* OpenTelemetry
* Distributed Storage
* Security Engineering
* Incident Response
* Disaster Recovery
* Failure Engineering

---

# Project Status

```text
[ ] Stage 0  — Application Foundation
[ ] Stage 1  — Native Linux Services
[ ] Stage 2  — Networking and TLS
[ ] Stage 3  — Ansible Automation
[ ] Stage 4  — Containerization
[ ] Stage 5  — CI/CD
[ ] Stage 6  — Kubernetes
[ ] Stage 7  — Observability
[ ] Stage 8  — SRE Practices
[ ] Stage 9  — PostgreSQL High Availability
[ ] Stage 10 — Distributed Messaging
[ ] Stage 11 — Distributed Storage
[ ] Final Production Simulation
```

---

# Engineering Philosophy

This repository intentionally avoids treating infrastructure technologies as isolated tools.

The project follows a progression:

```text
Understand manually
        ↓
Operate manually
        ↓
Observe failure
        ↓
Automate
        ↓
Containerize
        ↓
Orchestrate
        ↓
Measure
        ↓
Scale
        ↓
Break
        ↓
Recover
        ↓
Improve
```

The emphasis is not on memorizing commands.

The emphasis is on understanding **why production systems behave the way they do**.

---

# Disclaimer

This repository is an educational and portfolio laboratory.

It simulates production engineering practices but should not be treated as a complete production platform without additional:

* threat modeling
* security review
* capacity planning
* compliance controls
* load testing
* chaos testing
* organizational processes
* production-specific operational controls

---

# License

MIT License.
