package server

import (
	"encoding/json"
	"errors"
	"image"
	"io/fs"
	"os"
	"path/filepath"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
	"github.com/NickyBoy89/spatial-db/storage"
	"github.com/NickyBoy89/spatial-db/world"
)

const fileCacheSize = 8

type SimpleServer struct {
	StorageDir     string
	storageBackend storage.UnityFile
}

func (s *SimpleServer) SetStorageRoot(path string) {
	s.StorageDir = path

	var err error
	s.storageBackend, err = storage.OpenUnityFile(path, path+".metadata")
	if err != nil {
		panic(err)
	}
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
	chunkData, err := s.storageBackend.ReadChunk(pos)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return chunkData, storage.ChunkNotFoundError
		} else {
			return chunkData, err
		}
	}

	return chunkData, nil
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

func Abs[T int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func Max[T int](x, y T) T {
	if x > y {
		return x
	}
	return y
}

func Min[T int](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func RoundToNearest16[T int](x T) T {
	for x%16 != 0 {
		if x < 0 {
			x -= 1
		} else {
			x += 1
		}
	}
	return x
}

func PointsInRectHoriz(r image.Rectangle, y uint) []world.BlockPos {
	pts := make([]world.BlockPos, r.Dx()*r.Dy())
	ptIndex := 0

	for x := r.Min.X; x < r.Dx(); x++ {
		for z := r.Min.Y; z < r.Dy(); z++ {
			pts[ptIndex] = world.BlockPos{X: x, Y: y, Z: z}
			ptIndex++
		}
	}

	return pts
}

func CommonPoints(r1, r2 image.Rectangle) []image.Point {
	pts := []image.Point{}

	for x := r1.Min.X; x < r1.Dx(); x++ {
		for z := r1.Min.Y; z < r1.Dy(); z++ {
			testPoint := image.Point{X: x, y: z}
			if testPoint.X > r2.Min.X && testPoint.X < r2.Max.X && testPoint.Y > r2.Min.Y && testPoint.Y < r2.Max.Y {
				pts = append(pts, testPoint)
			}
		}
	}

	return pts
}

func (s *SimpleServer) ChangeBlockRange(
	targetState world.BlockID,
	start, end world.BlockPos,
) error {
	maxx, minx := Max(end.X, start.X), Min(end.X, start.X)
	maxz, minz := Max(end.Z, start.Z), Min(end.Z, start.Z)

	// Create a 2d rectangle looking down at the world
	worldSlice := image.Rectangle{
		Min: image.Point{X: minx, Y: minz},
		Max: image.Point{X: maxx, Y: maxz},
	}

	chunkMinx, chunkMaxx := RoundToNearest16(minx), RoundToNearest16(maxx)
	chunkMinz, chunkMaxz := RoundToNearest16(minz), RoundToNearest16(maxz)

	// Get a list of all the chunks that could intersect
	canditateChunks := []image.Rectangle{}
	for xpos := chunkMinx; xpos < chunkMaxx; xpos += 16 {
		for zpos := chunkMinz; zpos < chunkMaxz; zpos += 16 {
			canditateChunks = append(canditateChunks, image.Rectangle{
				Min: image.Point{X: xpos, Y: zpos},
				Max: image.Point{X: xpos + 16, Y: zpos + 16},
			})
		}
	}

	completeOverlaps := []image.Rectangle{}
	for _, chunkBounds := range canditateChunks {
		if chunkBounds.In(worldSlice) {
			completeOverlaps = append(completeOverlaps, chunkBounds)
		}
	}

	partialOverlaps := []image.Rectangle{}
	for _, chunkBounds := range canditateChunks {
		// Remove chunks that completely overlap
		if chunkBounds.In(worldSlice) {
			continue
		}

		if chunkBounds.Overlaps(worldSlice) {
			partialOverlaps = append(partialOverlaps, chunkBounds)
		}
	}

	// Completely fill chunks that overlap completely
	for _, complete := range completeOverlaps {
		chunk, err := s.FetchChunk(world.ChunkPos{X: complete.Min.X, Z: complete.Min.Y})
		if err != nil {
			return err
		}

		for _, section := range chunk.Sections {
			section.FillSection(targetState)
		}
	}

	// Overwrite chunks that partially overlap
	for _, partial := range partialOverlaps {
		chunk, err := s.FetchChunk(world.ChunkPos{X: partial.Min.X, Z: partial.Min.Y})
		if err != nil {
			return err
		}

		// Loop through the changed points, and overwrite them
		for _, point := range CommonPoints(partial, worldSlice) {
			pos := world.BlockPos{X: point.X, Z: point.Y}
			chunk.SectionFor(pos).UpdateBlock(pos, targetState)
		}
	}

	return nil
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
