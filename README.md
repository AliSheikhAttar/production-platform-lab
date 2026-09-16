# Production Platform Lab

A project-based **DevOps, Platform Engineering, Site Reliability Engineering, and Distributed Systems laboratory** that builds a production-style polyglot microservice platform from the ground up.

The project begins with small application services running directly on Linux and gradually evolves into an automated, containerized, cloud-provisioned, GitOps-delivered, observable, highly available Kubernetes platform.

The objective is **not simply to install DevOps technologies**.

The objective is to understand how production systems are:

* designed
* implemented
* deployed
* networked
* automated
* secured
* monitored
* scaled
* debugged
* optimized
* recovered from failure

Complexity is introduced gradually so that each technology solves or demonstrates a concrete production problem.

---

# Project Goals

This laboratory develops hands-on experience with:

## Software Engineering

* Go
* Python
* Bash
* SQL
* Git
* REST APIs
* gRPC
* Protocol Buffers
* concurrency
* testing
* service contracts

## Linux & Systems

* Linux administration
* processes
* signals
* systemd
* journald
* users and permissions
* filesystems
* resource management
* cgroups
* namespaces

## Networking

* TCP/IP
* UDP
* DNS
* HTTP
* HTTP/2
* HTTPS
* TLS
* reverse proxies
* load balancing
* service discovery
* network troubleshooting

## Infrastructure & DevOps

* KVM/QEMU (enough to run the lab, not a target of mastery)
* Libvirt
* Ansible
* Terraform / OpenTofu
* a real cloud provider account (AWS, GCP, or Azure)
* Docker
* OCI containers
* containerd
* GitLab CI/CD
* container registries
* Kubernetes
* Helm
* GitOps (ArgoCD or Flux)
* Infrastructure as Code, applied for real — not just read about

## Data & Distributed Systems

* PostgreSQL
* Redis
* RabbitMQ
* Kafka
* Ceph (optional / stretch)
* replication
* caching
* sharding
* consistency
* messaging
* event streaming
* distributed storage

## Observability & Reliability

* Prometheus
* Grafana
* Alertmanager
* OpenTelemetry
* centralized logging
* distributed tracing
* SLI
* SLO
* SLA
* error budgets
* incident response
* chaos engineering
* disaster recovery
* failure engineering

---

# Core Path vs. Stretch Path

This project is deliberately larger than any single career transition requires. Treat it as two paths, not one long checklist.

**Core path — Stage 0 through Stage 9.**
This is where almost all of the signal relevant to an SRE or platform engineering interview lives: the application itself, networking, real infrastructure-as-code against a cloud provider, containers, CI/CD, Kubernetes, GitOps, observability, and SRE practices — including deliberate failure injection and postmortems. Complete this path at solid depth before spending significant time further down.

**Stretch path — Stage 10 through Stage 13.**
PostgreSQL HA, Redis HA, deep messaging internals, and distributed storage extend the same distributed-systems muscles into more specialized territory. Valuable, but not required to be hireable, and safe to do unevenly, partially, or skip under time pressure. Stage 13 (distributed storage) is explicitly optional — see that stage for why.

If you have a few months rather than a year, finishing the core path well and picking one or two stretch stages to go deep on will produce a stronger portfolio than spreading thin across all fourteen stages.

---

# Application

The infrastructure hosts a distributed task-processing application called **TaskFlow**.

TaskFlow is intentionally simple from a product perspective because the primary purpose of the project is infrastructure and distributed-system engineering.

Users can:

* create tasks
* assign tasks
* update task status
* submit tasks for asynchronous processing
* inspect processing results
* query task history
* receive notifications
* view basic analytics

The application begins as a simple service and gradually evolves into a polyglot microservice system.

---

# Final Application Architecture

```text
                                   Clients
                                      │
                               HTTPS / REST
                                      │
                                      ▼
                           ┌────────────────────┐
                           │    API Gateway     │
                           │        Go          │
                           └─────────┬──────────┘
                                     │
                        ┌────────────┼────────────┐
                        │            │            │
                        │ gRPC       │ Redis      │ gRPC
                        ▼            ▼            ▼
                ┌──────────────┐ ┌─────────┐ ┌──────────────┐
                │ Task Service │ │  Redis  │ │ User Service │
                │      Go      │ │ Cluster │ │      Go      │
                └──────┬───────┘ └─────────┘ └──────┬───────┘
                       │                              │
                       ▼                              ▼
                   task_db                         user_db
                       │
                       │ ProcessTask command
                       ▼
                 ┌────────────┐
                 │  RabbitMQ  │
                 └─────┬──────┘
                       │
             ┌─────────┴──────────┐
             ▼                    ▼
     ┌───────────────┐    ┌───────────────┐
     │ Python Worker │    │ Python Worker │
     │       #1      │    │       #N      │
     └───────┬───────┘    └───────┬───────┘
             │                    │
             └──────────┬─────────┘
                        │
                        │ domain events
                        ▼
                  ┌───────────┐
                  │   Kafka   │
                  └─────┬─────┘
                        │
             ┌──────────┼────────────┐
             │          │            │
             ▼          ▼            ▼
       ┌───────────┐ ┌────────┐ ┌──────────────┐
       │ Analytics │ │ Audit  │ │ Notification │
       │  Python   │ │   Go   │ │  Go/Python   │
       └───────────┘ └────────┘ └──────────────┘
```

