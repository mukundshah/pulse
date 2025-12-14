# Pulse

> **Note:** This is an experimental learning project. Inspired by [Checkly](https://www.checklyhq.com/), it was built to understand how synthetic monitoring systems work. Active maintenance isn't guaranteed, though I may continue development if there's interest.

**Self-hosted synthetic monitoring that actually makes sense.**

Pulse is a self-hosted monitoring solution that gives you full control over your monitoring infrastructure. Monitor your APIs and services with complete data ownership, flexible deployment options, and no external dependencies.

Built with Go for developers who care about performance, simplicity, and control.

## Why Pulse?

**You own your data.** Complete control over where your data lives and how it's stored.

**It's fast.** Built in Go, runs on minimal resources, scales horizontally without breaking a sweat.

**It's simple.** One Docker command to deploy. Clean API. Intuitive interface. Straightforward setup.

**It's complete.** HTTP checks, performance tracking, time-series analytics, flexible alerting—everything you need, nothing you don't.

## ✨ What You Get

- **Lightning Fast** — Go-powered engine that uses minimal resources while handling thousands of checks
- **Time-Series Analytics** — ClickHouse integration for blazing-fast metrics queries and insights
- **Horizontal Scaling** — Worker-based architecture that grows with your infrastructure
- **Flexible Alerting** — Webhook integrations that fit into your existing notification stack
- **Docker-Ready** — Deploy in minutes with a single command, no configuration hell

## 🚀 Get Started in 60 Seconds

### Prerequisites

- **Docker & Docker Compose** (for production deployment)
- **Taskfile** (for running development tasks)
- **Go 1.25+** (only needed for local development)
- **PostgreSQL 18+**, **Redis 8+** (or Valkey), **ClickHouse** (optional, for advanced analytics)

### Production Deployment

```bash
git clone https://github.com/mukundshah/pulse.git
cd pulse
task docker:up
```

### Local Development

```bash
# Start infrastructure services (PostgreSQL, Valkey, ClickHouse)
task docker:infra:up

# Copy environment template
cp .env.example .env

# Run database migrations
task db:migrate

# Start api, worker and web
task dev:all
```

## About This Project

Pulse started as an experimental project to learn how synthetic monitoring systems work. Inspired by [Checkly](https://www.checklyhq.com/), it's a self-hosted implementation built from the ground up to understand the architecture, challenges, and design decisions that go into building a monitoring platform.

This is primarily a learning project, and active maintenance isn't guaranteed. That said, if there's interest and usage, I may continue developing it.

**Contributions, feedback, and discussions are still welcome!**

## 📄 License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**. See the [LICENSE](LICENSE) file for details.

---

**Built for developers who want control over their monitoring infrastructure.**
