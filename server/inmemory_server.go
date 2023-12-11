package server

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

type InMemoryServer struct {
	StorageDir string
	Chunks     map[world.ChunkPos]world.ChunkData
}

func (s *InMemoryServer) SetStorageRoot(path string) {
	s.StorageDir = path

	chunkFiles, err := os.ReadDir(s.StorageDir)
	if err != nil {
		panic(err)
	}

	s.Chunks = make(map[world.ChunkPos]world.ChunkData)

	validChunkFiles := []fs.DirEntry{}
	for _, chunkFile := range chunkFiles {
		if chunkFile.IsDir() || !strings.HasSuffix(chunkFile.Name(), ".chunk") {
			continue
		}
		validChunkFiles = append(validChunkFiles, chunkFile)
	}

	chunks := make([]world.ChunkData, len(validChunkFiles))

	for chunkIndex, chunkFile := range validChunkFiles {
		go func(index int, cf fs.DirEntry) {
			file, err := os.Open(filepath.Join(s.StorageDir, cf.Name()))
			if err != nil {
				panic(err)
			}

			chunkData, err := storage.ReadChunkFromFile(file)
			if err != nil {
				panic(err)
			}

			file.Close()

			chunks[index] = chunkData
		}(chunkIndex, chunkFile)
	}

	for _, chunkData := range chunks {
		s.Chunks[chunkData.Pos] = chunkData
	}
}

func (s *InMemoryServer) FetchOrCreateChunk(pos world.ChunkPos) (world.ChunkData, error) {
	return s.Chunks[pos], nil
}

func (s *InMemoryServer) FetchChunk(pos world.ChunkPos) (world.ChunkData, error) {
	chunkFileName := filepath.Join(s.StorageDir, pos.ToFileName())

	var chunkData world.ChunkData

	chunkFile, err := os.Open(chunkFileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return chunkData, storage.ChunkNotFoundError
		} else {
			return chunkData, err
		}
	}
	defer chunkFile.Close()

	return storage.ReadChunkFromFile(chunkFile)
}

// Voxel server implementation

func (s *InMemoryServer) ChangeBlock(
	worldPosition world.BlockPos,
	targetState world.BlockID,
) error {
	chunk, err := s.FetchOrCreateChunk(worldPosition.ToChunkPos())
	if err != nil {
		return err
	}

	chunk.SectionFor(worldPosition).UpdateBlock(worldPosition, targetState)

	return nil
}

func (s *InMemoryServer) ChangeBlockRange(
	targetState world.BlockID,
	start, end world.BlockPos,
) error {
	panic("ChangeBlockRange is unimplemented")
}

func (s *InMemoryServer) ReadBlockAt(pos world.BlockPos) (world.BlockID, error) {
	chunk, err := s.FetchOrCreateChunk(pos.ToChunkPos())
	if err != nil {
		return world.Empty, err
	}

	return chunk.SectionFor(pos).FetchBlock(pos), nil
}

func (s *InMemoryServer) ReadChunkAt(pos world.ChunkPos) (world.ChunkData, error) {
	return s.FetchOrCreateChunk(pos)
}
