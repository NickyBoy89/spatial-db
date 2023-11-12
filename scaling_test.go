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

func insertPointTemplate(testDir string, b *testing.B) {
	var server storage.SimpleServer

	server.StorageDir = testDir

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pos := world.RandomBlockPosWithRange(2048)
		if err := server.ChangeBlock(pos, world.Generic); err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkSmallSparse(b *testing.B) {
	insertPointTemplate(smallSparseDir, b)
}

func BenchmarkMedSparse(b *testing.B) {
	insertPointTemplate(medSparseDir, b)
}

func BenchmarkLgSparse(b *testing.B) {
	insertPointTemplate(largeSparseDir, b)
}

func BenchmarkSmallDense(b *testing.B) {
	insertPointTemplate(smallDenseDir, b)
}

func BenchmarkMedDense(b *testing.B) {
	insertPointTemplate(medDenseDir, b)
}

func BenchmarkLgDense(b *testing.B) {
	insertPointTemplate(largeDenseDir, b)
}
