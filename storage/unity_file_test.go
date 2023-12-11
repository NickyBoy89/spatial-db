package storage

import (
	"os"
	"path"
	"testing"
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
