package storage

import (
	"errors"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

type StorageServer interface {
	// Individual operations
	SetStorageRoot(path string)

	// Block-level interactions
	ChangeBlock(pos world.BlockPos, targetState world.BlockID) error
	ReadBlockAt(pos world.BlockPos) (world.BlockID, error)

	// Region-level interactions
	ChangeBlockRange(targetState world.BlockID, start, end world.BlockPos) error

	// Network-level operations
	ReadChunkAt(pos world.ChunkPos) (world.ChunkData, error)
}

var (
	ChunkNotFoundError = errors.New("chunk was not found in storage")
)
