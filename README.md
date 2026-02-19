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

### Deploy to Vercel (Recommended)

Vercel provides excellent support for React applications with serverless backend functions:

#### Quick Deploy
```bash
# Run our automated deployment script
./scripts/deploy-vercel.sh
```

#### Manual Deploy
1. **Create Vercel Account**
   - Go to [vercel.com](https://vercel.com)
   - Sign up with your GitHub account

2. **Deploy Your Repository**
   ```bash
   # Option 1: Import from GitHub (Recommended)
   # Go to vercel.com/new and import your GitHub repository
   # Vercel will automatically detect the configuration

   # Option 2: Deploy using Vercel CLI
   npm install -g vercel
   vercel login
   vercel --confirm
   vercel --prod
   ```

3. **Automatic Configuration**
   - Vercel automatically detects the `vercel.json` configuration
   - Frontend is built and deployed from `/frontend`  
   - Backend API functions are deployed from `/api`
   - No additional environment variables needed!

4. **Your App Will Be Live!**
   - Vercel provides a public URL like `https://poll-creator-username.vercel.app`
   - Frontend and API are served from the same domain
   - Automatic HTTPS and global CDN included

### Alternative Deployment Options

### Alternative Deployment Options

<details>
<summary>🚅 Railway (Docker-based)</summary>

1. Connect your GitHub repo to [Railway](https://railway.app)
2. Railway will use the `Dockerfile.production` for deployment
3. Set environment variables: `PORT=8080`, `NODE_ENV=production`
4. Single container serves both frontend and backend

</details>

<details>
<summary>🔧 Render</summary>

1. Connect your GitHub repo to [Render](https://render.com)  
2. Create a "Web Service"
3. Choose "Static Site" for frontend deployment
4. Choose "Web Service" with Docker for full-stack deployment
5. Use `Dockerfile.production` for build

</details>

<details> 
<summary>🌊 Netlify</summary>

1. Connect your GitHub repo to [Netlify](https://netlify.com)
2. Set build directory: `frontend/dist`
3. Set build command: `cd frontend && npm run build`
4. For API, use Netlify Functions (requires code adaptation)

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

