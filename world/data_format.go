package world

import (
	"fmt"
	"math/rand"
)

const (
	// Slice size is the total number of blocks in a horizontal slice of a chunk
	sliceSize = 16 * 16
)

type BlockPos struct {
	X int  `json:"x"`
	Y uint `json:"y"`
	Z int  `json:"z"`
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

type ChunkSection struct {
	// The count of full blocks in the chunk
	BlockCount  uint                  `json:"block_count"`
	BlockStates [16 * 16 * 16]BlockID `json:"block_states"`
}

func rem_euclid(a, b int) int {
	return (a%b + b) % b
}

func (cs *ChunkSection) IndexOfBlock(pos BlockPos) int {
	baseX := rem_euclid(pos.X, 16)
	baseY := rem_euclid(int(pos.Y), 16)
	baseZ := rem_euclid(pos.Z, 16)

	return (baseY * sliceSize) + (baseZ * 16) + baseX
}

func (cs *ChunkSection) UpdateBlockAtIndex(index int, targetState BlockID) {
	// TODO: Keep track of the block count

	cs.BlockStates[index] = targetState
}

func (cs *ChunkSection) UpdateBlock(pos BlockPos, targetState BlockID) {
	cs.BlockStates[cs.IndexOfBlock(pos)] = targetState
}

func (cs *ChunkSection) FetchBlock(pos BlockPos) BlockID {
	return cs.BlockStates[cs.IndexOfBlock(pos)]
}

type BlockID uint8

const (
	Empty BlockID = iota
	Generic
)

func (b *BlockID) UnmarshalJSON(data []byte) error {
	idName := string(data)

	if len(idName) < 2 {
		return fmt.Errorf("error decoding blockid, input was too short")
	}

	switch idName[1 : len(idName)-1] {
	case "Empty":
		*b = Empty
	case "Generic":
		*b = Generic
	default:
		return fmt.Errorf("unknown block id: %s", string(data))
	}

	return nil
}

func (b BlockID) MarshalJSON() ([]byte, error) {
	var encoded []byte

	switch b {
	case Empty:
		encoded = []byte("\"Empty\"")
	case Generic:
		encoded = []byte("\"Generic\"")
	default:
		return []byte{}, fmt.Errorf("could not turn block id %d into data", b)
	}

	return encoded, nil
}
