package langserver

import (
	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

type Server struct {
	Conn   *jsonrpc2.Conn
	client protocol.client
}

func notImplemented(method string) *jsonrpc2.Error {
	return jsonrpc2.NewErrorf("method %q no yet implemented", method)
}