---

# Service Responsibilities

## API Gateway — Go

The public entry point to TaskFlow.

Responsibilities:

* expose REST/HTTPS APIs
* route requests
* validate authentication
* enforce authorization
* rate limiting
* request validation
* correlation IDs
* trace propagation
* API versioning
* request timeouts

Internal communication primarily uses **gRPC**.

Redis is used for shared rate-limit state and selected short-lived data.

---

## Task Service — Go

Owns the task domain.

Responsibilities:

* create tasks
* retrieve tasks
* update tasks
* manage task lifecycle
* persist task metadata
* submit processing commands
* publish domain events
* cache selected read-heavy data

The Task Service owns `task_db`.

Redis is used through a **cache-aside strategy** for selected hot reads.

PostgreSQL remains the source of truth.

---

## User Service — Go

Owns users and identity-related information.

Responsibilities:

* user accounts
* user metadata
* identity information
* task ownership relationships
* authorization-related data

The User Service owns `user_db`.

Other services access user information through APIs rather than directly accessing its database.

---

## Processing Service — Python

Consumes asynchronous tasks from RabbitMQ.

Responsibilities:

* execute background jobs
* perform long-running processing
* acknowledge successfully completed work
* retry transient failures
* reject unrecoverable work
* produce processing results
* publish completion/failure events

Multiple workers can process jobs concurrently.

---

## Analytics Service — Python

Consumes events from Kafka.

Responsibilities:

* processing statistics
* historical analytics
* latency calculations
* failure-rate calculations
* aggregations
* reporting
* derived read models

---

## Audit Service — Go

Consumes domain events and builds an immutable audit history.

Example events:

```text
TaskCreated
TaskAssigned
TaskStarted
TaskCompleted
TaskFailed
UserCreated
```

---

## Notification Service — Go or Python

Consumes relevant Kafka events.

Responsibilities may include:

* task-completed notification
* task-failed notification
* retry exhaustion alerts
* user-facing event notifications

---

# Redis

Redis is a **first-class technology in this laboratory**.

It is not used as the primary source of truth.

PostgreSQL remains authoritative for durable relational application state.

Redis is introduced to study high-performance distributed in-memory systems.

## Application Uses

### Caching

```text
Client
   │
   ▼
Task Service
   │
   ├──── cache lookup ───► Redis
   │                         │
   │                     cache hit
   │                         │
   │◄────────────────────────┘
   │
   │ cache miss
   ▼
PostgreSQL
   │
   ▼
populate Redis
```

This demonstrates the **cache-aside pattern**.

---

### Distributed Rate Limiting

```text
Clients
   │
   ▼
API Gateway
   │
   ▼
Redis
   │
   ├── request counters
   ├── TTL windows
   └── shared rate-limit state
```

Multiple gateway replicas share the same rate-limiting state.

---

### Temporary Data

Redis can also hold:

* short-lived tokens
* transient metadata
* idempotency keys
* temporary request state

---

# Redis Topics

The Redis stage investigates:

* strings
* hashes
* lists
* sets
* sorted sets
* key expiration
* TTL
* memory management
* eviction policies
* cache-aside
* cache invalidation
* cache stampede
* cache penetration
* hot keys
* pipelining
* transactions
* Lua scripting
* optimistic locking
* persistence
* RDB snapshots
* AOF
* replication
* Sentinel
* Redis Cluster
* hash slots
* sharding
* failover
* split-brain considerations
* distributed locks
* Redlock discussion and limitations
* Pub/Sub
* Redis Streams
* monitoring
* latency analysis

Redis Pub/Sub and Streams are studied partly to compare them with RabbitMQ and Kafka.

---

# Communication Models

The project deliberately uses several communication mechanisms.

Each demonstrates a different distributed-system model.

---

## REST / HTTPS

Used between external clients and the API Gateway.

