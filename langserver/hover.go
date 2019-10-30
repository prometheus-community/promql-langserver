package langserver

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"

	// Do not remove! Side effects of init() needed
	_ "github.com/slrtbtfs/promql-lsp/langserver/documentation/functions_statik"
)

//nolint: gochecknoglobals
var functionDocumentationFS = initializeFunctionDocumentation()

func initializeFunctionDocumentation() http.FileSystem {
	ret, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	return ret
}

// Hover shows documentation on hover
// required by the protocol.Server interface
func (s *Server) Hover(_ context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc, err := s.cache.getDocument(params.TextDocumentPositionParams.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// FIXME: This is still a bit racy
	doc.compilers.Wait()
	doc.Mu.RLock()

	defer doc.Mu.RUnlock()

	if pos, err := doc.protocolPositionToTokenPos(params.TextDocumentPositionParams.Position); err != nil {
		return nil, err
	}

	markdown := ""

	if doc.compileResult.err == nil {
		node := getSmallestSurroundingNode(doc.compileResult.ast, pos)

		markdown = nodeToDocMarkdown(node)
	}

	if markdown != "" {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: markdown,
			},
		}, nil
	}

	return nil, nil
}

func nodeToDocMarkdown(node promql.Node) string {
	var ret bytes.Buffer

	if expr, ok := node.(promql.Expr); ok {
		_, err := ret.WriteString(fmt.Sprintf("Type: %v\n", expr.Type()))
		if err != nil {
			return ""
		}
	}

	var call *promql.Call
	call, ok = node.(*promql.Call)
	if ok {
		doc := funcDocStrings(call.Func.Name)
		_, err := ret.WriteString(doc)
		if err != nil {
			return ""
		}
		err = ret.WriteByte('\n')
		if err != nil {
			return ""
		}
	}
	return ret.String()
}

func funcDocStrings(name string) string {
	name = strings.ToLower(name)

	if file, err := functionDocumentationFS.Open(fmt.Sprintf("/%s.md", name)); err != nil {
		return ""
	}

	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return ""
	}
	ret := make([]byte, stat.Size())
	_, err = file.Read(ret)
	if err != nil {
		return ""
	}
	return string(ret)
}
