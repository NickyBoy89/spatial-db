package storage

import (
	"bytes"
	"encoding/json"
	"io"
	"io/fs"
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
	metadata map[world.ChunkPos]fileMetadata
}

type fileMetadata struct {
	StartOffset int `json:"start_offset"`
	FileSize    int `json:"file_size"`
}

func CreateUnityFile(fileName string) (UnityFile, error) {
	var u UnityFile

	f, err := os.Create(fileName)
	if err != nil {
		return u, err
	}

	u.fd = f
	u.metadata = make(map[world.ChunkPos]fileMetadata)

	return u, nil
}

func OpenUnityFile(fileName, metadataName string) (UnityFile, error) {
	var u UnityFile

	// Read the file
	f, err := os.Open(fileName)
	if err != nil {
		return u, err
	}
	u.fd = f

	if err := u.ReadMetadataFile(metadataName); err != nil {
		return u, err
	}

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
	u.fd.Seek(int64(u.fileSize), io.SeekStart)
	// Write the encoded contents to the file
	if _, err := u.fd.Write(encoded.Bytes()); err != nil {
		return err
	}

	// Update the metadata with the new file
	u.metadata[data.Pos] = fileMetadata{
		StartOffset: u.fileSize,
		FileSize:    encodedSize,
	}
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

func (u UnityFile) ReadChunk(pos world.ChunkPos) (world.ChunkData, error) {
	m, contains := u.metadata[pos]
	if !contains {
		return world.ChunkData{}, fs.ErrNotExist
	}

	u.fd.Seek(int64(m.StartOffset), io.SeekStart)

	fileReader := io.LimitReader(u.fd, int64(m.FileSize))

	var data world.ChunkData
	if err := json.NewDecoder(fileReader).Decode(&data); err != nil {
		return world.ChunkData{}, err
	}

	return data, nil
}

func (u UnityFile) ReadAllChunks() ([]world.ChunkData, error) {
	chunks := []world.ChunkData{}

	for pos := range u.metadata {
		chunk, err := u.ReadChunk(pos)
		if err != nil {
			return nil, err
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func (u *UnityFile) Close() error {
	return u.fd.Close()
}
