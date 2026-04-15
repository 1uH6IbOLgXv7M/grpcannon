# grpcannon

> A lightweight gRPC load testing CLI with configurable concurrency profiles and latency histograms

---

## Installation

```bash
go install github.com/yourorg/grpcannon@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/grpcannon.git && cd grpcannon && go build -o grpcannon .
```

---

## Usage

```bash
grpcannon [flags] <target>
```

### Example

```bash
grpcannon \
  --proto ./api/service.proto \
  --call helloworld.Greeter/SayHello \
  --data '{"name": "world"}' \
  --concurrency 50 \
  --requests 10000 \
  localhost:50051
```

### Key Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--proto` | | Path to the `.proto` file |
| `--call` | | Fully qualified method name |
| `--data` | | JSON request payload |
| `--concurrency` | `10` | Number of concurrent workers |
| `--requests` | `1000` | Total number of requests |
| `--duration` | | Run for a fixed duration (e.g. `30s`) |
| `--histogram` | `true` | Print latency histogram on completion |

### Sample Output

```
Summary:
  Total requests : 10000
  Duration       : 4.82s
  Throughput     : 2074 req/s

Latency (ms):
  p50  :  23.1
  p90  :  41.7
  p99  :  88.4
  p999 : 134.2

Histogram:
  [  0-10ms] ████░░░░░░  823
  [ 10-25ms] ████████░░ 4201
  [ 25-50ms] ██████░░░░ 3612
  [ 50-100ms] ██░░░░░░░  987
  [100ms+  ] ░░░░░░░░░░  377
```

---

## License

MIT © yourorg