```text
Client
   │
   │ HTTPS / REST
   ▼
API Gateway
```

Topics:

* HTTP
* REST
* TLS
* API design
* reverse proxies
* rate limiting
* authentication
* versioning

---

## gRPC

Used for synchronous internal communication.

```text
API Gateway
     │
     ├──── gRPC ───► Task Service
     │
     └──── gRPC ───► User Service
```

Topics:

* Protocol Buffers
* `.proto` files
* code generation
* unary RPC
* server streaming
* client streaming
* bidirectional streaming
* HTTP/2
* deadlines
* cancellation
* metadata
* status codes
* interceptors
* schema evolution
* service discovery
* load balancing

---

## RabbitMQ

Used for **commands and task distribution**.

```text
Task Service
     │
     │ ProcessTask
     ▼
  RabbitMQ
     │
     ▼
Python Workers
```

RabbitMQ represents:

> work that should be performed

Topics:

* exchanges
* queues
* bindings
* routing keys
* acknowledgements
* negative acknowledgements
* publisher confirms
* retry strategies
* dead-letter queues
* quorum queues
* prefetch
* backpressure
* poison messages
* duplicate delivery
* idempotent consumers

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
   ┌────┼───────────┐
   ▼    ▼           ▼
Audit Analytics Notification
```

Kafka represents:

> facts that have already happened

Topics:

* topics
* partitions
* brokers
* producers
* consumers
* consumer groups
* offsets
* replication
* ISR
* retention
* compaction
* ordering
* event replay
* consumer lag
* rebalancing
* delivery semantics

---

# Service Data Ownership

Each service owns its data.

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

Services must not directly manipulate another service's private tables.

Communication happens through:

* gRPC
* REST
* RabbitMQ
* Kafka

Initially, logical databases may share one PostgreSQL cluster:

```text
PostgreSQL Cluster
│
├── task_db
├── user_db
└── analytics_db
```

Logical ownership remains separate even when physical infrastructure is shared.

---

# Technology Stack

## Application

```text
Go
Python
REST
gRPC
Protocol Buffers
```

## Data

```text
PostgreSQL
Redis
```

## Messaging

```text
RabbitMQ
Kafka
```

## Infrastructure

```text
Linux
KVM / QEMU (lab only)
Libvirt (lab only)
Ansible
Terraform / OpenTofu
A major cloud provider (AWS, GCP, or Azure — pick one)
Docker
OCI
containerd
Kubernetes
Helm
ArgoCD or Flux
```

## Networking

```text
TCP/IP
UDP
DNS
HTTP
HTTP/2
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
Centralized Logging
Distributed Tracing
```

## Storage (optional / stretch)

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
                     ┌────────────┴────────────┐
                     │                         │
              Application Code           Infrastructure Code
                     │                         │
                     ▼                         ▼
               CI/CD Pipeline               Terraform
          ┌────────┼────────┐                 │
          │        │        │                 ▼
        Tests    Build    Scan          Cloud Provider
                   │                    (network, compute,
                   ▼                     Kubernetes cluster)
          Container Registry                   │
                   │                            │
                   ▼                            │
         Update GitOps Manifests                │
                   │                            │
                   ▼                            ▼
             ArgoCD / Flux ─────────► Kubernetes Cluster
             (reconciles cluster
              state from Git)
                                  │
 ┌────────────────────────────────┼────────────────────────────────┐
 │                                │                                │
 ▼                                ▼                                ▼
Gateway                       Go Services                    Python Workers
replicas                       replicas                       autoscaled
 │                                │                                │
 └───────────────┬────────────────┼────────────────┬───────────────┘
                 │                │                │
                 ▼                ▼                ▼
              Redis          PostgreSQL HA     RabbitMQ
                                │
                                ▼
                              Kafka
                                │
                    ┌───────────┼───────────┐
                    ▼           ▼           ▼
                Analytics     Audit    Notifications
                                │
                                ▼
                    Ceph (optional / stretch)
```

---

# Observability Architecture

Every production service emits:

```text
metrics
logs
traces
```

```text
                Go / Python Services
                        │
             ┌──────────┼──────────┐
             │          │          │
             ▼          ▼          ▼
          Metrics      Logs      Traces
             │          │          │
             ▼          ▼          ▼
        Prometheus   Log Store  OpenTelemetry
             │          │          │
             └──────────┼──────────┘
                        ▼
                      Grafana
```

Important identifiers propagate through service boundaries:

```text
request_id
trace_id
task_id
user_id
```

A request can therefore be traced across:

