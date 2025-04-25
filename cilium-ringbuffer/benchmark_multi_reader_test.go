package ciliumbuffer

import (
	"context"
	"fmt"
	"sync"
	"testing"

	v1 "github.com/cilium/cilium/pkg/hubble/api/v1"
	"github.com/cilium/cilium/pkg/hubble/container"
)

const (
	readersCount = 5
)

func BenchmarkRingReaderMultipleReaders(b *testing.B) {
	ir := container.NewRing(container.Capacity1023)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create and start multiple reader goroutines
	var wg sync.WaitGroup
	wg.Add(readersCount)

	reader := container.NewRingReader(ir, ir.LastWrite())

	for range readersCount {
		go func() {
			defer wg.Done()
			for {
				event := reader.NextFollow(ctx)
				if event == nil {
					return
				}
			}
		}()
	}

	// Reset the timer before starting the actual benchmark
	b.ResetTimer()

	for i := range b.N {
		ir.Write(&v1.Event{
			Event: fmt.Sprintf("event-%d", i),
		})
	}

	b.StopTimer()
	cancel()
	wg.Wait()
}

func BenchmarkStructChanMultipleReaders(b *testing.B) {
	// Create a channel with a size similar to the ring buffer
	ch := make(chan *v1.Event, 1023)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(readersCount)

	for range readersCount {
		go func() {
			defer wg.Done()
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
	}

	// Reset the timer before starting the actual benchmark
	b.ResetTimer()

	for i := range b.N {
		ch <- &v1.Event{
			Event: fmt.Sprintf("event-%d", i),
		}
	}

	// Wait for all readers to finish
	b.StopTimer()
	cancel()
	wg.Wait()
}
