package libdeflate

import (
	"bytes"
	"sync"
	"testing"
)

func TestPoolCapacityOverride(t *testing.T) {
	defer func() { PoolCapacity = 0 }()
	PoolCapacity = 3
	if got := poolCapacity(); got != 3 {
		t.Fatalf("poolCapacity() = %d, want 3", got)
	}
}

// The pooled convenience helpers must stay correct under concurrent use across
// multiple compression levels: each goroutine must get a compressor/decompressor
// it alone owns for the duration of the call, and round-trips must match.
func TestConvenienceHelpersConcurrent(t *testing.T) {
	data := bytes.Repeat(testData, 4)
	levels := []int{1, DefaultCompression, MaxCompression}

	var wg sync.WaitGroup
	for g := 0; g < 32; g++ {
		wg.Add(1)
		go func(g int) {
			defer wg.Done()
			level := levels[g%len(levels)]
			for i := 0; i < 200; i++ {
				comp, err := GzipCompress(data, level)
				if err != nil {
					t.Errorf("compress: %v", err)
					return
				}
				out, err := GzipDecompress(comp, len(data))
				if err != nil {
					t.Errorf("decompress: %v", err)
					return
				}
				if !bytes.Equal(out, data) {
					t.Error("round-trip mismatch")
					return
				}
			}
		}(g)
	}
	wg.Wait()
}

func BenchmarkGzipCompressOneShotSerial(b *testing.B) {
	data := bytes.Repeat(testData, 8)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := GzipCompress(data, DefaultCompression); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGzipCompressOneShotParallel(b *testing.B) {
	data := bytes.Repeat(testData, 8)
	b.SetBytes(int64(len(data)))
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if _, err := GzipCompress(data, DefaultCompression); err != nil {
				b.Fatal(err)
			}
		}
	})
}
