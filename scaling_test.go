package main

import (
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

// insertPointTemplate inserts a configurable variety of points into the server
func insertPointTemplate(testDir string, b *testing.B, pointSpread int) {
	var server storage.InMemoryServer

	server.SetStorageRoot(testDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(float64(pointSpread))
		if err := server.ChangeBlock(pos, world.Generic); err != nil {
			b.Error(err)
		}
	}
}

func fetchChunkTemplate(testDir string, b *testing.B) {
	var server storage.SimpleServer

	server.SetStorageRoot(testDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(2048).ToChunkPos()
		if _, err := server.ReadChunkAt(pos); err != nil {
			b.Error(err)
		}
	}
}

// Insert blocks

func BenchmarkInsertClusteredPoints(b *testing.B) {
	insertPointTemplate("test-world", b, 128)
}

func BenchmarkInsertSparserPoints(b *testing.B) {
	insertPointTemplate("test-world", b, 2048)
}
