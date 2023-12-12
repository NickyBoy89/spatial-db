package server

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

type InMemoryServer struct {
	StorageDir string
	Chunks     map[world.ChunkPos]world.ChunkData
}

func (s *InMemoryServer) SetStorageRoot(path string) {
	s.StorageDir = path

	u, err := storage.OpenUnityFile(s.StorageDir, s.StorageDir+".metadata")
	if err != nil {
		panic(err)
	}
	defer u.Close()

	chunks, err := u.ReadAllChunks()
	if err != nil {
		panic(err)
	}

	s.Chunks = make(map[world.ChunkPos]world.ChunkData, len(chunks))
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
