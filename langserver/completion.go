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
	"context"
	"fmt"
	"go/token"
	"os"
	"strings"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *Server) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	fmt.Fprintln(os.Stderr, "0")
	doc, docCtx, err := s.cache.GetDocument(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	var pos token.Pos

	pos, err = doc.ProtocolPositionToTokenPos(docCtx, params.TextDocumentPositionParams.Position)
	if err != nil {
		return nil, nil
	}

	var query *cache.CompiledQuery

	query, err = doc.GetQuery(docCtx, pos-1)
	if err != nil {
		return nil, nil
	}

	node := getSmallestSurroundingNode(query.Ast, pos-1)
	if node == nil {
		return nil, nil
	}

	return s.getCompletions(ctx, node, pos)
}

func (s *Server) getCompletions(ctx context.Context, node promql.Node, pos token.Pos) (*protocol.CompletionList, error) { // nolint:lll
	var metricName string

	switch n := node.(type) {
	case *promql.VectorSelector:
		metricName = n.Name
	case *promql.MatrixSelector:
		metricName = n.Name
	default:
		return nil, nil
	}

	if node.Pos()+token.Pos(len(metricName)) != pos {
		return nil, nil
	}

	if s.prometheus == nil {
		return nil, nil
	}

	api := v1.NewAPI(s.prometheus)

	allNames, _, err := api.LabelValues(ctx, "__name__")
	if err != nil {
		return nil, err
	}

	var items []protocol.CompletionItem

	for _, name := range allNames {
		if strings.HasPrefix(string(name), metricName) {
			item := protocol.CompletionItem{
				Label: string(name),
			}
			items = append(items, item)
		}
	}

	return &protocol.CompletionList{
		IsIncomplete: true,
		Items:        items,
	}, nil
}