```text
Client
  ↓
Gateway
  ↓
Task Service
  ↓
Redis / PostgreSQL
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
* MVCC
* streaming replication
* synchronous replication
* asynchronous replication
* replica lag
* connection pooling
* backups
* PITR
* promotion
* failover
* RPO
* RTO
* split-brain risks

---

# Redis High Availability

Redis evolves from a single node into fault-tolerant configurations.

## Replication

```text
                   Primary
                     │
              ┌──────┴──────┐
              ▼             ▼
           Replica 1     Replica 2
```

---

## Sentinel

```text
               Sentinel Nodes
              /       |       \
             /        |        \
        Primary ─── Replica ─── Replica
```

Topics:

* health monitoring
* leader election
* automatic failover
* client discovery

---

## Redis Cluster

```text
          Redis Cluster

       Node A        Node B
      Slots 0-X     Slots X-Y
         │             │
      Replica A      Replica B

              Node C
             Slots Y-Z
                │
             Replica C
```

Topics:

* horizontal partitioning
* hash slots
* resharding
* replica promotion
* node failure
* cluster topology

---

# Project Stages

Each stage builds on the problems and knowledge developed in earlier stages. See "Core Path vs. Stretch Path" above for how to prioritize these under time pressure.

---

# Stage 0 — Application Foundation

## Goal

Build the workload that the infrastructure will later operate.

## Contribution to the Project

This stage establishes the services, communication contracts, process behavior, and application-level failure modes that later infrastructure stages must deploy, automate, observe, and recover.

## Topics

* Go service architecture
* Python service architecture
* REST
* gRPC
* Protocol Buffers
* Go `net/http`
* Go concurrency
* `context.Context`
* graceful shutdown
* Linux signals
* configuration
* structured logging
* PostgreSQL persistence
* basic Redis integration
* RabbitMQ integration
* service ownership
* health endpoints
* unit tests
* integration tests

Tag:

```text
v0-application-foundation
```

---

# Stage 1 — Native Linux Services

## Goal

Operate the application directly on Linux without containers.

## Contribution

Establish the operational baseline that later stages automate and abstract.

## Topics

* users
* groups
* file permissions
* processes
* signals
* systemd
* systemd dependencies
* journald
* environment files
* service accounts
* SSH
* PostgreSQL
* Redis
* Nginx
* firewall configuration
* backup scripts
* restore procedures

Tag:

```text
v1-linux-service
```

---

# Stage 2 — Networking and TLS

## Goal

Move services across multiple machines and introduce real network boundaries.

## Contribution

Transforms local function/process communication into distributed communication subject to network failures.

## Topics

* IP addressing
* subnetting
* ARP
* routing
* TCP
* UDP
* DNS
* HTTP
* HTTP/2
* HTTPS
* TLS
* certificates
* gRPC networking
* reverse proxies
* load balancing
* firewalls
* network debugging
* `ss`
* `ip`
* `dig`
* `tcpdump`
* `traceroute`

Tag:

```text
v2-networking-tls
```

---

# Stage 3 — Infrastructure as Code

## Goal

Convert manually configured infrastructure into reproducible, version-controlled automation — both configuration (Ansible) and provisioning (Terraform), against a real cloud provider.

## Contribution

Eliminates configuration drift, makes environments reproducible, and introduces the provisioning workflow (plan → review → apply → state) that most SRE and platform roles are actually hired to operate.

## Topics

### Configuration Management (Ansible)

* inventories
* playbooks
* roles
* variables
* templates
* handlers
* Ansible Vault
* idempotence
* configuration drift
* secret management

### Provisioning (Terraform / OpenTofu)

* providers
* resources
* variables and outputs
* state
* remote state and state locking
* modules
* workspaces
* plan / apply / destroy
* drift detection
* importing existing infrastructure
* provisioning a small footprint on a real cloud provider (network, compute, and — later, in Stage 6 — a Kubernetes cluster)

This is the first stage where the lab stops being purely local. A single small node pool on a free-tier cloud account, provisioned and destroyed on demand, is enough — it does not need to run continuously.

Tag:

```text
v3-infrastructure-as-code
```

---

# Stage 4 — Containerization

## Goal

Package services and dependencies as OCI containers.

## Contribution

Creates portable application artifacts suitable for CI/CD and Kubernetes.

## Topics

* OCI
* Dockerfiles
* image layers
* multi-stage builds
* Go static binaries
* Python runtime images
* namespaces
* cgroups
* container networking
* persistent volumes
* non-root users
* resource limits
* health checks
* Docker Compose
* container security

Tag:

```text
v4-containers
```

---

# Stage 5 — CI/CD

## Goal

Automate build, validation, packaging, and deployment.

## Pipeline

```text
Lint
  ↓
