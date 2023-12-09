package server

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

const fileCacheSize = 8

type SimpleServer struct {
	StorageDir string
	cache      storage.FileCache
}

func (s *SimpleServer) SetStorageRoot(path string) {
	s.StorageDir = path
	s.cache = storage.NewFileCache(256)
}

// Filesystem operations

func (s *SimpleServer) FetchOrCreateChunk(pos world.ChunkPos) (world.ChunkData, error) {
	chunkFileName := filepath.Join(s.StorageDir, pos.ToFileName())

	var chunkData world.ChunkData

	chunkFile, err := os.Open(chunkFileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// There was no chunk that exists, create a blank one
			chunkFile, err = os.Create(chunkFileName)
			if err != nil {
				return chunkData, err
			}

			// Initilize the file with some blank data
			if err := json.NewEncoder(chunkFile).Encode(chunkData); err != nil {
				return chunkData, err
			}

			if _, err := chunkFile.Seek(0, 0); err != nil {
				return chunkData, err
			}
		} else {
			return chunkData, err
		}
	}
	defer chunkFile.Close()

	return storage.ReadChunkFromFile(chunkFile)
}

// `FetchChunk' fetches the chunk's data, given the chunk's position
func (s *SimpleServer) FetchChunk(pos world.ChunkPos) (world.ChunkData, error) {
	chunkFileName := filepath.Join(s.StorageDir, pos.ToFileName())

	var chunkData world.ChunkData

	chunkFile, err := s.cache.FetchFile(chunkFileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return chunkData, storage.ChunkNotFoundError
		} else {
			return chunkData, err
		}
	}

	return storage.ReadChunkFromFile(chunkFile)
}

// Voxel server implementation

func (s *SimpleServer) ChangeBlock(
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

func (s *SimpleServer) ChangeBlockRange(
	targetState world.BlockID,
	start, end world.BlockPos,
) error {
	panic("ChangeBlockRange is unimplemented")
}

func (s *SimpleServer) ReadBlockAt(pos world.BlockPos) (world.BlockID, error) {
	chunk, err := s.FetchChunk(pos.ToChunkPos())
	if err != nil {
		return world.Empty, err
	}

	return chunk.SectionFor(pos).FetchBlock(pos), nil
}

func (s *SimpleServer) ReadChunkAt(pos world.ChunkPos) (world.ChunkData, error) {
	return s.FetchOrCreateChunk(pos)
}
