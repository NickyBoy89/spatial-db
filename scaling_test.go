package main

import (
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func BenchmarkInsertSparsePoints(b *testing.B) {
	var server storage.SimpleServer

	tempDir := "./data"

	server.StorageDir = populateStorageDir(tempDir, 2048, 1_000)
	b.ResetTimer()

	b.Log("Finished generating directory")

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(2048)
		if err := server.ChangeBlock(pos, world.Generic); err != nil {
			b.Error(err)
		}
	}
}
