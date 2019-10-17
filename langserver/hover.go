package langserver

import (
	"context"
	"fmt"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/prometheus/promql"
)

func (s *Server) Hover(_ context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc, err := s.cache.getDocument(params.TextDocumentPositionParams.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// FIXME: This is still a bit racy
	doc.compilers.Wait()
	doc.Mu.RLock()
	defer doc.Mu.RUnlock()
	pos, err := doc.protocolPositionToTokenPos(params.TextDocumentPositionParams.Position)
	if err != nil {
		return nil, err
	}
	node := getSmallestSourroundingNode(doc.compileResult.ast, pos)

	expr, ok := node.(promql.Expr)

	if ok {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Type: %v", expr.Type()),
			},
		}, nil
	} else {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: fmt.Sprintf("Compile Error"),
			},
		}, nil
	}
}
