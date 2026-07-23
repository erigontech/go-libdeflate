// Package libdeflate provides Go bindings for the libdeflate C library.
// The C sources are bundled in this module — no system libdeflate.so required.
//
// libdeflate is a heavily optimized library for DEFLATE/zlib/gzip compression
// and decompression. See: https://github.com/ebiggers/libdeflate
package libdeflate

/*
#cgo CFLAGS: -I./libdeflate
#cgo CFLAGS: -I./libdeflate/lib

#include "libdeflate/libdeflate.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// DefaultCompression is the default compression level (6).
const DefaultCompression = 6

// MaxCompression is the maximum compression level (12).
const MaxCompression = 12

// Compressor wraps a libdeflate_compressor.
type Compressor struct {
	c *C.struct_libdeflate_compressor
}

// Decompressor wraps a libdeflate_decompressor.
type Decompressor struct {
	d *C.struct_libdeflate_decompressor
}

// NewCompressor allocates a new compressor with the given compression level (1-12).
// Use DefaultCompression (6) for a balanced trade-off.
func NewCompressor(level int) (*Compressor, error) {
	if level < 1 || level > 12 {
		return nil, fmt.Errorf("libdeflate: compression level must be between 1 and 12, got %d", level)
	}
	c := C.libdeflate_alloc_compressor(C.int(level))
	if c == nil {
		return nil, errors.New("libdeflate: failed to allocate compressor")
	}
	return &Compressor{c: c}, nil
}

// Close frees the compressor.
func (c *Compressor) Close() {
	if c.c != nil {
		C.libdeflate_free_compressor(c.c)
		c.c = nil
	}
}

// NewDecompressor allocates a new decompressor.
func NewDecompressor() (*Decompressor, error) {
	d := C.libdeflate_alloc_decompressor()
	if d == nil {
		return nil, errors.New("libdeflate: failed to allocate decompressor")
	}
	return &Decompressor{d: d}, nil
}

// Close frees the decompressor.
func (d *Decompressor) Close() {
	if d.d != nil {
		C.libdeflate_free_decompressor(d.d)
		d.d = nil
	}
}

// --- Deflate (raw) ---

// CompressDeflate compresses src into dst using raw DEFLATE format.
// dst must be large enough; use DeflateCompressBound to get the upper bound.
// Returns the number of bytes written, or an error if dst is too small.
func (c *Compressor) CompressDeflate(dst, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	n := C.libdeflate_deflate_compress(
		c.c,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
	)
	if n == 0 {
		return 0, errors.New("libdeflate: deflate compress failed: output buffer too small")
	}
	return int(n), nil
}

// DeflateCompressBound returns the maximum compressed size for srcLen bytes.
func (c *Compressor) DeflateCompressBound(srcLen int) int {
	return int(C.libdeflate_deflate_compress_bound(c.c, C.size_t(srcLen)))
}

// DecompressDeflate decompresses a raw DEFLATE stream.
// outLen must be the exact uncompressed size.
// Returns actual bytes written.
func (d *Decompressor) DecompressDeflate(dst []byte, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	var actualOut C.size_t
	result := C.libdeflate_deflate_decompress(
		d.d,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
		&actualOut,
	)
	if result != C.LIBDEFLATE_SUCCESS {
		return 0, fmt.Errorf("libdeflate: deflate decompress failed: %d", int(result))
	}
	return int(actualOut), nil
}

// --- Gzip ---

// CompressGzip compresses src into dst using the gzip format.
func (c *Compressor) CompressGzip(dst, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	n := C.libdeflate_gzip_compress(
		c.c,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
	)
	if n == 0 {
		return 0, errors.New("libdeflate: gzip compress failed: output buffer too small")
	}
	return int(n), nil
}

// GzipCompressBound returns the maximum gzip-compressed size for srcLen bytes.
func (c *Compressor) GzipCompressBound(srcLen int) int {
	return int(C.libdeflate_gzip_compress_bound(c.c, C.size_t(srcLen)))
}

// DecompressGzip decompresses a gzip stream.
// Returns actual bytes written.
func (d *Decompressor) DecompressGzip(dst []byte, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	var actualOut C.size_t
	result := C.libdeflate_gzip_decompress(
		d.d,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
		&actualOut,
	)
	if result != C.LIBDEFLATE_SUCCESS {
		return 0, fmt.Errorf("libdeflate: gzip decompress failed: %d", int(result))
	}
	return int(actualOut), nil
}

// --- Zlib ---

// CompressZlib compresses src into dst using the zlib format.
func (c *Compressor) CompressZlib(dst, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	n := C.libdeflate_zlib_compress(
		c.c,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
	)
	if n == 0 {
		return 0, errors.New("libdeflate: zlib compress failed: output buffer too small")
	}
	return int(n), nil
}

// ZlibCompressBound returns the maximum zlib-compressed size for srcLen bytes.
func (c *Compressor) ZlibCompressBound(srcLen int) int {
	return int(C.libdeflate_zlib_compress_bound(c.c, C.size_t(srcLen)))
}

// DecompressZlib decompresses a zlib stream.
// Returns actual bytes written.
func (d *Decompressor) DecompressZlib(dst []byte, src []byte) (int, error) {
	if len(src) == 0 {
		return 0, nil
	}
	var actualOut C.size_t
	result := C.libdeflate_zlib_decompress(
		d.d,
		unsafe.Pointer(&src[0]), C.size_t(len(src)),
		unsafe.Pointer(&dst[0]), C.size_t(len(dst)),
		&actualOut,
	)
	if result != C.LIBDEFLATE_SUCCESS {
		return 0, fmt.Errorf("libdeflate: zlib decompress failed: %d", int(result))
	}
	return int(actualOut), nil
}

// --- Convenience helpers ---

// GzipCompress compresses src and returns the gzip-compressed bytes.
func GzipCompress(src []byte, level int) ([]byte, error) {
	c, err := getCompressor(level)
	if err != nil {
		return nil, err
	}
	defer putCompressor(level, c)
	dst := make([]byte, c.GzipCompressBound(len(src)))
	n, err := c.CompressGzip(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// GzipDecompress decompresses a gzip stream, assuming the uncompressed size
// is at most maxOutLen bytes.
func GzipDecompress(src []byte, maxOutLen int) ([]byte, error) {
	d, err := getDecompressor()
	if err != nil {
		return nil, err
	}
	defer putDecompressor(d)
	dst := make([]byte, maxOutLen)
	n, err := d.DecompressGzip(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// ZlibCompress compresses src and returns the zlib-compressed bytes.
func ZlibCompress(src []byte, level int) ([]byte, error) {
	c, err := getCompressor(level)
	if err != nil {
		return nil, err
	}
	defer putCompressor(level, c)
	dst := make([]byte, c.ZlibCompressBound(len(src)))
	n, err := c.CompressZlib(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}

// ZlibDecompress decompresses a zlib stream.
func ZlibDecompress(src []byte, maxOutLen int) ([]byte, error) {
	d, err := getDecompressor()
	if err != nil {
		return nil, err
	}
	defer putDecompressor(d)
	dst := make([]byte, maxOutLen)
	n, err := d.DecompressZlib(dst, src)
	if err != nil {
		return nil, err
	}
	return dst[:n], nil
}
