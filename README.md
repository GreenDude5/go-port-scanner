# go-port-scanner

A concurrent TCP port scanner written in Go with PostgreSQL persistence and an HTTP API.

## Features

- Scan TCP ports on a target host concurrently
- Save results to PostgreSQL
- HTTP API to query scan results
- Configurable port range and concurrency

## Usage

Start PostgreSQL (requires Docker):

```
docker compose up -d
```

Run a scan:

```
go run ./cmd/scanner --host scanme.nmap.org --start 1 --end 1024 --threads 100
```

Start the HTTP API server:

```
go run ./cmd/scanner --mode server
```

Then query results at `http://localhost:8080/results`.

## Flags

| Flag      | Default            | Description                |
|-----------|--------------------|----------------------------|
| `--mode`  | `scan`             | Run mode: `scan` or `server` |
| `--host`  | `scanme.nmap.org`  | Target hostname            |
| `--start` | `1`                | Start port                 |
| `--end`   | `1024`             | End port                   |
| `--threads` | `100`            | Number of concurrent workers |
