# Pulse: Mini Checkly

A lightweight synthetic monitoring system that executes HTTP checks on a schedule, logs performance, and triggers alerts on failure.

## Features

- **Scheduled HTTP Checks**: Execute HTTP requests on configurable intervals
- **Background Job Processing**: Redis-based queue system for scalable check execution
- **Time-Series Analytics**: ClickHouse integration for high-performance metrics storage and queries
- **Webhook Alerts**: Configurable alert thresholds with webhook notifications
- **Health Monitoring**: Built-in health, readiness, and metrics endpoints
- **Horizontal Scaling**: Worker-based architecture supports multiple worker instances
- **Performance Tracking**: Detailed latency and status tracking for all checks
- **Self-Monitoring**: System exposes metrics to monitor itself

## Architecture

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

## Prerequisites

- **Go 1.24+**
- **PostgreSQL 15+**
- **Redis 7+** (or Valkey)
- **ClickHouse** (optional, for analytics)
- **Docker and Docker Compose** (for infrastructure services)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/mukund/mini-checkly.git
cd mini-checkly
```

### 2. Start Infrastructure

#### Option A: Docker Compose (Recommended)

Start everything with a single command:

```bash
docker compose up -d
```

This starts:

- PostgreSQL on port `5432`
- Valkey (Redis) on port `6379`
- ClickHouse on ports `8123` (HTTP) and `9000` (Native protocol)
- API server on port `8080`
- Worker process

To view logs:

```bash
docker compose logs -f
```

To stop:

```bash
docker compose down
```

#### Option B: Local Development

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

   > **Note**: `CLICKHOUSE_DSN` is optional. If not provided or if ClickHouse is unavailable, the service will continue to work without analytics.

3. Start server (runs migrations automatically):

   ```bash
   make run
   ```

4. Start worker in a separate terminal:

   ```bash
   make run-worker
   ```

### 3. Test the Setup

Create your first check:

```bash
curl -X POST http://localhost:8080/checks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google Homepage",
    "url": "https://www.google.com",
    "method": "GET",
    "expected_status": 200,
    "interval_seconds": 60,
    "alert_threshold": 3,
    "timeout_ms": 10000,
    "webhook_url": "https://your-webhook-url.com/alerts"
  }'
```

## API Overview

### Check Management

#### Create Check

```bash
curl -X POST http://localhost:8080/checks \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Google Homepage",
    "url": "https://www.google.com",
    "method": "GET",
    "expected_status": 200,
    "interval_seconds": 60,
    "alert_threshold": 3,
    "timeout_ms": 10000,
    "webhook_url": "https://your-webhook-url.com/alerts"
  }'
