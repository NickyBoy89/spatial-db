package main

import (
	"errors"
	"io/fs"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"git.nicholasnovak.io/nnovak/spatial-db/server"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func populateStorageDir(
	dirName string,
	maxSpread float64,
	numPoints int,
	cleanup bool,
) {
	log.Debug("Generating new storage directory")

	// Make sure that another directory is not already at that location
	if _, err := os.Stat(dirName); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Debugf("Making new directory at %s", dirName)
			if err := os.Mkdir(dirName, 0755); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		log.Debug("Directory already exists, skipping generation")
		return
	}

	var server server.SimpleServer

	server.StorageDir = dirName
	if cleanup {
		defer os.RemoveAll(server.StorageDir)
	}

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
	log.Info("Done generating")
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
