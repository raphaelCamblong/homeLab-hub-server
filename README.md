# Homelab Server

![Go](https://img.shields.io/badge/Go-1.20-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## Overview

Homelab Server is a Go-based application designed to manage and serve various functionalities for your home lab. It is built following the principles of Clean Architecture to ensure maintainability, testability, and scalability.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Building and Running](#building-and-running)
- [Directory Structure](#directory-structure)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Features

- **Clean Architecture**: Ensures a scalable and maintainable codebase.
- **Fast and Lightweight**: Built with Go, known for its performance and efficiency.
- **Extensible**: Easily add new features and modules.
- **Secure**: Implements best practices for security in a home lab environment.

## Getting Started

### Prerequisites

Ensure you have the following installed:

- Go 1.20 or later
- Docker (for containerized builds)
- `just` (a handy command runner)

### Installation

Clone the repository:

```sh
git clone https://github.com/yourusername/homelab_server.git
cd homelab_server
```

## API:
```
/api
├── /system
│   ├── /info                 # Get system info (hostname, OS, uptime)
│   ├── /health               # Check if system is healthy
│
├── /network
│   ├── /devices              # List devices on the network (DHCP scan)
│   ├── /interfaces           # List network interfaces and stats
│   ├── /interfaces/{id}      # Get details for a specific network interface
│   ├── /firewall             # View firewall rules (iptables/nftables)
│   ├── /firewall/rules       # Modify firewall rules
│   ├── /dns                  # View and modify DNS settings
│   ├── /services             # View and see available and deployed services
│
├── /nas
│   ├── /storage              # List storage pools and usage
│   ├── /storage/{id}         # Details of a specific storage pool
│   ├── /zfs                  # ZFS-related metrics (if using ZFS)
│   ├── /files                # Browse NAS filesystem
│   ├── /files/upload         # Upload a file to the NAS
│   ├── /files/download       # Download a file
│   ├── /backup               # Trigger a backup job
│   ├── /metrics
│
├── /cluster 
│   ├── /nodes                # List Kubernetes nodes
│   ├── /pods                 # List running pods
│   ├── /pods/{namespace}     # List pods in a specific namespace
│   ├── /services             # List Kubernetes services
│   ├── /deployments          # List Kubernetes deployments
│   ├── /logs/{pod}           # Get logs from a specific pod
│   ├── /events               # Get cluster events
│   ├── /apply                # Apply a YAML manifest
│   ├── /delete               # Delete a resource
│   ├── /metrics
│
├── /users
│   ├── /login                # User authentication
│   ├── /logout               # Log out user
│   ├── /permissions          # View user roles & permissions
│   ├── /audit-logs           # View audit logs (user actions)
│   ├── /register             # Register a new user
```