package langserver

import (
	"context"
	"os"
	"sync"

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.Client

	state   serverState
	stateMu sync.Mutex
}

type serverState int

const (
	serverCreated      = serverState(iota)
	serverInitializing // set once the server has received "initialize" request
	serverInitialized  // set once the server has received "initialized" request
	serverShutDown
)

func (s *Server) Run(ctx context.Context) error {
	return s.Conn.Run(ctx)
}

func ServerFromStream(ctx context.Context, stream jsonrpc2.Stream) (context.Context, *Server) {
	s := &Server{}
	ctx, s.Conn, s.client = protocol.NewServer(ctx, stream, s)
	return ctx, s
}

func StdioServer(ctx context.Context) (context.Context, *Server) {
	stream := jsonrpc2.NewHeaderStream(os.Stdin, os.Stdout)
	stream = protocol.LoggingStream(stream, os.Stderr)
	return ServerFromStream(ctx, stream)
}
