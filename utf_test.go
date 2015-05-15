package utfutil

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"testing"
)

func BenchmarkEncode(b *testing.B) {
	encodeChunk(1024*128, b)
}

func BenchmarkDecode(b *testing.B) {
	decodeChunk(1024*128, b)
}

func encodeChunk(size int, b *testing.B) {
	p := make([]byte, size)
	io.ReadFull(rand.Reader, p)

	for i := 0; i < b.N; i++ {
		EncodeSlice(p, binary.LittleEndian)
	}

	b.SetBytes(int64(len(p)))
}

func decodeChunk(size int, b *testing.B) {
	p := make([]byte, size)
	io.ReadFull(rand.Reader, p)

	p = EncodeSlice(p, binary.LittleEndian)

	for i := 0; i < b.N; i++ {
		DecodeSlice(p)
	}

	b.SetBytes(int64(len(p)))
}
