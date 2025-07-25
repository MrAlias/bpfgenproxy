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
$ docker-compose up
[+] Running 4/4
 ✔ Network bpfgenproxy_my_app_network  Created                                                                                                                           0.1s
 ✔ Container goproxy.opentelemetry.io  Created                                                                                                                           0.1s
 ✔ Container go.opentelemetry.io       Created                                                                                                                           0.1s
 ✔ Container bpfgenproxy-client-1      Created                                                                                                                           0.0s
Attaching to client-1, go.opentelemetry.io, goproxy.opentelemetry.io
client-1                  | go: downloading go.opentelemetry.io/auto v0.22.1
goproxy.opentelemetry.io  | time=2025-07-25T16:47:52.735Z level=DEBUG source=/app/main.go:50 msg=Downloading Fetcher.path=go.opentelemetry.io/auto Fetcher.version=v0.22.1
goproxy.opentelemetry.io  | time=2025-07-25T16:47:52.735Z level=DEBUG source=/app/main.go:52 msg="Serving local files for go.opentelemetry.io/auto v0.22.1"
# [...] 
client-1                  | 2025/07/25 16:49:14 failed to create instrumentation: invalid ID: -1
client-1                  | exit status 1
client-1 exited with code 1
```

The logs of the client and module proxy show the requests being made.
The client will fail to run because the target PID for the instrumentation is invalid.
However, it is important to note that the module was successfully downloaded and it started, meaning the proxy served the module packaged with the eBPF object files correctly.
If the eBPF object files were not present, we would have seen an error like `pattern bpf_x86_bpfel.o: no matching files found` when the client tried to run.

### Clean Up

```shell
docker-compose down
```
