package main

import (
	"fmt"
	"testing"

	"git.nicholasnovak.io/nnovak/spatial-db/server"
)

var hash server.HashServer
var hashInit bool

func initHash() {
	if !hashInit {
		fmt.Println("Initializing hash server")
		hash.SetStorageRoot("skygrid-save")
		hashInit = true
	}
}

func BenchmarkReadWithin128Hashtable(b *testing.B) {
	initHash()
	readBlockTemplate(&hash, b, 128)
}

func BenchmarkReadWithin512Hashtable(b *testing.B) {
	initHash()
	readBlockTemplate(&hash, b, 512)
}

func BenchmarkReadWithin2048Hashtable(b *testing.B) {
	initHash()
	readBlockTemplate(&hash, b, 2048)
}

func BenchmarkReadWithin65536Hashtable(b *testing.B) {
	initHash()
	readBlockTemplate(&hash, b, 65536)
}
