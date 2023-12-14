package server

import (
	"github.com/NickyBoy89/spatial-db/storage"
	"github.com/NickyBoy89/spatial-db/world"
	log "github.com/sirupsen/logrus"
)

type HashServer struct {
	blocks map[world.BlockPos]world.BlockID
}

func (hs *HashServer) SetStorageRoot(path string) {
	hs.blocks = make(map[world.BlockPos]world.BlockID)

	u, err := storage.OpenUnityFile(path, path+".metadata")
	if err != nil {
		panic(err)
	}
	defer u.Close()

	chunks, err := u.ReadAllChunks()
	if err != nil {
		panic(err)
	}

	for chunkIndex, data := range chunks {
		// Load in each data point from disk
		log.Infof("Reading in chunk %d of %d", chunkIndex, len(chunks))

		for _, section := range data.Sections {
			for blockIndex, blockState := range section.BlockStates {
				pos := data.IndexToBlockPos(blockIndex)
				hs.blocks[pos] = section.Palette.State(blockState)
			}
		}
	}
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
	return hs.blocks[pos], nil
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
