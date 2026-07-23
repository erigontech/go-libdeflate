package libdeflate

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"
	"testing"
)

// Compressing empty input must yield a valid, decodable stream — not zero
// bytes. Previously the compress methods short-circuited empty input to
// (0, nil), so GzipCompress(nil) returned an empty slice that a standard gzip
// decoder rejects as truncated. We decode with the stdlib here specifically to
// prove the output is a conformant stream, not just something this library can
// read back.
func TestEmptyGzipIsValidStandardStream(t *testing.T) {
	compressed, err := GzipCompress(nil, DefaultCompression)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("empty input compressed to zero bytes; not a valid gzip stream")
	}
	r, err := gzip.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("stdlib gzip.NewReader rejected our output: %v", err)
	}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stdlib gzip read: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("decoded %d bytes, want 0", len(got))
	}
}

func TestEmptyZlibIsValidStandardStream(t *testing.T) {
	compressed, err := ZlibCompress(nil, DefaultCompression)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	if len(compressed) == 0 {
		t.Fatal("empty input compressed to zero bytes; not a valid zlib stream")
	}
	r, err := zlib.NewReader(bytes.NewReader(compressed))
	if err != nil {
		t.Fatalf("stdlib zlib.NewReader rejected our output: %v", err)
	}
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stdlib zlib read: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("decoded %d bytes, want 0", len(got))
	}
}

func TestEmptyDeflateIsValidStandardStream(t *testing.T) {
	c, err := NewCompressor(DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	dst := make([]byte, c.DeflateCompressBound(0))
	n, err := c.CompressDeflate(dst, nil)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	if n == 0 {
		t.Fatal("empty input compressed to zero bytes; not a valid DEFLATE stream")
	}
	r := flate.NewReader(bytes.NewReader(dst[:n]))
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("stdlib flate read rejected our output: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("decoded %d bytes, want 0", len(got))
	}
}

// A zero-length destination buffer must return an error, never panic with a
// slice out-of-range from &dst[0]. This covers both the low-level methods and
// the convenience decompress helper (called with maxOutLen == 0).
func TestZeroLengthDstReturnsErrorNotPanic(t *testing.T) {
	c, err := NewCompressor(DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	d, err := NewDecompressor()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	if _, err := c.CompressGzip(nil, testData); err == nil {
		t.Error("CompressGzip with empty dst: expected error, got nil")
	}
	if _, err := c.CompressZlib(nil, testData); err == nil {
		t.Error("CompressZlib with empty dst: expected error, got nil")
	}
	if _, err := c.CompressDeflate(nil, testData); err == nil {
		t.Error("CompressDeflate with empty dst: expected error, got nil")
	}

	compressed, err := GzipCompress(testData, DefaultCompression)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := d.DecompressGzip(nil, compressed); err == nil {
		t.Error("DecompressGzip with empty dst: expected error, got nil")
	}
	if _, err := GzipDecompress(compressed, 0); err == nil {
		t.Error("GzipDecompress with maxOutLen 0: expected error, got nil")
	}
}
