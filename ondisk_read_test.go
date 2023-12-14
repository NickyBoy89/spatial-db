package main

import (
	"fmt"
	"testing"

	"github.com/NickyBoy89/spatial-db/server"
)

var disk server.SimpleServer
var diskInit bool

func initDisk() {
	if !diskInit {
		fmt.Println("Initializing disk server")
		disk.SetStorageRoot("skygrid-save")
		diskInit = true
	}
}

func BenchmarkReadWithin128OnDisk(b *testing.B) {
	initDisk()
	readBlockTemplate(&disk, b, 128)
}

func BenchmarkReadWithin512OnDisk(b *testing.B) {
	initDisk()
	readBlockTemplate(&disk, b, 512)
}

func BenchmarkReadWithin2048OnDisk(b *testing.B) {
	initDisk()
	readBlockTemplate(&disk, b, 2048)
}

func BenchmarkReadWithin65536OnDisk(b *testing.B) {
	initDisk()
	readBlockTemplate(&disk, b, 65536)
}