Unit Tests
  ↓
Integration Tests
  ↓
gRPC Contract Tests
  ↓
Build Images
  ↓
Security Scan
  ↓
Push Images
  ↓
Update GitOps Manifest (staging)
  ↓
Smoke Tests
  ↓
Approval
  ↓
Update GitOps Manifest (production)
```

At this stage, "deploy" means committing a new image tag to a manifest repository — not pushing directly to the cluster. The actual deployment mechanism, a controller reconciling that manifest into the cluster, is covered in Stage 7.

## Topics

* GitLab Runner
* pipeline jobs
* artifacts
* caching
* registries
* protected variables
* versioning
* image scanning
* software supply chain
* deployment strategies
* rollback

Tag:

```text
v5-cicd
```

---

# Stage 6 — Kubernetes

## Goal

Move the platform onto a multi-node orchestrated environment.

## Contribution

Introduces reconciliation, declarative infrastructure, scheduling, service discovery, self-healing, and horizontal scaling.

## Topics

* Pods
* ReplicaSets
* Deployments
* StatefulSets
* Services
* Ingress
* ConfigMaps
* Secrets
* PersistentVolumes
* StorageClasses
* requests
* limits
* probes
* HPA
* RBAC
* ServiceAccounts
* NetworkPolicy
* rolling deployments
* Pod disruption
* Kubernetes DNS
* gRPC service discovery
* Helm
* etcd
* cluster upgrades
* etcd backups

Tag:

```text
v6-kubernetes
```

---

# Stage 7 — GitOps & Progressive Delivery

## Goal

Replace direct `kubectl apply` / `helm upgrade` from the pipeline with a Git-driven, reconciled deployment model.

## Contribution

This is close to the default deployment pattern in real platform and SRE roles: the CI pipeline builds, tests, and pushes an image, then updates a manifest in a Git repository; a controller running in the cluster continuously reconciles actual state to match that repository, rather than being told to deploy once. It also gives you an audit trail (every deployment is a commit) and a trivial rollback path (revert the commit).

## Topics

* declarative vs. imperative deployment
* the reconciliation loop
* ArgoCD or Flux (pick one — the pattern matters more than the specific tool)
* application manifests vs. Helm charts vs. Kustomize overlays
* environment promotion (dev → staging → production) via branches or directories
* drift detection and self-healing
* automated sync vs. manual approval gates
* progressive delivery: canary or blue/green rollouts
* automated rollback on failed health checks
* secrets in a GitOps world (sealed-secrets, external-secrets, or a vault integration)

## Failure Experiments

* deploy a broken manifest and watch the controller's health checks catch it
* manually edit a live resource (`kubectl edit`) and watch GitOps revert the drift
* force a rollback via Git revert instead of `kubectl rollout undo`

Tag:

```text
v7-gitops
```

---

# Stage 8 — Observability

## Goal

Make the distributed platform measurable and debuggable.

## Topics

* Prometheus
* PromQL
* Grafana
* Alertmanager
* Node Exporter
* Go metrics
* Python metrics
* PostgreSQL metrics
* Redis metrics
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
v8-observability
```

---

# Stage 9 — SRE Practices

## Goal

Turn monitoring into measurable reliability engineering, and formalize the failure injection that's been happening informally since Stage 0.

## Topics

* SLIs
* SLOs
* SLAs
* error budgets
* burn-rate alerts
* incident severity
* on-call procedures
* runbooks
* postmortems
* toil
* capacity planning
* chaos engineering (Chaos Mesh or Litmus) — turn manual "kill the primary" experiments into scheduled, automated, and progressively more targeted fault injection
* game days — a scheduled exercise where you page yourself, work from the runbook, and time your own response

Example SLO:

```text
99.9% of valid API requests must complete successfully
within 500 milliseconds over a rolling 30-day period.
```

## Portfolio Note

This stage produces your strongest interview material. Publish two or three of your best postmortems as short write-ups — a blog post, or a docs page in this repo: here's the failure you injected, here's what confused you during triage, here's what you changed afterward. That artifact demonstrates judgment under uncertainty directly, rather than leaving an interviewer to infer it from a list of technologies.

Tag:

```text
v9-sre
```

---

# Stage 10 — PostgreSQL High Availability *(stretch)*

## Goal

Understand durable state, replication, database failure, and recovery.

## Topics

* WAL
* LSN
* MVCC
* streaming replication
* primary/standby
* synchronous replication
* asynchronous replication
* PgBouncer
* backup automation
* PITR
* replica promotion
* failover
* replication lag
* RPO
* RTO

