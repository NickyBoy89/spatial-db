package connector

import (
	"io"
)

// This is the size that Netty apparently uses to read from the incoming
// connection at a time
//
// This is documented in a mod that removes this limit:
// https://www.curseforge.com/minecraft/mc-mods/xl-packets
const nettyMaxPacketSize = 2_097_152 // 16MiB

func handleGameConnection(conn io.Reader) {
	panic("Unhandled game traffic")
}
