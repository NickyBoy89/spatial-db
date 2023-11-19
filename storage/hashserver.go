package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

type HashServer struct {
	blocks map[world.BlockPos]world.BlockID
}

func (hs *HashServer) SetStorageRoot(path string) {
	hs.blocks = make(map[world.BlockPos]world.BlockID)

	chunkFiles, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for chunkIndex, chunkFile := range chunkFiles {
		var data world.ChunkData

		log.Infof("Reading in chunk %d of %d", chunkIndex, len(chunkFiles))

		f, err := os.Open(filepath.Join(path, chunkFile.Name()))
		if err != nil {
			panic(err)
		}

		// Read each file from disk
		if err := json.NewDecoder(f).Decode(&data); err != nil {
			panic(err)
		}

		// Load in each data point from disk
		for _, section := range data.Sections {
			for blockIndex, blockState := range section.BlockStates {
				pos := data.IndexToBlockPos(blockIndex)
				hs.blocks[pos] = blockState
			}
		}

		f.Close()
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
