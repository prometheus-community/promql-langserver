package langserver

import (
	"context"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

func (s *Server) Hover(_ context.Context, _ *protocol.HoverParams) (*protocol.Hover, error) {
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: "Example *Hovertext*",
		},
	}, nil
}
