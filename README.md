# ☁️ Bin-Packing Cloud Resource Scheduler

A custom, microservice-based auto-scaling orchestration engine built from scratch in Go. 

This project acts as a "Local Cloud Orchestrator." It dynamically scales worker containers up and down based on real-time workload calculations using the **First-Fit Decreasing (FFD) bin-packing algorithm**. It bridges Go's concurrency with Terraform's Infrastructure-as-Code (IaC) and RabbitMQ's message brokering to create a resilient, production-like scaling daemon.

## 🏗️ How It Works (The Architecture)

1. **Ingestion (REST API):** A Go-based HTTP server receives simulated high-fidelity rendering jobs (represented by integer weights).
2. **Message Broker (RabbitMQ):** To prevent data loss and decouple the ingestion from processing, jobs are safely published to a durable RabbitMQ queue (`render_jobs`).
3. **The Master Daemon:** A Go background worker polls RabbitMQ every 30 seconds to fetch batches of queued jobs.
4. **Algorithmic Scheduling:** The First-Fit Decreasing (FFD) algorithm sorts the job weights and calculates the absolute mathematical minimum number of servers (bins) required to handle the batch without exceeding capacity.
5. **Infrastructure Provisioning:** The daemon dynamically triggers **Terraform** to provision the exact number of custom Alpine-based Docker worker nodes required to process the load.

## 🛠️ Tech Stack

* **Core Logic & API:** Go (Golang)
* **Infrastructure as Code (IaC):** Terraform
* **Containerization:** Docker & Multi-stage Docker Builds
* **Message Broker:** RabbitMQ
* **Algorithm:** First-Fit Decreasing (FFD)

## 📂 Project Structure

```text
├── application/
│   ├── cmd/autoscaler/         # Main daemon entry point (Orchestrator)
│   ├── internal/
│   │   ├── api/                # REST API server & RabbitMQ publisher
│   │   └── worker/             # Custom payload running inside Docker containers
│   ├── pkg/
│   │   ├── provisioner/        # Terraform execution and bridge logic
│   │   └── scheduler/          # FFD bin-packing algorithm implementation
│   ├── .gitignore
│   ├── Dockerfile              # Multi-stage build for the custom worker node
│   ├── go.mod                  # Go module dependencies
│   └── main.tf                 # Terraform configuration for Docker provider
└── ffd.cpp                     # Initial C++ implementation of the FFD logic

```

## 🚀 Getting Started

### Prerequisites

* [Go](https://golang.org/dl/) (v1.22+)
* [Docker Desktop](https://www.docker.com/products/docker-desktop)
* [Terraform](https://www.terraform.io/downloads.html)

### Installation & Setup

**1. Start the RabbitMQ Message Broker**
Spin up a local RabbitMQ instance using Docker:

```bash
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

```

**2. Build the Custom Worker Image**
Build the Docker image that Terraform will use for the worker nodes:

```bash
docker build -t video-worker:v1 .

```

**3. Initialize Terraform**

```bash
terraform init -upgrade

```

**4. Start the Orchestration Daemon**

```bash
go run cmd/autoscaler/main.go

```

### Triggering the Auto-Scaler

Once the daemon is running, open a new terminal and send jobs to the API. (Examples below use PowerShell):

```powershell
Invoke-RestMethod -Uri "http://localhost:8080/job" -Method POST -ContentType "application/json" -Body '{"weight": 8}'
Invoke-RestMethod -Uri "http://localhost:8080/job" -Method POST -ContentType "application/json" -Body '{"weight": 4}'
Invoke-RestMethod -Uri "http://localhost:8080/job" -Method POST -ContentType "application/json" -Body '{"weight": 6}'

```

Watch the master daemon output. On its next 30-second tick, it will pull these jobs, calculate the required infrastructure, and you will see new worker containers spin up in Docker Desktop!

## 🛑 Graceful Shutdown & Cleanup

The daemon is built to handle OS interrupt signals (e.g., `CTRL+C`). When terminated, it intercepts the kill signal, halts new scheduling, and safely triggers `terraform destroy -auto-approve` to tear down all active worker containers before exiting, preventing orphaned cloud resources.

```

```