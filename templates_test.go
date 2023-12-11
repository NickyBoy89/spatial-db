package main

import (
	"errors"
	"io"
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
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
			if errors.Is(err, storage.ChunkNotFoundError) || errors.Is(err, io.EOF) {
				continue
			} else {
				b.Error(err)
			}
		}
	}
}
