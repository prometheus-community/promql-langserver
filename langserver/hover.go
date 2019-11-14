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
	"os"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"

	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	// Do not remove! Side effects of init() needed
	_ "github.com/slrtbtfs/promql-lsp/langserver/documentation/functions_statik"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
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
func (s *server) Hover(ctx context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
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

	var hoverRange *protocol.Range

	if compileResult != nil && compileResult.Ast != nil {
		node := getSmallestSurroundingNode(compileResult.Ast, pos)

		markdown = s.nodeToDocMarkdown(ctx, node)

		if node != nil {
			start, err := doc.PosToProtocolPostion(docCtx, node.Pos())
			if err != nil {
				return nil, nil
			}

			end, err := doc.PosToProtocolPostion(docCtx, node.EndPos())
			if err != nil {
				return nil, nil
			}

			hoverRange = &protocol.Range{
				Start: start,
				End:   end,
			}
		}
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  "markdown",
			Value: markdown,
		},
		Range: hoverRange,
	}, nil
}

func (s *server) nodeToDocMarkdown(ctx context.Context, node promql.Node) string {
	var ret bytes.Buffer

	if call, ok := node.(*promql.Call); ok {
		doc := funcDocStrings(call.Func.Name)

		if _, err := ret.WriteString(doc); err != nil {
			return ""
		}

		if err := ret.WriteByte('\n'); err != nil {
			return ""
		}
	}

	if vector, ok := node.(*promql.VectorSelector); ok {
		metric := vector.Name

		doc, err := s.getMetricDocs(ctx, metric)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get metric data: ", err.Error())
		}

		if _, err := ret.WriteString(doc); err != nil {
			return ""
		}
	}

	if matrix, ok := node.(*promql.MatrixSelector); ok {
		metric := matrix.Name

		doc, err := s.getMetricDocs(ctx, metric)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get metric data: ", err.Error())
		}

		if _, err := ret.WriteString(doc); err != nil {
			return ""
		}
	}

	if expr, ok := node.(promql.Expr); ok {
		_, err := ret.WriteString(fmt.Sprintf("__PromQL Type:__ %v\n\n", expr.Type()))
		if err != nil {
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

func (s *server) getMetricDocs(ctx context.Context, metric string) (string, error) {
	var ret strings.Builder

	fmt.Fprintf(&ret, "### %s\n\n", metric)

	if s.prometheus == nil {
		return ret.String(), nil
	}

	api := v1.NewAPI(s.prometheus)

	metadata, err := api.TargetsMetadata(ctx, "", metric, "1")
	if err != nil {
		return ret.String(), err
	} else if len(metadata) == 0 {
		return ret.String(), nil
	}

	if metadata[0].Help != "" {
		fmt.Fprintf(&ret, "__Metric Help:__ %s\n\n", metadata[0].Help)
	}

	if metadata[0].Type != "" {
		fmt.Fprintf(&ret, "__Metric Type:__  %s\n\n", metadata[0].Type)
	}

	if metadata[0].Unit != "" {
		fmt.Fprintf(&ret, "__Metric Unit:__  %s\n\n", metadata[0].Unit)
	}

	return ret.String(), nil
}
