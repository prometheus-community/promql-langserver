package langserver

import (
	"context"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

func (s *Server) DidOpen(_ context.Context, params *protocol.DidOpenTextDocumentParams) error {
	s.cache.addDocument(&params.TextDocument)
	return nil
}

func (s *Server) DidClose(_ context.Context, params *protocol.DidCloseTextDocumentParams) error {
	s.cache.removeDocument(params.TextDocument.URI)
	return nil
}
