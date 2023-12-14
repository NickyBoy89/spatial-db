package storage

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/NickyBoy89/spatial-db/world"
)

func ReadChunkFromFile(chunkFile *os.File) (world.ChunkData, error) {
	var chunkData world.ChunkData

	if err := json.NewDecoder(chunkFile).Decode(&chunkData); err != nil {
		return chunkData, err
	}

	return chunkData, nil
}

func ReadParallelFromDirectory(dirName string) ([]world.ChunkData, error) {
	chunkFiles, err := os.ReadDir(dirName)
	if err != nil {
		panic(err)
	}

	// Filter invalid chunks

	validChunkFiles := []fs.DirEntry{}
	for _, chunkFile := range chunkFiles {
		if chunkFile.IsDir() || !strings.HasSuffix(chunkFile.Name(), ".chunk") {
			continue
		}
		validChunkFiles = append(validChunkFiles, chunkFile)
	}

	chunks := make([]world.ChunkData, len(validChunkFiles))

	var wg sync.WaitGroup
	wg.Add(len(validChunkFiles))

	for fileIndex, chunkFile := range validChunkFiles {
		// Avoid implicit copies
		chunkFile := chunkFile
		fileIndex := fileIndex

		go func() {
			defer wg.Done()
			file, err := os.Open(filepath.Join(dirName, chunkFile.Name()))
			if err != nil {
				panic(err)
			}

			chunkData, err := ReadChunkFromFile(file)
			if err != nil {
				panic(err)
			}

			file.Close()

			chunks[fileIndex] = chunkData
		}()
	}

	wg.Wait()

	return chunks, nil
}
