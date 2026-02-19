# Contributing to Poll Creator

Thank you for your interest in contributing to Poll Creator! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/) & [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/)
- [Go](https://golang.org/dl/) 1.21+ (for backend development)
- [Node.js](https://nodejs.org/) 20+ (for frontend development)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/Poll-Creator.git
   cd Poll-Creator
   ```

2. **Install Dependencies**
   ```bash
   # Backend
   cd backend && go mod download
   
   # Frontend
   cd ../frontend && npm install
   ```

3. **Start Development Environment**
   ```bash
   # Docker (recommended)
   docker compose up --build
   
   # OR run locally
   # Terminal 1: Backend
   cd backend && go run main.go
   
   # Terminal 2: Frontend
   cd frontend && npm run dev
   ```

## 📋 Development Guidelines

### Code Style

**Go Backend:**
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Handle errors appropriately

**React Frontend:**
- Use functional components with hooks
- Follow ESLint rules
- Use semantic HTML elements
- Keep components small and focused

### Git Workflow

1. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b bugfix/issue-number
   ```

2. **Make Changes**
   - Keep commits atomic and focused
   - Write descriptive commit messages
   - Test your changes thoroughly

3. **Commit Format**
   ```
   type(scope): description
   
   feat(backend): add poll closing functionality
   fix(frontend): resolve vote submission bug
   docs(readme): update installation instructions
   ```

4. **Push and Create PR**
   ```bash
   git push origin your-branch-name
   ```

### Testing

- **Backend**: Write unit tests for new functions
- **Frontend**: Test user interactions and API calls
- **Integration**: Ensure Docker setup works correctly

### Pull Request Process

1. Update documentation if needed
2. Add tests for new functionality
3. Ensure all CI checks pass
4. Fill out the PR template completely
5. Request review from maintainers

## 🐛 Bug Reports

Use the bug report template and include:
- Steps to reproduce
- Expected vs actual behavior
- Environment details
- Screenshots/logs if applicable

## 💡 Feature Requests

Use the feature request template and include:
- Clear description of the feature
- Use case and benefits
- Implementation suggestions (if any)

## 📁 Project Structure

```
poll-creator/
├── backend/           # Go API server
│   ├── main.go       # Main application file
│   └── go.mod        # Go dependencies
├── frontend/         # React application
│   ├── src/
│   │   ├── components/  # Reusable components
│   │   ├── App.jsx      # Main app component
│   │   └── api.js       # API client
│   └── package.json
├── .github/          # GitHub templates and workflows
└── docker-compose.yml
```

## 🔧 Available Scripts

**Backend:**
```bash
go run main.go          # Start development server
go test ./...          # Run tests
go build -o poll-api   # Build binary
```

**Frontend:**
```bash
npm run dev            # Start development server
npm run build          # Build for production
npm test               # Run tests
npm run lint           # Check code style
```

**Docker:**
```bash
docker compose up --build  # Build and start all services
docker compose down        # Stop all services
```

## 🤝 Community

- Be respectful and welcoming
- Follow the GitHub Community Guidelines
- Help others when you can
- Ask questions if you're stuck

## 📄 License

By contributing, you agree that your contributions will be licensed under the same license as this project.

---

Thank you for contributing to Poll Creator! 🎉