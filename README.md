# Poll Creator

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
| App      | http://localhost:3000    |
| API      | http://localhost:5000    |

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
curl -X POST http://localhost:5000/api/polls \
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

Vite dev server proxies `/api` requests to `localhost:5000`.