Tag:

```text
v10-postgresql-ha
```

---

# Stage 11 — Redis & Distributed Caching *(stretch)*

## Goal

Understand distributed caching, in-memory data systems, replication, sharding, and cache failure modes.

## Contribution

Adds a low-latency distributed state layer and exposes caching-specific consistency and availability problems.

## Topics

### Data Structures

* strings
* hashes
* lists
* sets
* sorted sets
* streams

### Caching

* cache-aside
* TTL
* invalidation
* eviction
* LRU/LFU concepts
* cache stampede
* cache penetration
* hot keys
* stale data
* cache warming

### Coordination

* atomic operations
* transactions
* optimistic locking
* Lua scripting
* idempotency keys
* rate limiting
* distributed locks
* Redlock analysis

### Persistence

* RDB
* AOF
* fsync strategies
* restart recovery

### High Availability

* replication
* Sentinel
* automatic failover
* Redis Cluster
* hash slots
* sharding
* resharding
* replica promotion

### Failure Experiments

* primary failure
* replica lag
* stale cache
* network partition
* cache flush
* memory exhaustion
* eviction storm
* hot-key overload
* cache stampede

Tag:

```text
v11-redis
```

---

# Stage 12 — Distributed Messaging *(stretch)*

## Goal

Study queues and distributed event logs in depth.

## RabbitMQ Topics

* exchanges
* queues
* routing keys
* bindings
* acknowledgements
* publisher confirms
* retry
* dead-letter queues
* quorum queues
* prefetch
* consumer crashes
* poison messages
* duplicate delivery
* idempotency
* backpressure

## Kafka Topics

* brokers
* topics
* partitions
* producers
* consumers
* consumer groups
* offsets
* replication
* ISR
* retention
* compaction
* ordering
* replay
* consumer lag
* rebalancing
* delivery semantics

## Comparison

The project explicitly compares:

```text
Redis Pub/Sub
Redis Streams
RabbitMQ
Kafka
```

to understand why these technologies solve different problems.

## If Time-Constrained

Both systems are worth understanding conceptually, and the comparison itself is a useful distributed-systems lesson. If you need to cut scope, take one to full operational depth — Kafka is the more common requirement in current postings — and keep the other to its core topics plus a single deliberate failure experiment, rather than working the full topic list for both.

Tag:

```text
v12-messaging
```

---

# Stage 13 — Distributed Storage *(optional / stretch)*

## Goal

Get hands-on with distributed storage internals — not to master Ceph operations, but to understand what a managed cloud block/object store is abstracting away from you.

## Why This Is Optional

Most companies run on managed cloud storage (EBS/S3, Persistent Disk/GCS, Azure Disk/Blob) rather than operating their own Ceph clusters, and hands-on Ceph administration appears in relatively few SRE postings outside storage-specialist or on-prem-heavy roles. Full CRUSH-map tuning and placement-group rebalancing are a large time investment for a comparatively narrow slice of the job market.

## Recommended Scope

Do a conceptual pass rather than the full topic list: bring up a small Ceph cluster, mount it through the Kubernetes CSI driver, fail one OSD, and watch recovery happen. That's enough to talk about distributed storage intelligently in an interview.

## Topics (recommended scope)

* Ceph MON / MGR / OSD, at a conceptual level
* RADOS and CRUSH, at a conceptual level
* replication and failure domains
* Kubernetes CSI integration
* one deliberate OSD-failure experiment

## Full Topic List (if you choose to go deeper anyway)

* pools
* placement groups
* RBD
* CephFS
* resharding / rebalancing internals
* disk-failure recovery at scale

Tag:

```text
v13-storage-optional
```

---

# Failure Engineering

Failure is intentionally introduced throughout the laboratory — manually at first, then increasingly automated through the chaos engineering tooling introduced in Stage 9.

---

## Application Failures

* Go panic
* Python worker crash
* incorrect configuration
* memory exhaustion
* CPU saturation
* slow dependency
* invalid release
* deadlock
* goroutine leak

---

## Network Failures

* DNS failure
* packet loss
* high latency
* connection reset
* blocked port
* TLS failure
* gRPC deadline exceeded
* unreachable service
* network partition

---

## Redis Failures

* cache miss storm
* stale cache
* primary failure
* replica failure
* Sentinel failover
* Redis Cluster node loss
* hot-key overload
* memory exhaustion
* eviction storm
* RDB/AOF recovery
* cache/database inconsistency

---

## PostgreSQL Failures

* primary failure
* replication lag
* full disk
* broken connection pool
* failed backup
* failed restore
* incorrect promotion
* connection exhaustion

