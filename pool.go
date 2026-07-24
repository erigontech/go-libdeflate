package libdeflate

import (
	"fmt"
	"runtime"
	"sync"
)

// PoolCapacity is the max idle contexts the convenience helpers keep per pool.
// 0 (default) means auto: 2*GOMAXPROCS clamped to [4, 32]. A positive value is
// used verbatim. Set it before the first compress/decompress call.
var PoolCapacity int

var (
	compressorPools  [MaxCompression + 1]chan *Compressor
	compressorOnce   [MaxCompression + 1]sync.Once
	decompressorCh   chan *Decompressor
	decompressorOnce sync.Once
)

func poolCapacity() int {
	if PoolCapacity > 0 {
		return PoolCapacity
	}
	switch n := 2 * runtime.GOMAXPROCS(0); {
	case n < 4:
		return 4
	case n > 32:
		return 32
	default:
		return n
	}
}

func getCompressor(level int) (*Compressor, error) {
	if level < 1 || level > MaxCompression {
		return nil, fmt.Errorf("libdeflate: compression level must be between 1 and %d, got %d", MaxCompression, level)
	}
	compressorOnce[level].Do(func() {
		compressorPools[level] = make(chan *Compressor, poolCapacity())
	})
	select {
	case c := <-compressorPools[level]:
		return c, nil
	default:
		return NewCompressor(level)
	}
}

func putCompressor(level int, c *Compressor) {
	select {
	case compressorPools[level] <- c:
	default:
		c.Close()
	}
}

func getDecompressor() (*Decompressor, error) {
	decompressorOnce.Do(func() {
		decompressorCh = make(chan *Decompressor, poolCapacity())
	})
	select {
	case d := <-decompressorCh:
		return d, nil
	default:
		return NewDecompressor()
	}
}

func putDecompressor(d *Decompressor) {
	select {
	case decompressorCh <- d:
	default:
		d.Close()
	}
}
