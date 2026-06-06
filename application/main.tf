terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 3.0.0"
    }
  }
}

provider "docker" {
}

variable "server_count" {
  description = "Number of worker containers to run"
  type        = number
  default     = 0
}

resource "docker_container" "worker_node" {
  count = var.server_count 
  name  = "worker-node-${count.index + 1}"
  image = "video-worker:v1"
}