---

## RabbitMQ Failures

* broker failure
* queue backlog
* consumer crash
* poison message
* lost acknowledgement
* duplicate delivery

---

## Kafka Failures

* broker failure
* consumer lag
* partition failure
* rebalance
* duplicate processing
* slow consumer

---

## Kubernetes Failures

* Pod crash
* CrashLoopBackOff
* bad readiness probe
* failed node
* broken Service
* incorrect NetworkPolicy
* exhausted resources
* failed rollout

---

## GitOps Failures

* manual drift from desired state (`kubectl edit` on a GitOps-managed resource)
* broken manifest merged to the tracked branch
* failed reconciliation / sync error
* stuck rollout with a failing health check

---

## Ceph Failures *(optional — see Stage 13)*

* OSD failure
* node loss
* degraded placement groups
* disk failure
* recovery
* rebalancing

---

# Incident Documentation

Every meaningful failure should produce:

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
│   │   └── tests/
│   │
│   ├── analytics-service/
│   ├── audit-service/
│   └── notification-service/
│
├── proto/
│   ├── task/
│   └── user/
│
├── infrastructure/
│   ├── virtual-machines/
│   ├── cloud-init/
│   ├── networking/
│   ├── ansible/
│   └── terraform/
│       ├── modules/
│       └── environments/
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
├── gitops/
│   ├── argocd/            # or flux/
│   └── applications/
│
├── database/
│   ├── migrations/
│   ├── replication/
│   ├── backup/
│   └── restore/
│
├── redis/
│   ├── standalone/
│   ├── sentinel/
│   ├── cluster/
│   ├── benchmarks/
│   └── failure-labs/
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
├── chaos/
│   └── experiments/
│
├── storage/
│   └── ceph/               # optional / stretch — see Stage 13
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

# Local Laboratory & Cloud Footprint

Stages 0–2 run entirely on local Linux virtual machines — deliberately cheap and fast to iterate on.

Example local topology:

```text
Linux Host
│
├── VM 1: GitLab / CI Runner / Registry
├── VM 2: Load Balancer
├── VM 3: Kubernetes Control Plane
├── VM 4: Kubernetes Worker
├── VM 5: Kubernetes Worker
└── VM 6: Monitoring / Storage
```

Virtualization:

* KVM/QEMU
* Libvirt
* Proxmox, VirtualBox, or VMware Workstation (any one is fine)

Treat the hypervisor layer as a means to an end, not a skill to master — enough to reliably run the lab, no more. The time investment belongs in what runs on top of it.

From Stage 3 onward, a real cloud provider account (AWS, GCP, or Azure — pick one and stay there) is introduced for anything provisioned through Terraform, including, by Stage 6, the Kubernetes cluster itself. This does not need to run continuously: provision with `terraform apply`, do the stage's work, then `terraform destroy`. A free-tier account is sufficient for a small single-node-pool cluster run for a few hours at a time.

The environment simulates:

* private networks
* cloud instances
* security groups
* block storage
* internal DNS
* load balancers
* private services
* public endpoints
* node failures
* network partitions

Heavy components do not need to run simultaneously.

---

# Security Principles

* no secrets committed to Git
* SSH key authentication
* least-privilege service accounts
* least-privilege IAM roles for CI/CD and Terraform service accounts
* TLS
* non-root containers
* minimal images
* Kubernetes RBAC
* NetworkPolicies
* resource limits
* protected CI/CD variables
* secret rotation
* isolated database credentials
* isolated Redis credentials
* Terraform state stored remotely with locking, never committed to Git
* backup encryption where applicable
* dependency scanning
* container image scanning

---

# Definition of Done

A stage is complete only when it includes:

* working implementation
* repeatable installation
* infrastructure provisioned via Terraform, not created by hand
* automation where applicable
* architecture documentation
* health checks
* metrics
* logs
* traces where relevant
* security considerations
* backup procedures where relevant
* restore testing
* failure simulations
* troubleshooting runbooks
* clear documentation

Installing a technology alone does **not** complete a stage.

---

# Final Production Simulation

The final demonstration exercises the full software and infrastructure lifecycle.

