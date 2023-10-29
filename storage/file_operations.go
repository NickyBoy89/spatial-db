package storage

import (
	"encoding/json"
	"os"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func ReadChunkFromFile(chunkFile *os.File) (world.ChunkData, error) {
	var chunkData world.ChunkData

	if err := json.NewDecoder(chunkFile).Decode(&chunkData); err != nil {
		return chunkData, err
	}

	return chunkData, nil
}
