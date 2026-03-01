# hostmonitor

A lightweight host monitoring dashboard built with **SvelteKit** + **Node.js**.

## Architecture

- **Frontend:** SvelteKit (Svelte 5, TailwindCSS)
- **Backend:** SvelteKit server routes (`+server.ts`) — no separate backend process
- **Metrics collection:** [`systeminformation`](https://systeminformation.io/) npm package
- **Runtime:** Node.js (via `@sveltejs/adapter-node`)

The app is a single Node.js process serving both the UI and the `/api/hostinfo` endpoint.

## Endpoints

| Route | Description |
|-------|-------------|
| `GET /` | Dashboard UI |
| `GET /api/hostinfo` | JSON system metrics |

## Metrics collected

- CPU: physical/virtual core count, usage %, load averages, temperature
- Memory: total, free, cached, swap
- Disk: total/free on `/`
- OS: hostname, uptime, platform, kernel version, boot time
- Network interfaces list
- Running processes count
- Container detection

## Development

```bash
npm install
npm run dev        # dev server on http://localhost:5173
```

## Production build

```bash
npm run build
node build         # runs on port 3000
```

Set `PORT` env var to override the default port (3000).

## Docker

```bash
docker build -t hostmonitor .
docker run -p 3000:3000 hostmonitor
```

## Environment variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT`   | `3000`  | HTTP port   |
| `HOST`   | `0.0.0.0` | Bind address |
| `ORIGIN` | — | Required in production behind a reverse proxy |
