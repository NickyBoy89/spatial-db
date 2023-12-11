package world

import (
	"fmt"
	"math/rand"
)

type BlockPos struct {
	X int  `json:"x"`
	Y uint `json:"y"`
	Z int  `json:"z"`
}

func (b BlockPos) String() string {
	return fmt.Sprintf("BlockPos { X: %v, Y: %v, Z: %v }", b.X, b.Y, b.Z)
}

func RandomBlockPosWithRange(maxRange float64) BlockPos {
	return BlockPos{
		X: int(rand.NormFloat64() * maxRange),
		Y: uint(rand.NormFloat64() * maxRange),
		Z: int(rand.NormFloat64() * maxRange),
	}
}

func (b BlockPos) ToChunkPos() ChunkPos {
	return ChunkPos{
		X: b.X / 16,
		Z: b.Z / 16,
	}
}

type ChunkPos struct {
	X int `json:"x"`
	Z int `json:"z"`
}

func (cp ChunkPos) ToFileName() string {
	return fmt.Sprintf("p.%d.%d.chunk", cp.X, cp.Z)
}

func (cp ChunkPos) StringCoords() string {
	return fmt.Sprintf("%d, %d", cp.X, cp.Z)
}

type BlockID uint8

const (
	Empty BlockID = iota
	Generic
)
