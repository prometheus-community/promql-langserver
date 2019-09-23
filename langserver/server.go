package langserver

import (
	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
)

type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.client
}
