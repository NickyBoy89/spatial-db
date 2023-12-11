package storage

import (
	"bytes"
	"encoding/json"
	"os"

	"git.nicholasnovak.io/nnovak/spatial-db/world"
)

// A `UnityFile` is a collection of chunks, stored as a single file on disk
//
// This is done because with worlds that span millions of chunks, the speed to
// access individual files slows down
type UnityFile struct {
	fd       *os.File
	fileSize int
	// metadata maps the position of a chunk to its start index within the file
	metadata map[world.ChunkPos]int
}

func CreateUnityFile(fileName string) (UnityFile, error) {
	var u UnityFile

	f, err := os.Create(fileName)
	if err != nil {
		return u, err
	}

	u.fd = f
	u.metadata = make(map[world.ChunkPos]int)

	return u, nil
}

func (u UnityFile) Size() int {
	return u.fileSize
}

func (u *UnityFile) WriteChunk(data world.ChunkData) error {
	var encoded bytes.Buffer

	// Encode the chunk first
	if err := json.NewEncoder(&encoded).Encode(data); err != nil {
		return err
	}

	encodedSize := encoded.Len()

	// Go to the end of the file
	u.fd.Seek(0, u.fileSize)
	// Write the encoded contents to the file
	if _, err := u.fd.Write(encoded.Bytes()); err != nil {
		return err
	}

	// Update the metadata with the new file
	u.metadata[data.Pos] = u.fileSize
	u.fileSize += encodedSize

	return nil
}

func (u UnityFile) WriteMetadataFile(fileName string) error {
	fd, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fd.Close()

	if err := json.NewEncoder(fd).Encode(&u.metadata); err != nil {
		return err
	}

	return nil
}

func (u *UnityFile) ReadMetadataFile(fileName string) error {
	fd, err := os.Open(fileName)
	if err != nil {
		return err
	}

	if err := json.NewDecoder(fd).Decode(&u.metadata); err != nil {
		return err
	}

	return nil
}

func (u UnityFile) ReadChunk() world.ChunkData {
	return world.ChunkData{}
}
