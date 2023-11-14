package connector

import (
	"errors"
	"fmt"
	"io"
	"net"

	mcnet "github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	inputPort  int
	outputPort int
)

func init() {
	ProxyPortCommand.Flags().IntVar(&inputPort, "in", -1, "The input port, or the port to listen to a server on")
	ProxyPortCommand.Flags().IntVar(&outputPort, "out", -1, "The output port, or the port for the client to listen on")
	ProxyPortCommand.MarkFlagRequired("in")
	ProxyPortCommand.MarkFlagRequired("out")
}

var ProxyPortCommand = &cobra.Command{
	Use:   "proxy",
	Short: "Proxies the connection between the input and output ports",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Listen for connections to the local server
		l, err := net.Listen("tcp", fmt.Sprintf(":%d", inputPort))
		if err != nil {
			return err
		}

		log.Infof("Listening for incoming connections on port %d", inputPort)
		for {
			conn, err := l.Accept()
			if err != nil {
				return err
			}

			go handleConn(conn)
		}
	},
}

func handleConn(clientConn net.Conn) {
	defer log.Info("Closed all connections")
	defer clientConn.Close()
	log.Infof("Received connection from %v", clientConn.RemoteAddr())
	// Open a connection to the remote server
	serverConn, err := net.Dial("tcp", fmt.Sprintf(":%d", outputPort))
	if err != nil {
		log.Errorf("Could not connect to remote server: %v", err)
		return
	}
	defer serverConn.Close()

	// Wrap the server's connection into a mc conn to read packets
	wrappedServerConn := mcnet.WrapConn(serverConn)
	defer wrappedServerConn.Close()

	// Writes to the client connection any data read from the server connection
	wrappedServerConn.Reader = io.TeeReader(wrappedServerConn.Reader, clientConn)

	go func() {
		log.Info("Listening for packets")
		var p pk.Packet
		for {
			if err := wrappedServerConn.ReadPacket(&p); err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				log.Errorf("Error reading packet from server: %v", err)
			}

			// log.Infof("Received packet with id %.2x", p.ID)

			switch p.ID {
			case 0x20:
				panic("Incoming chunk")
			// From here: https://wiki.vg/Protocol#Chunk_Data_and_Update_Light
			case 0x25:
				var (
				// chunkX     pk.Int
				// chunkZ     pk.Int
				// heightmaps struct {
				// 	MotionBlocking []int64       `nbt:"MOTION_BLOCKING"`
				// 	WorldSurface   []pk.NBTField `nbt:"WORLD_SURFACE"`
				// }
				// heightmaps nbt.RawMessage
				// chunkData pk.ByteArray
				)
				panic("Chunk upload called")
			}
		}
	}()

	// Now, we need to copy the network data from the client to the server
	log.Info("Copying data from the client to the server")
	if _, err := io.Copy(wrappedServerConn, clientConn); err != nil {
		log.Errorf("Error copying data to the server: %v", err)
	}
}
