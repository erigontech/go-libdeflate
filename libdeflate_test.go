package libdeflate

import (
	"bytes"
	"testing"
)

var testData = []byte("Hello, libdeflate! This is a test string for compression. " +
	"libdeflate is a heavily optimized C library for DEFLATE/zlib/gzip.")

func TestGzipRoundtrip(t *testing.T) {
	compressed, err := GzipCompress(testData, DefaultCompression)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	t.Logf("original=%d compressed=%d", len(testData), len(compressed))

	decompressed, err := GzipDecompress(compressed, len(testData)*2)
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	if !bytes.Equal(decompressed, testData) {
		t.Fatal("decompressed data does not match original")
	}
}

func TestZlibRoundtrip(t *testing.T) {
	compressed, err := ZlibCompress(testData, DefaultCompression)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	decompressed, err := ZlibDecompress(compressed, len(testData)*2)
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	if !bytes.Equal(decompressed, testData) {
		t.Fatal("decompressed data does not match original")
	}
}

func TestDeflateRoundtrip(t *testing.T) {
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

	dst := make([]byte, c.DeflateCompressBound(len(testData)))
	n, err := c.CompressDeflate(dst, testData)
	if err != nil {
		t.Fatalf("compress: %v", err)
	}
	compressed := dst[:n]

	out := make([]byte, len(testData))
	n, err = d.DecompressDeflate(out, compressed)
	if err != nil {
		t.Fatalf("decompress: %v", err)
	}
	if !bytes.Equal(out[:n], testData) {
		t.Fatal("decompressed data does not match original")
	}
}

func TestInvalidLevel(t *testing.T) {
	_, err := NewCompressor(0)
	if err == nil {
		t.Fatal("expected error for level 0")
	}
	_, err = NewCompressor(13)
	if err == nil {
		t.Fatal("expected error for level 13")
	}
}

func BenchmarkGzipCompress(b *testing.B) {
	data := bytes.Repeat(testData, 100)
	c, _ := NewCompressor(DefaultCompression)
	defer c.Close()
	dst := make([]byte, c.GzipCompressBound(len(data)))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.CompressGzip(dst, data) //nolint:errcheck
	}
}

func BenchmarkGzipDecompress(b *testing.B) {
	data := bytes.Repeat(testData, 100)
	compressed, _ := GzipCompress(data, DefaultCompression)
	d, _ := NewDecompressor()
	defer d.Close()
	dst := make([]byte, len(data))
	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.DecompressGzip(dst, compressed) //nolint:errcheck
	}
}
