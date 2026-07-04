# ForgeFlow Release Checklist

## Completed Features
- [x] Database Schema & GORM Models (UUIDs, Soft Deletes)
- [x] Repository Layer with `SELECT FOR UPDATE SKIP LOCKED`
- [x] Service Layer with strict RBAC and Job State Machine
- [x] Gin REST API (Auth, Projects, Pipelines, Jobs, Queues)
- [x] JWT Authentication & Token Bucket Rate Limiting
- [x] Centralized Configuration (`pkg/config`)
- [x] Prometheus Metrics Endpoint
- [x] Gorilla WebSocket Hub for real-time updates
- [x] Scheduler Daemon for Dependency & Queue management
- [x] Execution Engine & Worker Daemon
- [x] React + Vite + Tailwind Enterprise Frontend
- [x] Docker Compose hardened with MinIO, Redis, RabbitMQ
- [x] GitHub Actions CI Pipeline

## Remaining TODOs
- [ ] Connect MinIO client to the Worker Execution Engine to upload build logs.
- [ ] Migrate the internal in-memory `EventBus` to Redis Pub/Sub so WebSocket events scale across multiple API nodes.
- [ ] Integrate OpenTelemetry spans into the `middleware.Logging()` wrapper.

## Known Limitations
- The API's in-memory rate limiter does not share state across horizontally scaled nodes.
- Local `pkg/eventbus` will only broadcast to WebSockets connected to the exact same API container that triggered the event.

## Future Roadmap
- Kubernetes Helm Chart manifests deployment.
- True Docker-in-Docker (dind) execution sandboxing within the Worker Engine.
- Audit Log persistence to DB on every organization change.
