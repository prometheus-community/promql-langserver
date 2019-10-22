package langserver

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/prometheus/promql"

	// Do not remove! Side effects of init() needed
	_ "github.com/slrtbtfs/promql-lsp/langserver/documentation/functions_statik"
)

var functionDocumentationFS http.FileSystem

func init() {
	var err error
	functionDocumentationFS, err = fs.New()
	if err != nil {
		log.Fatal(err)
	}
}

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

	markdown := ""
	if doc.compileResult.err == nil {
		node := getSmallestSourroundingNode(doc.compileResult.ast, pos)

		markdown = nodeToDocMarkdown(node)
	}
	if markdown != "" {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: markdown,
			},
		}, nil
	} else {
		return nil, nil
	}
}

func nodeToDocMarkdown(node promql.Node) string {
	var ret bytes.Buffer
	expr, ok := node.(promql.Expr)
	if ok {
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
	return string(ret.Bytes())
}

func funcDocStrings(name string) string {
	name = strings.ToLower(name)
	file, err := functionDocumentationFS.Open(fmt.Sprintf("/%s.md", name))
	if err != nil {
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
