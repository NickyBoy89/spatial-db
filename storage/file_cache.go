package storage

import (
	"math/rand"
	"os"
)

// A `FileCache` stores open file descriptors for all files that have been
// requested, and returns them when requested
type FileCache struct {
	cacheSize      int
	fileNames      map[string]int
	indexesOfFiles map[int]string
	files          []*os.File
}

func NewFileCache(cacheSize int) FileCache {
	var c FileCache

	c.cacheSize = cacheSize
	c.fileNames = make(map[string]int)
	c.indexesOfFiles = make(map[int]string)
	c.files = make([]*os.File, cacheSize)

	return c
}

// `FetchFile` takes in a file's name and returns its file descriptor
//
// If the file is not found in the cache, the file is added to the cache
func (fc *FileCache) FetchFile(fileName string) (*os.File, error) {
	// The file is already in the cache
	if index, contains := fc.fileNames[fileName]; contains {
		return fc.files[index], nil
	}

	// The file was not in the cache, try and insert it

	// Random cache eviction
	indexToRemove := rand.Intn(fc.cacheSize)

	// Remove the old value
	oldName, present := fc.indexesOfFiles[indexToRemove]

	// If there is already a value, delete its mapping
	if present {
		delete(fc.fileNames, oldName)
		delete(fc.indexesOfFiles, indexToRemove)
	}

	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	// Insert the new value
	fc.files[indexToRemove] = f
	fc.fileNames[fileName] = indexToRemove
	fc.indexesOfFiles[indexToRemove] = fileName

	return f, nil
}
