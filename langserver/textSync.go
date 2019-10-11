package langserver

import (
	"context"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

func (s *Server) DidOpen(_ context.Context, params *protocol.DidOpenTextDocumentParams) error {
	return s.cache.addDocument(&params.TextDocument)
}

func (s *Server) DidClose(_ context.Context, params *protocol.DidCloseTextDocumentParams) error {
	return s.cache.removeDocument(params.TextDocument.URI)
}
