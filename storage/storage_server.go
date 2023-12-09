package storage

import (
	"errors"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

type StorageServer interface {
	// Individual operations
	SetStorageRoot(path string)

	// Block-level interactions
	ChangeBlock(targetState world.BlockID, world_position world.BlockPos) error
	ReadBlockAt(pos world.BlockPos) (world.BlockID, error)

	// Region-level interactions
	ChangeBlockRange(targetState world.BlockID, start, end world.BlockPos) error

	// Network-level operations
	ReadChunkAt(pos world.ChunkPos) error
}

var (
	ChunkNotFoundError = errors.New("chunk was not found in storage")
)
