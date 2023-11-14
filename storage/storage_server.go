package storage

import "git.nicholasnovak.io/nnovak/spatial-db/world"

type StorageServer interface {
	// Individual operations
	SetStorageRoot(path string)

	// Individual block-level interactions
	ChangeBlock(targetState world.BlockID, world_position world.BlockPos) error
	ChangeBlockRange(targetState world.BlockID, start, end world.BlockPos) error
	ReadBlockAt(pos world.BlockPos) error

	// Network-level operations
	ReadChunkAt(pos world.ChunkPos) error
}
