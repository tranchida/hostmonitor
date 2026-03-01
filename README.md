# HostMonitor

A real-time system monitoring dashboard built with **Go** (backend) and **SvelteKit** (frontend).

## Architecture

```
┌─────────────────────────────────┐
│         Go (Echo) :8080         │
│                                 │
│  GET /api/hostinfo → JSON API   │
│  GET /*            → SvelteKit  │
│                    static files │
│  GET /legacy       → htmx page │
│  GET /host         → HTML frag  │
└─────────────────────────────────┘
         ▲
         │ fetch /api/hostinfo (every 10s)
         │
┌─────────────────────────────────┐
│     SvelteKit Frontend (SPA)    │
│                                 │
│  - TailwindCSS v4 styling       │
│  - Reactive polling via Svelte  │
│  - Static adapter (embedded)    │
└─────────────────────────────────┘
```

### Backend (Go)

- **labstack/echo** as the web framework
- **gopsutil** for system information collection
- **go-humanize** for human-readable units
- Serves the SvelteKit build as embedded static files
- JSON API at `GET /api/hostinfo`
- Legacy htmx routes kept at `/legacy` and `/host`

### Frontend (SvelteKit)

- **SvelteKit** with static adapter (SPA mode)
- **TailwindCSS v4** for styling
- Client-side polling every 10 seconds
- Dark theme (gray-900) with sky-600 accent
- Responsive grid layout (1/2/4 columns)

### Data Exposed

The `/api/hostinfo` endpoint returns a JSON object with:

| Field | Type | Description |
|-------|------|-------------|
| `Hostname` | string | Machine hostname |
| `OS` | string | Operating system |
| `Platform` / `PlatformVersion` | string | Platform details |
| `KernelVersion` | string | Kernel version |
| `BootTime` | string | Boot time (RFC3339) |
| `Uptime` | string | Human-readable uptime |
| `CPUP` / `CPUV` | int | Physical / virtual CPU count |
| `CPUUsage` | string | Current CPU usage % |
| `CPUTemperature` | string | CPU temperature |
| `RunningProcesses` | int | Number of running processes |
| `LoadAverage1/5/15` | string | Load averages |
| `TotalMemory` / `FreeMemory` / `CacheMemory` | string | Memory stats |
| `TotalDiskSpace` / `FreeDiskSpace` | string | Disk stats |
| `TotalSwap` / `FreeSwap` | string | Swap stats |
| `NetworkInterfaces` | string[] | Network interface names |
| `IsRunningInContainer` | bool | Container detection |
| `CurrentTime` | string | Server time (RFC3339) |

## Development

### Prerequisites

- Go 1.24+
- Node.js 22+

### Run the backend

```bash
go run .
# Listens on :8080
```

### Run the frontend (dev mode with hot reload)

```bash
cd frontend
npm install
npm run dev
# Listens on :5173, proxies /api to :8080
```

### Build for production

```bash
# 1. Build the frontend
cd frontend && npm ci && npm run build && cd ..

# 2. Build the Go binary (embeds frontend/build/)
go build -o hostmonitor .
```

### Docker

```bash
docker build -t hostmonitor .
docker run -p 8080:8080 hostmonitor
```

## Legacy

The original htmx-based UI is still available at `/legacy` for backward compatibility.
