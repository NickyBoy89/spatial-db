package world

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

func (cp ChunkPos) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%d %d", cp.X, cp.Z)), nil
}

func (cp *ChunkPos) UnmarshalText(text []byte) error {
	words := strings.Split(string(text), " ")
	x, err := strconv.Atoi(words[0])
	if err != nil {
		return err
	}

	z, err := strconv.Atoi(words[1])
	if err != nil {
		return err
	}

	cp.X = x
	cp.Z = z

	return nil
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

func (id BlockID) String() string {
	switch id {
	case Empty:
		return "Empty"
	case Generic:
		return "Generic"
	default:
		panic(fmt.Sprintf("Unknown block id: %v", uint8(id)))
	}
}
