package main

import (
	"errors"
	"testing"

	"github.com/NickyBoy89/spatial-db/storage"
	"github.com/NickyBoy89/spatial-db/world"
)

func readBlockTemplate(
	storageServer storage.StorageServer,
	b *testing.B,
	pointSpread int,
) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(float64(pointSpread))
		if _, err := storageServer.ReadBlockAt(pos); err != nil {
			if errors.Is(err, storage.ChunkNotFoundError) {
				continue
			} else {
				b.Error(err)
			}
		}
	}
}
