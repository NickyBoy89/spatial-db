package main

import (
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

var server storage.SimpleServer

func init() {
	server.SetStorageRoot("skygrid-save")
}

// insertPointTemplate inserts a configurable variety of points into the server
func insertPointTemplate(testDir string, b *testing.B, pointSpread int) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(float64(pointSpread))
		if err := server.ChangeBlock(pos, world.Generic); err != nil {
			b.Error(err)
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

func BenchmarkInsertClusteredPoints(b *testing.B) {
	insertPointTemplate("imperial-test", b, 128)
}

func BenchmarkInsertSparserPoints(b *testing.B) {
	insertPointTemplate("imperial-test", b, 2048)
}
