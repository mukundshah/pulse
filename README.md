# Pulse

> A modern, self-hosted synthetic monitoring platform built for developers who want control, simplicity, and performance.

**Pulse** is a lightweight monitoring system that continuously checks your APIs and services, tracks performance metrics, and alerts you when things go wrong. Built with Go, designed for scale.

## Features

- 🚀 **Fast & Lightweight**: Built in Go for minimal resource usage
- 📊 **Time-Series Analytics**: Powered by ClickHouse for fast metrics queries
- 🔄 **Horizontally Scalable**: Worker-based architecture
- 🔌 **Webhook Integration**: Flexible alerting
- 📈 **Self-Monitoring**: Built-in metrics and health endpoints
- 🐳 **Docker-Ready**: One-command deployment

## Quick Start

### Prerequisites

- **Go 1.25+** (for local development)
- **Docker & Docker Compose** (for containerized deployment)
- **PostgreSQL 18+**, **Redis 8+** (or Valkey), **ClickHouse** (optional)

### Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/mukund/pulse.git
cd pulse

# Start all services
docker compose -f docker-compose.local.yml up -d

# View logs
docker compose -f docker-compose.local.yml logs -f
```

This starts PostgreSQL, Valkey/Redis, ClickHouse, API server (port 8080), and worker process.

### Local Development

```bash
# Start infrastructure services
docker compose -f docker-compose.local.yml up -d

# Create .env file
cp .env.example .env

# Start server (runs migrations automatically)
make run

# In another terminal, start worker
make run-worker
```

### Project Structure

```text
pulse/
├── cmd/
│   ├── db/              # DB management commands
│   ├── server/          # HTTP API server entrypoint
│   └── worker/          # Background worker entrypoint
├── internal/
│   ├── alerter/         # Alert processing and webhook delivery
│   ├── auth/            # Authentication and authorization
│   ├── checker/         # HTTP check execution engine
│   ├── clickhouse/      # ClickHouse client and queries
│   ├── config/          # Configuration management
│   ├── db/              # PostgreSQL connection and migrations
│   ├── email/           # Email service integration
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware (auth, logging)
│   ├── metrics/         # Metrics tracking
│   ├── models/          # Data models and migrations
│   ├── redis/           # Redis client for job queue
│   ├── scheduler/       # Check scheduling logic
│   ├── store/           # Data access layer
│   └── worker/          # Worker process logic
├── web/                 # Nuxt.js frontend application
├── specs/               # OpenAPI specifications
├── templates/           # Email templates
├── Dockerfile.*         # Container definitions
└── docker-compose*.yml  # Docker Compose configurations
```

## Contributing

Contributions are welcome! Whether it's bug fixes, new features, documentation improvements, or architectural enhancements.

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See the [LICENSE](LICENSE) file for details.

---

**Built with ❤️ for developers who want control over their monitoring infrastructure.**
