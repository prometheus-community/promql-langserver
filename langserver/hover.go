// Copyright 2019 Tobias Guggenmos
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	"github.com/slrtbtfs/promql-lsp/langserver/cache"
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
	doc, docCtx, err := s.cache.GetDocument(params.TextDocumentPositionParams.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	pos, err := doc.ProtocolPositionToTokenPos(docCtx, params.TextDocumentPositionParams.Position)
	if err != nil {
		return nil, err
	}

	markdown := ""

	var compileResult *cache.CompiledQuery

	compileResult, err = doc.GetQuery(docCtx, pos)
	if err != nil {
		return nil, nil
	}

	if compileResult.Err == nil {
		node := getSmallestSurroundingNode(compileResult.Ast, pos)

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

	if call, ok := node.(*promql.Call); ok {
		doc := funcDocStrings(call.Func.Name)

		if _, err := ret.WriteString(doc); err != nil {
			return ""
		}

		if err := ret.WriteByte('\n'); err != nil {
			return ""
		}
	}

	return ret.String()
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
