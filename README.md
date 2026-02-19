# Poll Creator

[![CI/CD Pipeline](https://github.com/CoderXAyush/Poll-Creator/actions/workflows/ci.yml/badge.svg)](https://github.com/CoderXAyush/Poll-Creator/actions/workflows/ci.yml)

A full-stack polling application built with **Go** and **React**. Create polls, share them via link, and collect votes with real-time results visualization.

## Tech Stack

| Layer    | Technology                  |
| -------- | --------------------------- |
| Backend  | Go (net/http, no framework) |
| Frontend | React 18, Vite |
| Serve    | Nginx (reverse proxy + SPA) |
| Infra    | Docker, Docker Compose      |

## Features

- Create polls with 2–8 options
- Vote with duplicate prevention (session/IP based)
- Live results with animated bar charts
- Share polls via copyable link
- Close polls to stop voting
- Fully containerized with multi-stage Docker builds

## Getting Started

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)

### Run

```bash
docker compose up --build -d
```

| Service  | URL                     |
| -------- | ----------------------- |
| App      | http://localhost:3001    |
| API      | http://localhost:5001    |

### Stop

```bash
docker compose down
```

## API Endpoints

```
POST   /api/polls            Create a poll
GET    /api/polls            List all polls
GET    /api/polls/:id        Get poll details
POST   /api/polls/:id/vote   Submit a vote
GET    /api/polls/:id/results Get poll results
PATCH  /api/polls/:id/close  Close a poll
GET    /api/health           Health check
```

### Example — Create a Poll

```bash
curl -X POST http://localhost:5001/api/polls \
  -H "Content-Type: application/json" \
  -d '{"question": "Tabs or spaces?", "options": ["Tabs", "Spaces"]}'
```

## Project Structure

```
poll-creator/
├── backend/
│   ├── main.go           # Go API server
│   ├── go.mod
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── App.jsx
│   │   ├── api.js
│   │   ├── components/
│   │   │   ├── CreatePoll.jsx
│   │   │   ├── PollList.jsx
│   │   │   ├── PollView.jsx
│   │   │   └── ResultsChart.jsx
│   │   └── index.css
│   ├── nginx.conf
│   └── Dockerfile
└── docker-compose.yml
```

## Local Development (without Docker)

**Backend:**

```bash
cd backend
go run main.go
```

**Frontend:**

```bash
cd frontend
npm install
npm run dev
```

Vite dev server proxies `/api` requests to `localhost:5001`.

## 🚀 Deployment

### Deploy to Railway

Railway provides the easiest deployment for this Docker-based application:

#### Quick Deploy (Recommended)
```bash
# Run our automated deployment script
./scripts/deploy-railway.sh
```

#### Manual Deploy
1. **Create Railway Account**
   - Go to [railway.app](https://railway.app)
   - Sign up with your GitHub account

2. **Deploy Your Repository**
   ```bash
   # Option 1: Deploy from GitHub (Recommended)
   # Connect your GitHub repo at railway.app/new

   # Option 2: Deploy using Railway CLI
   npm install -g @railway/cli
   railway login
   railway init
   railway up
   ```

3. **Configure Environment Variables**
   - In Railway dashboard, go to your project
   - Add these environment variables:
     - `PORT=8080`
     - `NODE_ENV=production`

4. **Your App Will Be Live!**
   - Railway provides a public URL like `https://poll-creator-production.up.railway.app`
   - The app serves both frontend and API from the same URL

### Alternative Deployment Options

<details>
<summary>🔧 Render</summary>

1. Connect your GitHub repo to [Render](https://render.com)
2. Choose "Web Service"
3. Set build command: `docker build -f Dockerfile.production .`
4. Set start command: `./start.sh`
5. Add environment variable: `PORT=10000`

</details>

<details>
<summary>🌊 DigitalOcean App Platform</summary>

1. Create a new app in [DigitalOcean](https://cloud.digitalocean.com/apps)
2. Connect your GitHub repository
3. Choose "Dockerfile" build method
4. Set Dockerfile path: `Dockerfile.production`
5. Configure environment variables as needed

</details>

<details>
<summary>🐳 Docker Self-Hosting</summary>

```bash
# Build production image
docker build -f Dockerfile.production -t poll-creator .

# Run container
docker run -p 8080:8080 -e PORT=8080 poll-creator
```

</details>