```

#### List All Checks

```bash
curl -X GET http://localhost:8080/checks
```

#### Get Check Details

```bash
curl -X GET http://localhost:8080/checks/{check-id}
```

#### Get Check Runs

```bash
curl -X GET http://localhost:8080/checks/{check-id}/runs?limit=100
```

Returns check execution history from ClickHouse. Optional `limit` query parameter (default: 100, max: 1000).

#### Get Check Alerts

```bash
curl -X GET http://localhost:8080/checks/{check-id}/alerts?limit=100
```

Returns all alerts for a check. Optional `limit` query parameter (default: 100, max: 1000).

#### Get Check Webhook Attempts

```bash
curl -X GET http://localhost:8080/checks/{check-id}/webhooks?limit=100
```

Returns all webhook delivery attempts for a check, including retries. Optional `limit` query parameter (default: 100, max: 1000).

Each webhook attempt includes:
- Request details (URL, body, headers)
- Response details (status code, body, headers)
- Error messages (if any)
- Latency and timeout information
- Retry number

### Health & Observability

#### Health Check

```bash
curl -X GET http://localhost:8080/health
```

Returns overall health status with dependency checks:

- **Status**: `ok`, `degraded`, or `down` based on dependency health
- **Uptime**: Process uptime in seconds
- **Dependencies**: Health status of PostgreSQL, Redis, and ClickHouse
- **Version**: Application version
- **Build Time**: Process start time

#### Readiness Probe

```bash
curl -X GET http://localhost:8080/ready
```

Returns `200 OK` only when all critical dependencies (PostgreSQL, Redis) are reachable. Returns `503 Service Unavailable` if any dependency is unreachable. Used by orchestration systems to determine when the service is ready to accept traffic.

#### System Metrics

```bash
curl -X GET http://localhost:8080/metrics
```

Returns system-wide operational metrics:

- **checks_executed_total**: Total number of checks executed since startup
- **checks_failed_total**: Total number of failed checks
- **alerts_sent_total**: Total number of alerts sent
- **avg_latency_ms**: Average latency across all check runs (from ClickHouse)
- **worker.active_jobs**: Current number of active jobs being processed
- **worker.queue_depth**: Number of jobs waiting in the Redis queue
- **worker.uptime**: Worker process uptime in seconds

> **Note**: Both the API server (port 8080) and worker process (port 8081) expose these endpoints for self-monitoring.

## Execution Flow

1. **Scheduler** polls PostgreSQL every 10 seconds for due checks (`next_run_at <= now()`)
2. **Enqueue** check IDs to Redis queue and update `next_run_at` in PostgreSQL
3. **Workers** dequeue jobs from Redis and execute HTTP checks
4. **Results** saved to ClickHouse for analytics, state updated in PostgreSQL
5. **Alerter** sends webhooks when failure thresholds are reached

## Design Decisions

- **Redis queue**: Enables horizontal scaling and decouples scheduler from workers
- **ClickHouse for time-series**: Optimized for high-volume writes and analytics queries
- **PostgreSQL for relational data**: Stores checks and alerts metadata
- **Interval-based scheduling**: Simpler than cron, easier to reason about
- **Webhook alerts**: Flexible integration with any external system
- **Self-monitoring**: System exposes metrics to monitor itself, ensuring operational failures cause visible failure modes

## Analytics

The service automatically records check execution metrics to ClickHouse when available:

- **Check Results**: Response times, status codes, success/failure status
- **Performance Metrics**: Latency tracking, timeout detection
- **Alert History**: Webhook delivery status and timing

Analytics are **optional** - if ClickHouse is not configured or unavailable, the service continues to function normally without analytics.

### Analytics Tables

- `check_results`: Individual check execution results with timestamps
- `check_metrics`: Aggregated performance metrics

## Known Limitations

- No job retry mechanism or timeouts
- Single scheduler instance (no distributed locking)
- No authentication or rate limiting on API
- Basic error handling for webhooks

## Scaling Considerations

At 10M checks/day (~115 checks/second), potential bottlenecks:

- **PostgreSQL**: Connection pooling and read replicas
- **Redis**: Cluster mode, shard queues by check ID
- **ClickHouse**: Batch inserts, cluster scaling
- **Workers**: Limit concurrency, use worker pools, rate limit requests

Scaling strategies: horizontal worker scaling, database sharding, Redis caching, async alert processing.

## Configuration

All configuration is done via environment variables or a `.env` file:

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://user:password@localhost:5432/pulse?sslmode=disable` | PostgreSQL connection string |
| `REDIS_URL` | `localhost:6379` | Redis connection string |
| `CLICKHOUSE_DSN` | `clickhouse://default@localhost:9000/default` | ClickHouse connection string (optional) |
| `PORT` | `8080` | API server port |

> **Note**: `CLICKHOUSE_DSN` is optional. If not provided or if ClickHouse is unavailable, the service will continue to work without analytics.

## Development

### Build

```bash
# Build server
make build

# Build worker
make build-worker
```

### Run Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

## Project Structure

```text
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
│   ├── handlers/   # HTTP request handlers
│   ├── metrics/    # Metrics tracking with atomic counters
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

## Security Considerations

- **No Authentication**: API endpoints are currently unauthenticated (add authentication for production use)
- **Webhook URLs**: Validate webhook URLs to prevent SSRF attacks
- **Rate Limiting**: Consider implementing rate limiting to prevent abuse
- **Input Validation**: All user inputs should be validated and sanitized

## Contributing

Contributions are welcome! Feel free to open issues or submit pull requests. Whether it's bug fixes, new features, documentation improvements, or examples, your contributions help make this project better.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). See the [LICENSE](LICENSE) file for details.
