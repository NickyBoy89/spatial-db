package main

import (
	"fmt"
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/server"
)

var inmemory server.InMemoryServer
var inmemoryInit bool

func initInMemory() {
	if !inmemoryInit {
		fmt.Println("Initializing in-memory server")
		inmemory.SetStorageRoot("skygrid-save")
		inmemoryInit = true
	}
}

func BenchmarkReadWithin128InMemory(b *testing.B) {
	initInMemory()
	readBlockTemplate(&inmemory, b, 128)
}

func BenchmarkReadWithin512InMemory(b *testing.B) {
	initInMemory()
	readBlockTemplate(&inmemory, b, 512)
}

func BenchmarkReadWithin2048InMemory(b *testing.B) {
	initInMemory()
	readBlockTemplate(&inmemory, b, 2048)
}

func BenchmarkReadWithin65536InMemory(b *testing.B) {
	initInMemory()
	readBlockTemplate(&inmemory, b, 65536)
}
