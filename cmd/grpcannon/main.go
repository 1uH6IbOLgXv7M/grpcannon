// Command grpcannon is a lightweight gRPC load testing CLI with configurable
// concurrency profiles and latency histograms.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/example/grpcannon/internal/config"
	"github.com/example/grpcannon/internal/output"
	"github.com/example/grpcannon/internal/runner"
	"github.com/example/grpcannon/internal/signal"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("grpcannon", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		target      = fs.String("target", "", "gRPC server address (host:port)")
		call        = fs.String("call", "", "fully-qualified method name (package.Service/Method)")
		data        = fs.String("data", "{}", "request payload as JSON")
		concurrency = fs.Int("concurrency", 10, "number of concurrent workers")
		rps         = fs.Int("rps", 0, "requests per second limit (0 = unlimited)")
		duration    = fs.Duration("duration", 10*time.Second, "test duration")
		timeout     = fs.Duration("timeout", 5*time.Second, "per-request timeout")
		warmup      = fs.Int("warmup", 0, "number of warmup requests before recording")
		insecure    = fs.Bool("insecure", false, "skip TLS verification")
		format      = fs.String("format", "text", "output format: text or json")
		protoFile   = fs.String("proto", "", "path to .proto file for reflection fallback")
	)

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return nil
		}
		return err
	}

	cfg := config.Default()
	cfg.Target = *target
	cfg.Call = *call
	cfg.Data = *data
	cfg.Concurrency = *concurrency
	cfg.RPS = *rps
	cfg.Duration = *duration
	cfg.Timeout = *timeout
	cfg.WarmupRequests = *warmup
	cfg.Insecure = *insecure
	cfg.ProtoFile = *protoFile

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background())
	defer stop()

	out, err := output.New(*format, os.Stdout)
	if err != nil {
		return fmt.Errorf("output: %w", err)
	}

	result, err := runner.Run(ctx, cfg)
	if err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return out.Write(result)
}
