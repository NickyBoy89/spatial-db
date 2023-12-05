package main

import (
	"errors"
	"io"
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

var server storage.InMemoryServer

func init() {
	server.SetStorageRoot("skygrid-save")
}

func readBlockTemplate(rootDir string, b *testing.B, pointSpread int) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(float64(pointSpread))
		if _, err := server.ReadBlockAt(pos); err != nil {
			if errors.Is(err, storage.ChunkNotFoundError) || errors.Is(err, io.EOF) {
				continue
			} else {
				b.Error(err)
			}
		}
	}
}

func fetchChunkTemplate(testDir string, b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(2048).ToChunkPos()
		if _, err := server.ReadChunkAt(pos); err != nil {
			b.Error(err)
		}
	}
}

// Insert blocks

func BenchmarkReadClusteredPoints(b *testing.B) {
	readBlockTemplate("skygrid-test", b, 128)
}

func BenchmarkReadSparserPoints(b *testing.B) {
	readBlockTemplate("skygrid-test", b, 2048)
}

func BenchmarkReadSparserPoints1(b *testing.B) {
	readBlockTemplate("skygrid-test", b, 65536)
}
