# WordGames

A word games application built with React (frontend) and Go/Gin (backend).

## Quick Start

### Option 1: Using npm script (Recommended)
```bash
# Install dependencies for both frontend and backend
npm install

# Start both servers simultaneously
npm run dev
```

### Option 2: Using Windows batch script
```bash
# Double-click or run the batch file
start-dev.bat
```

### Option 3: Manual start
```bash
# Terminal 1 - Start backend
cd backend
go run main.go

# Terminal 2 - Start frontend
cd frontend
npm run dev
```

## Available Scripts

- `npm run dev` - Start both frontend and backend servers
- `npm run dev:frontend` - Start only the frontend server
- `npm run dev:backend` - Start only the backend server
- `npm run install` - Install dependencies for both frontend and backend

## Development URLs

- **Frontend**: http://localhost:5173
- **Backend API**: http://localhost:8080

## Project Structure

```
wordgames/
├── frontend/          # React application
├── backend/           # Go/Gin server
├── package.json       # Root package.json with dev scripts
├── start-dev.bat      # Windows batch script
└── README.md          # This file
```

## Features

- Header component with customizable title
- Server communication button that fetches data from Go backend
- CORS enabled for local development
- Environment-based configuration 