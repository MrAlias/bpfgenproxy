# bpfgenproxy

Proof-of-concept running a Go module proxy to serve eBPF object files dynamically.

## Architecture Diagram

```diagram
+--------+                                              +-------------------+
|        |-- 1. go mod tidy --------------------------->|                   |
|        |                                              |                   |
|        |                                              | Go Proxy (Vanity) |
|        |                                              |                   |
|        |<- 2. "http://goproxy.opentelemetry.io/auto" -|                   |
|        |                                              +-------------------+
| Client |
|        |                                              +-------------------+
|        |-- 3. Client requests module from proxy ----->|                   |
|        |                                              |                   |
|        |                                              | Module Proxy      |
|        |                                              |                   |
|        |<- 4. Zip file of go.opentelemetry.io/auto ---|                   |
+--------+                                              +-------------------+
```

1. The Go client runs `go mod tidy` and requests the module information for `go.opentelemetry.io/auto` from the Go proxy server serving the vanity URL.
2. The Go proxy server responds, pointing the client to the module proxy and tells the client to use the `mod` protocol for the download (instead of a VCS).
3. The client requests the module from from the module proxy.
4. The module proxy serves a cached copy of the module, which includes pre-generated eBPF object files.

All other modules requested from the module proxy are treated as pass-through requests to the [upstream Go proxy](https://proxy.golang.org).

## Requirements

- Docker
- docker-compose

## Usage

### Running

```shell
docker-compose up -d
```

### Inspect

The logs of the client and module proxy show the requests being made.

```shell
docker-compose logs

```

### Clean Up

```shell
docker-compose down
```
