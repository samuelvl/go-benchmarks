package ciliumbuffer

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/cilium/cilium/pkg/hubble/api/v1"
	"github.com/cilium/cilium/pkg/hubble/container"
)

func BenchmarkRingReader(b *testing.B) {
	ir := container.NewRing(container.Capacity1023)

	reader := container.NewRingReader(ir, ir.LastWrite())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			event := reader.NextFollow(ctx)
			if event == nil {
				return
			}
		}
	}()

	// Reset the timer before starting the actual benchmark
	b.ResetTimer()

	for i := range b.N {
		ir.Write(&v1.Event{
			Event: fmt.Sprintf("event-%d", i),
		})
	}

	cancel()
	<-done
}

func BenchmarkStructChan(b *testing.B) {
	// Create a channel with a size similar to the ring buffer
	ch := make(chan *v1.Event, 1023)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			select {
			case event := <-ch:
				if event == nil {
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	// Reset the timer before starting the actual benchmark
	b.ResetTimer()

	for i := range b.N {
		ch <- &v1.Event{
			Event: fmt.Sprintf("event-%d", i),
		}
	}

	cancel()
	<-done
}
