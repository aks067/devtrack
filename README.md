# devtrack

> A lightweight developer productivity tracker written in Go — automatically monitors your coding activity and displays your stats in the terminal.

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![Platform](https://img.shields.io/badge/platform-Linux-blue)
![License](https://img.shields.io/badge/license-MIT-green)

---

## What it does

devtrack runs silently in the background and tracks how much time you spend on each project by watching for file changes using **inotify** (zero CPU when idle). It also fetches your GitHub commit stats and displays everything in a clean terminal dashboard.

```
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  📁 devtrack          ⏱ 2h34m   📦 12 commits
  ████████████████████░░░░░░░░░░

  📁 NeuralNetworkInC  ⏱ 1h15m   📦 9 commits
  ████████████░░░░░░░░░░░░░░░░░░

  📁 Snake-AI          ⏱ 0h42m   📦 7 commits
  ███████░░░░░░░░░░░░░░░░░░░░░░░
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
```

---

## Features

- **Automatic time tracking** — detects which project you're working on based on file changes, no manual start/stop
- **Recursive directory watching** — monitors all subdirectories using inotify, zero polling overhead
- **Session persistence** — sessions are saved to a local JSON file and accumulated over time
- **GitHub integration** — fetches your recent commits per repository via the GitHub API
- **Clean terminal dashboard** — projects sorted by time spent, with ASCII progress bars
- **Graceful shutdown** — Ctrl+C saves the current session before exiting

---

## Installation

```bash
git clone https://github.com/aks067/devtrack
cd devtrack
go build -o devtrack .
```

---

## Usage

**Start the tracker** (watches your projects folder in the background):
```bash
./devtrack daemon
```

**View your stats:**
```bash
./devtrack status
```

**Stop the tracker:**
```
Ctrl+C  — saves the current session and exits cleanly
```

---

## Configuration

In `cmd/commands.go`, set your projects directory and GitHub username:

```go
tracker.Start("/home/yourname/projects")
```

```go
githubTrack.FetchCommits("your_github_username")
```

> A config file (~/.config/devtrack/config.json) is planned for a future release.

---

## How it works

```
devtrack/
├── main.go                      # CLI entry point
├── cmd/
│   └── commands.go         # Starts the file watcher and Renders the dashboard
├── internal/
│   ├── tracker/tracker.go       # inotify watcher, session logic, JSON persistence
│   └── github/github.go         # GitHub REST API client
└── data/
    └── logs.json                # Local session storage
```

1. The daemon uses [fsnotify](https://github.com/fsnotify/fsnotify) to watch your projects directory recursively
2. When a file is modified, it detects the active project (first-level subfolder)
3. After 10 minutes of inactivity, the session is automatically closed and saved
4. `devtrack status` reads the local logs and combines them with live GitHub data

---

## Tech stack

- **Go** — core language
- **fsnotify** — cross-platform filesystem notifications (inotify on Linux)
- **GitHub REST API** — commit stats, no authentication required for public repos
- **encoding/json** — session persistence
- **os/signal** — graceful shutdown on SIGINT/SIGTERM

---

## Upcoming Features

- [ ] Config File
- [ ] Web dashboard (local, self-hosted)
- [ ] GitHub streak tracking
- [ ] Support for private repositories (token-based auth)

---

## Licence 

OpenSource
