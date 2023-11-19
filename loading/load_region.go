package loading

import (
	"path/filepath"
	"strconv"
	"strings"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
	"github.com/Tnze/go-mc/save"
	"github.com/Tnze/go-mc/save/region"
)

// LoadRegionFile loads a single region file into an array of chunks
//
// A region is a 32x32 grid of chunks, although the final output can store less
func LoadRegionFile(fileName string) ([]world.ChunkData, error) {
	regionFile, err := region.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer regionFile.Close()

	// Parse the name of the region to find its position within the world
	nameParts := strings.Split(filepath.Base(fileName), ".")
	regionX, err := strconv.Atoi(nameParts[1])
	if err != nil {
		return nil, err
	}
	regionY, err := strconv.Atoi(nameParts[2])
	if err != nil {
		return nil, err
	}

	// A region file is a 32x32 grid of chunks
	chunks := []world.ChunkData{}

	for i := 0; i < 32; i++ {
		for j := 0; j < 32; j++ {
			if regionFile.ExistSector(i, j) {
				sectorFile, err := regionFile.ReadSector(i, j)
				if err != nil {
					return nil, err
				}

				// Read each chunk from disk
				var chunk save.Chunk
				if err := chunk.Load(sectorFile); err != nil {
					return nil, err
				}

				// Convert each chunk into the database's format
				var chunkData world.ChunkData
				chunkData.FromMCAChunk(chunk)

				chunkData.Pos.X = regionX*32 + i
				chunkData.Pos.Z = regionY*32 + j

				chunks = append(chunks, chunkData)
			}
		}
	}

	return chunks, nil
}
