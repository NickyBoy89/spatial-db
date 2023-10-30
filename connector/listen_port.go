package connector

import (
	"bytes"
	"fmt"
	"io"
	"net"

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

func handleConn(conn net.Conn) {
	defer conn.Close()
	log.Infof("Received connection from %v", conn.RemoteAddr())
	// Open a connection to the remote server
	serverConn, err := net.Dial("tcp", fmt.Sprintf(":%d", outputPort))
	if err != nil {
		panic(err)
	}
	defer serverConn.Close()

	var sidecarServerDataStream bytes.Buffer

	// Divert any data read from the server to be copied into the sidecar data stream
	serverReader := io.TeeReader(serverConn, &sidecarServerDataStream)

	go handleGameConnection(&sidecarServerDataStream)

	// Start copying data from the server to the client
	go func() {
		if _, err := io.Copy(conn, serverReader); err != nil {
			panic(err)
		}
	}()

	// Copy data from the client to the server
	if _, err := io.Copy(serverConn, conn); err != nil {
		panic(err)
	}
}
