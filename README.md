# ForgeFlow 🚀

ForgeFlow is a modern, lightweight, and language-agnostic Continuous Integration & Continuous Deployment (CI/CD) system. It features a scalable microservice architecture, auto-language detection, and a beautiful React-based dashboard.

## Features ✨

- **Auto-Language Detection**: Simply link a GitHub repository, and ForgeFlow's engine will automatically detect the language (Node.js, Go, Python, Rust, etc.) and configure the build/test commands.
- **Beautiful Dashboard**: A modern, Vercel-like dashboard to monitor pipeline runs, metrics, and organization activity in real-time.
- **Scalable Architecture**: Decoupled API and Worker services using PostgreSQL as a centralized queue and state store.
- **Docker-Native Execution**: Pipelines are executed within sandboxed environments.
- **Live Metrics**: Monitor language distribution, pipeline success rates, and activity over time.

## Tech Stack 🛠️

- **Backend (API & Worker)**: Go
- **Frontend**: React (Vite), TailwindCSS, Recharts, Lucide Icons
- **Database**: PostgreSQL (GORM)
- **Deployment**: Docker & Docker Compose

## Getting Started 🚀

### Prerequisites
- Docker and Docker Compose
- Git

### How to Run

1. **Clone the repository** (if you haven't already):
   ```bash
   git clone https://github.com/SharmilaSugumar/cicd.git
   cd cicd
   ```

2. **Start the services using Docker Compose**:
   This single command will spin up the PostgreSQL database, the Go API backend, the Go worker engine, and the React frontend.
   ```bash
   docker-compose up -d --build
   ```

3. **Access the Dashboard**:
   Open your browser and navigate to:
   ```
   http://localhost
   ```

### Troubleshooting

- **Database Issues**: If you ever need to reset your database, you can wipe the PostgreSQL volume by running:
  ```bash
  docker-compose down -v
  docker-compose up -d --build
  ```
- **Checking Logs**: To see what's happening under the hood (for example, in the worker):
  ```bash
  docker-compose logs -f worker
  ```

## Architecture 🏗️

- **API Service** (`/cmd/api`): Handles REST endpoints for organizations, projects, pipelines, jobs, and metrics.
- **Worker Service** (`/cmd/worker`): Polls the database for pending jobs, clones repositories into a local workspace, executes auto-detected build/test commands, captures standard output/error, and updates job statuses.
- **Web Frontend** (`/web`): Communicates with the API to render pipeline statuses and statistical graphs.

## License 📄
MIT License
