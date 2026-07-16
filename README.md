# LocalBox

> Cloud on your own machine. File storage with a web interface, fully under your control.

![Status](https://img.shields.io/badge/status-in%20development-blue)
![Language](https://img.shields.io/badge/language-Go-00ADD8?logo=go)
![Storage](https://img.shields.io/badge/storage-SQLite%20%2B%20Redis-success)

## The Idea

LocalBox is a private cloud storage that runs locally. Upload files, organize them into folders, download whenever you need. No third-party servers, just you and your data on one machine.

Built in Go with a vanilla JS frontend. Fast, simple, no unnecessary bloat.

## Structure

```
cmd/
  main.go             application entry point

deployment/docker/ Docker deployment

internal/app/
  app.go              initialization, route configuration
  endpoint/           HTTP handlers (login, register, files)
  service/            business logic for files and authentication
  mw/                 middleware for session validation
  repository/
    database/         SQLite: user and file storage
    redis/            Redis: session management and caching

website/
  LocalCloudMain.html login/register menu
  login.html          login form
  register.html       signup form
  dashboard.html      cloud interface (all JS magic happens here)
```

## How It Works

1. **Request hits** the Echo server (port 8080)
2. **Middleware checks** cookies (valid session in Redis?)
3. **Endpoint handles** the request (signup, file upload, download)
4. **Service talks** to the database (SQLite stores users and metadata, Redis caches sessions)
5. **Frontend sends** AJAX, updates the UI

## Get Started

```bash
# 1. Clone
git clone https://github.com/rukiamuq-hard/LocalBox.git
cd LocalBox

# 2. start docker
make up

# 3. Open in browser
# http://localhost:8080

```

That's it. Sign up, upload files, use it.

## Stack

- **Go 1.26.3** — language and runtime
- **Echo v5** — web framework
- **SQLite** (modernc.org) — user and file database
- **Redis** (go-redis) — session storage
- **HTML + CSS + Vanilla JS** — interface (no frameworks)
- **Docker** - for deployment, no interface

## Features

| Feature | Status |
|---------|--------|
| Sign up and login | ✓ |
| File upload | ✓ |
| Folder organization | ✓ |
| File download | ✓ |
| Drag-and-drop move | ✓ |
| Multi-user support | ✓ |
