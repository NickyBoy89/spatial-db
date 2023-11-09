package main

import (
	"math/rand"
	"os"
	"time"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func populateStorageDir(dir string, maxSpread float64, numPoints int) string {
	var server storage.SimpleServer

	server.StorageDir = dir
	defer os.RemoveAll(server.StorageDir)

	points := make([]world.BlockPos, numPoints)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numPoints; i++ {
		points[i] = world.BlockPos{
			X: int(r.NormFloat64() * maxSpread),
			Y: uint(r.NormFloat64() * maxSpread),
			Z: int(r.NormFloat64() * maxSpread),
		}
	}

	for _, point := range points {
		if err := server.ChangeBlock(point, world.Generic); err != nil {
			panic(err)
		}
	}

	return server.StorageDir
}

// func BenchmarkInsertSomePoints(b *testing.B) {
// 	var server storage.SimpleServer
//
// 	stdDev := 65536
//
// 	storage.ChunkFileDirectory = setupStorageDir()
// 	defer os.RemoveAll(storage.ChunkFileDirectory)
//
// 	points := make([]world.BlockPos, b.N)
//
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//
// 	for i := 0; i < b.N; i++ {
// 		points[i] = world.BlockPos{
// 			X: int(r.NormFloat64() * float64(stdDev)),
// 			Y: uint(r.NormFloat64() * float64(stdDev)),
// 			Z: int(r.NormFloat64() * float64(stdDev)),
// 		}
// 	}
//
// 	b.ResetTimer()
//
// 	for _, point := range points {
// 		if err := server.ChangeBlock(point, world.Generic); err != nil {
// 			b.Error(err)
// 		}
// 	}
//
// 	fmt.Println(os.ReadDir(storage.ChunkFileDirectory))
// }
