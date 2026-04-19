# go-libdeflate

Go bindings for [ebiggers/libdeflate](https://github.com/ebiggers/libdeflate) v1.25.

**No system `libdeflate.so` required.** The C sources are bundled in this module and compiled via CGo at `go build` time — same pattern used by [`erigontech/secp256k1`](https://github.com/erigontech/secp256k1) and [`erigontech/mdbx-go`](https://github.com/erigontech/mdbx-go).

## Requirements

- CGo enabled (default)
- GCC or Clang

## Install

```bash
go get github.com/erigontech/go-libdeflate
```

## Usage

```go
import libdeflate "github.com/erigontech/go-libdeflate"

// Convenience API
compressed, err := libdeflate.GzipCompress(data, libdeflate.DefaultCompression)
decompressed, err := libdeflate.GzipDecompress(compressed, maxOutLen)

compressed, err = libdeflate.ZlibCompress(data, libdeflate.DefaultCompression)
decompressed, err = libdeflate.ZlibDecompress(compressed, maxOutLen)

// Low-level API — reuse Compressor/Decompressor via sync.Pool for best performance
c, _ := libdeflate.NewCompressor(libdeflate.DefaultCompression)
defer c.Close()
dst := make([]byte, c.GzipCompressBound(len(src)))
n, err := c.CompressGzip(dst, src)
compressed = dst[:n]

d, _ := libdeflate.NewDecompressor()
defer d.Close()
out := make([]byte, originalLen)
n, err = d.DecompressGzip(out, compressed)
```

## Compression levels

| Level | Description |
|-------|-------------|
| 1 | Fastest |
| 6 | Default (balanced) |
| 12 | Best compression |

## Supported formats

- **Gzip** — `CompressGzip` / `DecompressGzip`
- **Zlib** — `CompressZlib` / `DecompressZlib`
- **Deflate (raw)** — `CompressDeflate` / `DecompressDeflate`

## Benchmarks

### Low-level API (library)

```
BenchmarkGzipCompress-16     GB/s  (libdeflate level 6, ~3-5x faster than compress/gzip)
BenchmarkGzipDecompress-16   GB/s
```

Run with: `go test -bench=. -benchmem`

### One-shot compression vs stdlib (real-world JSON payloads, ~30 KB)

| Implementation        | Latency    | Throughput | Allocs/req |
|-----------------------|------------|------------|------------|
| libdeflate (one-shot) | 4.6 ms/req | 122 MB/s   | 63 KB      |
| stdlib compress/gzip  | 8.0 ms/req |  70 MB/s   | 66 KB      |
| **Speedup**           | **~1.75x** | **~1.74x** | —          |

Under high concurrency the advantage grows further: libdeflate uses less CPU per request, leaving more headroom for parallel workloads.

## Repository structure

```
go-libdeflate/
├── libdeflate.go          # CGo bindings + Go API
├── libdeflate_test.go     # Tests and benchmarks
├── libdeflate_*.c         # CGo bridge files (one per translation unit)
└── libdeflate/            # C sources from ebiggers/libdeflate v1.25
    ├── libdeflate.h
    ├── common_defs.h
    └── lib/
        ├── deflate_compress.c / deflate_decompress.c
        ├── gzip_compress.c / gzip_decompress.c
        ├── zlib_compress.c / zlib_decompress.c
        ├── adler32.c, crc32.c, utils.c
        └── x86/, arm/, riscv/   (arch-specific SIMD optimizations)
```

## License

The Go bindings are MIT licensed. The bundled libdeflate C library is MIT licensed by Eric Biggers — see [`libdeflate/LICENSE`](libdeflate/LICENSE) or the [upstream repository](https://github.com/ebiggers/libdeflate).
