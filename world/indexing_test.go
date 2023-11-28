package world

import (
	"testing"
)

func TestUniqueDataPoints(t *testing.T) {
	points := make(map[int]bool)

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {

				pos := BlockPos{
					X: x,
					Y: uint(y),
					Z: z,
				}

				points[IndexOfBlock(pos)] = true
			}
		}
	}

	if len(points) != 4096 {
		t.Fatalf("Expected %d unique points, got %d", 4096, len(points))
	}
}

func TestCorrectIndexReversal(t *testing.T) {
	points := make(map[int]BlockPos)

	for x := 0; x < 16; x++ {
		for y := 0; y < 16; y++ {
			for z := 0; z < 16; z++ {

				pos := BlockPos{
					X: x,
					Y: uint(y),
					Z: z,
				}

				points[IndexOfBlock(pos)] = pos
			}
		}
	}

	var chunk ChunkData

	for index, blockPos := range points {
		testBlock := chunk.IndexToBlockPos(index)
		if testBlock != blockPos {
			t.Fatalf("Expected block %v, got %v", blockPos, testBlock)
		}
	}
}
