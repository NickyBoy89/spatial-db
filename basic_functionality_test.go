package main

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func setupStorageDir() string {
	dir, err := os.MkdirTemp("", "spatial-db-persistence")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Temporary directory is at %s\n", dir)

	storage.ChunkFileDirectory = dir

	return dir
}

func BenchmarkInsertSomePoints(b *testing.B) {
	var server storage.SimpleServer

	stdDev := 65536

	storage.ChunkFileDirectory = setupStorageDir()
	defer os.RemoveAll(storage.ChunkFileDirectory)

	points := make([]world.BlockPos, b.N)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < b.N; i++ {
		points[i] = world.BlockPos{
			X: int(r.NormFloat64() * float64(stdDev)),
			Y: uint(r.NormFloat64() * float64(stdDev)),
			Z: int(r.NormFloat64() * float64(stdDev)),
		}
	}

	b.ResetTimer()

	for _, point := range points {
		if err := server.ChangeBlock(point, world.Generic); err != nil {
			b.Error(err)
		}
	}

	fmt.Println(os.ReadDir(storage.ChunkFileDirectory))
}
