package storage

import "git.nicholasnovak.io/nnovak/spatial-db/world"

type HashServer struct {
	blocks map[world.BlockPos]world.BlockID
}

func (hs *HashServer) SetStorageRoot(path string) {
	hs.blocks = make(map[world.BlockPos]world.BlockID)
}

func (hs *HashServer) FetchChunk(pos world.ChunkPos) (world.ChunkData, error) {
	panic("Unimplemented")
}

func (hs *HashServer) ChangeBlock(
	worldPosition world.BlockPos,
	targetState world.BlockID,
) error {
	hs.blocks[worldPosition] = targetState
	return nil
}

func (hs *HashServer) ChangeBlockRange(
	targetState world.BlockID,
	start, end world.BlockPos,
) error {
	panic("Unimplemented")
}

func (hs *HashServer) ReadBlockAt(pos world.BlockPos) (world.BlockID, error) {
	panic("Unimplemented")
}

func (hs *HashServer) ReadChunkAt(pos world.ChunkPos) (world.ChunkData, error) {
	var data world.ChunkData
	data.Pos = pos
	for blockPos, state := range hs.blocks {
		if blockPos.ToChunkPos() == pos {
			sec := data.SectionFor(blockPos)
			sec.UpdateBlock(blockPos, state)
		}
	}

	return data, nil
}
