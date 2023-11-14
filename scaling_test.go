package main

import (
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

var (
	emptyDir       = "./empty"
	smallSparseDir = "./small-sparse"
	medSparseDir   = "./med-sparse"
	largeSparseDir = "./lg-sparse"
	smallDenseDir  = "./small-dense"
	medDenseDir    = "./med-dense"
	largeDenseDir  = "./lg-dense"
)

// Point densities
const (
	sparse = 1_000
	dense  = 10_000
)

var pointCounts = []int{200, 1_000, 2_500}

var dirs = []string{emptyDir, smallSparseDir, medSparseDir, largeSparseDir,
	smallDenseDir, medDenseDir, largeDenseDir}

func init() {
	for index, dir := range dirs {
		if index > 2 {
			populateStorageDir(dir, sparse, pointCounts[index%3], false)
		} else {
			populateStorageDir(dir, dense, pointCounts[index%3], false)
		}
	}
}

// insertPointTemplate inserts a configurable variety of points into the server
func insertPointTemplate(testDir string, b *testing.B) {
	var server storage.SimpleServer

	server.SetStorageRoot(testDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(2048)
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

func BenchmarkInsertSmallSparse(b *testing.B) {
	insertPointTemplate(smallSparseDir, b)
}

func BenchmarkInsertMedSparse(b *testing.B) {
	insertPointTemplate(medSparseDir, b)
}

func BenchmarkInsertLgSparse(b *testing.B) {
	insertPointTemplate(largeSparseDir, b)
}

func BenchmarkInsertSmallDense(b *testing.B) {
	insertPointTemplate(smallDenseDir, b)
}

func BenchmarkInsertMedDense(b *testing.B) {
	insertPointTemplate(medDenseDir, b)
}

func BenchmarkInsertLgDense(b *testing.B) {
	insertPointTemplate(largeDenseDir, b)
}

// Fetch chunks

func BenchmarkFetchChunkSmallSparse(b *testing.B) {
	fetchChunkTemplate(smallSparseDir, b)
}

func BenchmarkFetchChunkMedSparse(b *testing.B) {
	fetchChunkTemplate(medSparseDir, b)
}

func BenchmarkFetchChunkLgSparse(b *testing.B) {
	fetchChunkTemplate(largeSparseDir, b)
}

func BenchmarkFetchChunkSmallDense(b *testing.B) {
	fetchChunkTemplate(smallDenseDir, b)
}

func BenchmarkFetchChunkMedDense(b *testing.B) {
	fetchChunkTemplate(medDenseDir, b)
}

func BenchmarkFetchChunkLgDense(b *testing.B) {
	fetchChunkTemplate(largeDenseDir, b)
}
