terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = "ap-south-1" 
}

variable "server_count" {
  description = "Number of servers dictated by the Go FFD algorithm"
  type        = number
}

resource "aws_instance" "render_node" {
  count         = var.server_count
  ami           = "ami-0287a05f0ef0e9d9a" 
  instance_type = "t2.micro"              
  
  tags = {
    Name = "Smart-Render-Node-${count.index + 1}"
    ManagedBy = "Golang-FFD-Scheduler"
  }
}