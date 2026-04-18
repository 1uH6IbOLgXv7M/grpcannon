package ticker_test

import (
	"context"
	"testing"
	"time"

	"github.com/example/grpcannon/internal/ticker"
)

func BenchmarkTicker_1ms(b *testing.B) {
	tk := ticker.New(time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go tk.Run(ctx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		select {
		case <-tk.C():
		case <-time.After(time.Second):
			b.Fatal("tick timeout")
		}
	}
}