```text
 0. Provision the cloud footprint and Kubernetes cluster with `terraform apply`.

 1. Push a source-code change.

 2. Run linting.

 3. Run Go and Python unit tests.

 4. Run integration tests.

 5. Run gRPC contract tests.

 6. Build container images.

 7. Scan container images.

 8. Push images to the registry.

 9. Update the GitOps manifest with the new image tag and commit.

10. Watch ArgoCD/Flux reconcile staging automatically.

11. Run smoke tests.

12. Promote the release by merging to the production manifest path.

13. Watch ArgoCD/Flux reconcile production.

14. Generate application traffic.

15. Observe Redis cache hits and misses.

16. Observe PostgreSQL queries.

17. Observe RabbitMQ queues.

18. Observe Kafka consumer lag.

19. Trace a request across services.

20. View service metrics.

21. Search centralized logs.

22. Introduce a faulty application release.

23. Detect the incident through alerts.

24. Diagnose the problem.

25. Roll back the deployment via a Git revert, not `kubectl rollout undo`.

26. Confirm recovery.

27. Fail a Redis node.

28. Verify Redis failover.

29. Fail the PostgreSQL primary.

30. Verify database recovery/failover.

31. Fail a Kubernetes worker.

32. Observe workload rescheduling.

33. (Optional) Fail a Ceph OSD.

34. (Optional) Observe distributed-storage recovery.

35. Publish incident postmortems.

36. Run `terraform destroy`, then confirm you can bring everything back with `terraform apply` plus a GitOps sync alone.
```

---

# Distributed Systems Concepts Demonstrated

The laboratory intentionally explores:

* partial failure
* unreliable networks
* timeouts
* cancellation
* retries
* exponential backoff
* jitter
* duplicate requests
* idempotency
* eventual consistency
* strong vs weak consistency
* replication
* leader/follower systems
* failover
* quorum concepts
* sharding
* caching
* cache invalidation
* stale reads
* distributed coordination
* message ordering
* event replay
* consumer coordination
* backpressure
* service discovery
* load balancing
* horizontal scaling
* distributed tracing
* graceful degradation
* fault isolation
* data ownership
* stateless vs stateful workloads

---

# Skills Demonstrated

This project demonstrates practical experience in:

* Software Engineering
* Go
* Python
* Bash
* SQL
* Linux Engineering
* DevOps
* Platform Engineering
* Site Reliability Engineering
* Distributed Systems
* Networking
* TCP/IP
* HTTP/2
* gRPC
* Protocol Buffers
* Infrastructure Automation
* Infrastructure as Code (Terraform / OpenTofu)
* Cloud Platforms (AWS / GCP / Azure)
* CI/CD
* GitOps (ArgoCD / Flux)
* Containers
* Kubernetes
* PostgreSQL
* Redis
* RabbitMQ
* Kafka
* Event-Driven Architecture
* Database High Availability
* Distributed Caching
* Distributed Storage (optional/stretch)
* Prometheus
* Grafana
* OpenTelemetry
* Observability
* Security Engineering
* Failure Engineering
* Chaos Engineering
* Incident Response
* Disaster Recovery

---

# Project Status

```text
[ ] Stage 0  — Application Foundation
[ ] Stage 1  — Native Linux Services
[ ] Stage 2  — Networking and TLS
[ ] Stage 3  — Infrastructure as Code (Ansible + Terraform)
[ ] Stage 4  — Containerization
[ ] Stage 5  — CI/CD
[ ] Stage 6  — Kubernetes
[ ] Stage 7  — GitOps & Progressive Delivery
[ ] Stage 8  — Observability
[ ] Stage 9  — SRE Practices
    ── core path complete above this line ──
[ ] Stage 10 — PostgreSQL High Availability   (stretch)
[ ] Stage 11 — Redis & Distributed Caching    (stretch)
[ ] Stage 12 — Distributed Messaging          (stretch)
[ ] Stage 13 — Distributed Storage            (stretch, optional)
[ ] Final Production Simulation
```

---

# Engineering Philosophy

The project follows this progression:

```text
Understand
    ↓
Implement manually
    ↓
Observe behavior
    ↓
Introduce failure
    ↓
Debug
    ↓
Automate
    ↓
Provision as code
    ↓
Containerize
    ↓
Orchestrate
    ↓
Deliver continuously (GitOps)
    ↓
Observe
    ↓
Scale
    ↓
Fail
    ↓
Recover
    ↓
Improve
```

The emphasis is not on memorizing commands or collecting technologies.

The emphasis is on understanding **how production distributed systems actually behave**.

---

# Disclaimer

This repository is an educational and portfolio laboratory.

It simulates production engineering practices but should not be considered a complete production platform without additional:

* threat modeling
* security review
* compliance controls
* full capacity planning at production scale
* performance/load testing at scale
* continuous, large-scale chaos engineering programs (this lab covers introductory chaos tooling in Stage 9, not a mature program)
* organizational processes
* production-specific operational controls
* cloud cost management at scale

---

# License

MIT License.
