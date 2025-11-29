# Pulse: Mini Checkly

A lightweight synthetic monitoring system that executes HTTP checks on a schedule, logs performance, and triggers alerts on failure.

## Overview

Pulse demonstrates background job scheduling, HTTP reliability engineering, persistence layer design with PostgreSQL and ClickHouse, and alerting workflows.

## System Architecture

```mermaid
flowchart TD
    Scheduler[Scheduler] -->|Polls every 10s| PostgreSQL[(PostgreSQL)]
    Scheduler -->|Enqueue check IDs| Redis[(Redis Queue)]
    Redis -->|Dequeue jobs| Workers[Workers]
    Workers -->|Load config| PostgreSQL
    Workers -->|Save results| ClickHouse[(ClickHouse)]
    Workers -->|Update state| PostgreSQL
    Workers -->|Trigger alerts| Alerter[Alerter]
    Alerter -->|Send webhooks| Webhook[Webhook Endpoint]
```

## Execution Flow

1. Scheduler polls PostgreSQL for due checks (`next_run_at <= now()`)
2. Enqueues check IDs to Redis and updates `next_run_at`
3. Workers dequeue jobs and execute HTTP checks
4. Results saved to ClickHouse, state updated in PostgreSQL
5. Alerter sends webhooks when failure thresholds are reached

## Design Decisions

- **Redis queue**: Enables horizontal scaling and decouples scheduler from workers
- **ClickHouse for time-series**: Optimized for high-volume writes and analytics queries
- **PostgreSQL for relational data**: Stores checks and alerts metadata
- **Interval-based scheduling**: Simpler than cron, easier to reason about
- **Webhook alerts**: Flexible integration with any external system

## Known Limitations

- No job retry mechanism or timeouts
- Single scheduler instance (no distributed locking)
- No authentication or rate limiting on API
- Basic error handling for webhooks
- No metrics or observability (Prometheus, tracing)

## Scaling Considerations

At 10M checks/day (~115 checks/second), potential bottlenecks:

- **PostgreSQL**: Connection pooling and read replicas
- **Redis**: Cluster mode, shard queues by check ID
- **ClickHouse**: Batch inserts, cluster scaling
- **Workers**: Limit concurrency, use worker pools, rate limit requests

Scaling strategies: horizontal worker scaling, database sharding, Redis caching, async alert processing.

## Quick Start

### Option 1: Docker Compose (Recommended)

Start everything with a single command:

```bash
docker compose up -d
```

This starts PostgreSQL, Valkey (Redis), ClickHouse, the API server, and worker. The server will be available at `http://localhost:8080`.

To view logs:

```bash
docker compose logs -f
```

To stop:

```bash
docker compose down
```

### Option 2: Local Development

**Prerequisites:** Go 1.24+, Docker and Docker Compose

1. Start infrastructure services:

   ```bash
   docker compose -f docker-compose.infrastructure.yml up -d
   ```

2. Create `.env` file:

   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/pulse?sslmode=disable
   REDIS_URL=localhost:6379
   CLICKHOUSE_DSN=clickhouse://default@localhost:9000/default
   PORT=8080
   ```

3. Start server (runs migrations automatically):

   ```bash
   make run
   ```

4. Start worker in a separate terminal:

   ```bash
   make run-worker
   ```

## API Endpoints

- `GET /checks` - List all checks
- `GET /checks/:id` - Get check details
- `POST /checks` - Create new check

**Example request:**

```json
{
  "name": "Google Homepage",
  "url": "https://www.google.com",
  "method": "GET",
  "expected_status": 200,
  "interval_seconds": 60,
  "alert_threshold": 3,
  "timeout_ms": 10000,
  "webhook_url": "https://your-webhook-url.com/alerts"
}
```

## Development

```bash
make build        # Build server
make build-worker # Build worker
make run          # Run server
make run-worker   # Run worker
make test         # Run tests
```

## Project Structure

```
pulse/
├── cmd/
│   ├── server/     # HTTP API server
│   └── worker/     # Background worker process
├── internal/
│   ├── alerter/    # Alert processing and webhook sending
│   ├── checker/    # HTTP check execution
│   ├── clickhouse/ # ClickHouse client
│   ├── config/     # Configuration management
│   ├── db/         # PostgreSQL connection and migrations
│   ├── models/     # Data models
│   ├── redis/      # Redis client for job queue
│   ├── scheduler/  # Check scheduling logic
│   ├── store/      # Data access layer
│   └── worker/     # Worker process logic
├── Dockerfile.server
├── Dockerfile.worker
├── docker-compose.yml
└── docker-compose.infrastructure.yml
```
