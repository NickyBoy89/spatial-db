package main

import (
	"math/rand"
	"testing"
	"time"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func BenchmarkInsertSomePoints(b *testing.B) {
	var server storage.SimpleServer

	points := make([]world.BlockPos, b.N)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		points[i] = world.BlockPos{
			X: int(r.NormFloat64()),
			Y: uint(r.NormFloat64()),
			Z: int(r.NormFloat64()),
		}
	}

	b.ResetTimer()

	for _, point := range points {
		if err := server.ChangeBlock(point, world.Generic); err != nil {
			b.Error(err)
		}
	}
}
