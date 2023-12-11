package storage

import (
	"os"
	"path"
	"reflect"
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

func TestCreateUnityFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "unity")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an empty file
	u, err := CreateUnityFile(path.Join(tempDir, "test-unity"))
	if err != nil {
		t.Fatalf("Error creating unity file: %v", err)
	}

	if u.Size() != 0 {
		t.Fatalf("Expected size of file to be %v, got %v", 0, u.Size())
	}

	// Save the metadata
	if err := u.WriteMetadataFile(path.Join(tempDir, "test-unity.metadata")); err != nil {
		t.Fatalf("Got an error saving the empty metadata: %v", err)
	}
}

func TestWriteSingleFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "unity")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	u, err := CreateUnityFile(path.Join(tempDir, "test-unity"))
	if err != nil {
		t.Fatalf("Error creating unity file: %v", err)
	}

	// Write a single file
	var data world.ChunkData
	data.Sections[0].BlockStates[0] = 2

	if err := u.WriteChunk(data); err != nil {
		t.Fatalf("Error writing chunk: %v", err)
	}

	// Read the chunk back
	readChunk, err := u.ReadChunk(data.Pos)
	if err != nil {
		t.Fatalf("Error reading chunk: %v", err)
	}

	// Compare the chunks directly
	if !reflect.DeepEqual(data, readChunk) {
		t.Fatalf("Chunks differed, sent %v, received %v", data, readChunk)
	}
}

func TestWriteMultipleFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "unity")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	u, err := CreateUnityFile(path.Join(tempDir, "test-unity"))
	if err != nil {
		t.Fatalf("Error creating unity file: %v", err)
	}

	var (
		chunk1 world.ChunkData
		chunk2 world.ChunkData
		chunk3 world.ChunkData
	)
	chunk1.Pos = world.ChunkPos{
		X: 0,
		Z: 0,
	}
	chunk1.Sections[0].BlockStates[0] = 2
	chunk2.Sections[0].BlockStates[0] = 3
	chunk2.Pos = world.ChunkPos{
		X: 1,
		Z: 0,
	}
	chunk3.Sections[0].BlockStates[0] = 4
	chunk3.Pos = world.ChunkPos{
		X: 2,
		Z: 0,
	}

	chunks := []world.ChunkData{chunk1, chunk2, chunk3}

	// Write all chunks
	for _, data := range chunks {
		if err := u.WriteChunk(data); err != nil {
			t.Fatalf("Error writing chunk: %v", err)
		}
	}

	// Read the chunks back
	for _, data := range chunks {
		readChunk, err := u.ReadChunk(data.Pos)
		if err != nil {
			t.Fatalf("Error reading chunk: %v", err)
		}

		// Compare the chunks directly
		if !reflect.DeepEqual(data, readChunk) {
			t.Fatalf("Chunks differed, sent %v, received %v", data, readChunk)
		}
	}
}
