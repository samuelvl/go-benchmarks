# go-benchmarks

The benchmarks are run on a MacBook Pro with an Apple M3 Pro chip:

```shell
goos: darwin
goarch: arm64
cpu: Apple M3 Pro
```

## cilium-ringbuffer

### Single reader

- The ring buffer is significantly faster than channels (1.89x faster).
- The ring buffer uses more memory (108 bytes vs 64 bytes).
- The ring buffer makes one more allocation per operation.

```shell
BenchmarkRingReader
BenchmarkRingReader-12    	11001187	       107.1 ns/op	     108 B/op	       5 allocs/op
BenchmarkStructChan
BenchmarkStructChan-12    	 6310552	       202.5 ns/op	      64 B/op	       4 allocs/op
```

### Multiple readers

- The ring buffer is significantly faster than channels (2.48x faster).
- The ring buffer uses slightly more memory (85 bytes vs 64 bytes).
- Same number of allocations.

```shell
BenchmarkRingReaderMultipleReaders
BenchmarkRingReaderMultipleReaders-12    	10699506	       104.4 ns/op	      85 B/op	       4 allocs/op
BenchmarkStructChanMultipleReaders
BenchmarkStructChanMultipleReaders-12    	 4622535	       258.6 ns/op	      64 B/op	       4 allocs/op
```
