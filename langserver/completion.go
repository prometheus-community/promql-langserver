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
	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *server) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	fmt.Fprintln(os.Stderr, "0")
	doc, err := s.cache.GetDocument(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	var pos token.Pos

	pos, err = doc.ProtocolPositionToTokenPos(params.TextDocumentPositionParams.Position)
	if err != nil {
		return nil, nil
	}

	var query *cache.CompiledQuery

	query, err = doc.GetQuery(pos)
	if err != nil {
		return nil, nil
	}

	node := getSmallestSurroundingNode(query.Ast, pos)
	if node == nil {
		return nil, nil
	}

	return s.getCompletions(ctx, doc, node, pos)
}

// nolint:funlen
func (s *server) getCompletions(ctx context.Context, doc *cache.DocumentHandle, node promql.Node, pos token.Pos) (*protocol.CompletionList, error) { // nolint:lll
	var metricName string

	var noLabelSelectors bool

	switch n := node.(type) {
	case *promql.VectorSelector:
		metricName = n.Name

		if n.LBrace == token.NoPos {
			noLabelSelectors = true
		}
	case *promql.MatrixSelector:
		metricName = n.Name
	default:
		return nil, nil
	}

	if node.Pos()+token.Pos(len(metricName)) >= pos {
		return s.completeMetricName(ctx, doc, node, metricName, noLabelSelectors)
	}

	return nil, nil
}

// nolint:funlen
func (s *server) completeMetricName(ctx context.Context, doc *cache.DocumentHandle, node promql.Node, metricName string, noLabelSelectors bool) (*protocol.CompletionList, error) { // nolint:lll
	if s.prometheus == nil {
		return nil, nil
	}

	api := v1.NewAPI(s.prometheus)

	allNames, _, err := api.LabelValues(ctx, "__name__")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get metric data from prometheus: %s", err.Error())

		allNames = nil
	}

	var editRange protocol.Range

	editRange.Start, err = doc.PosToProtocolPosition(node.Pos())
	if err != nil {
		return nil, err
	}

	editRange.End, err = doc.PosToProtocolPosition(node.Pos() + token.Pos(len(metricName)))
	if err != nil {
		return nil, err
	}

	var items []protocol.CompletionItem

	if noLabelSelectors {
		for name := range promql.Functions {
			if strings.HasPrefix(strings.ToLower(name), metricName) {
				item := protocol.CompletionItem{
					Label:            name,
					SortText:         "__1__" + name,
					Kind:             3, //Function
					InsertTextFormat: 2, //Snippet
					TextEdit: &protocol.TextEdit{
						Range:   editRange,
						NewText: name + "($1)",
					},
					Command: &protocol.Command{
						// This might create problems with non VS Code clients
						Command: "editor.action.triggerParameterHints",
					},
				}
				items = append(items, item)
			}
		}
	}

	for _, name := range allNames {
		if strings.HasPrefix(string(name), metricName) {
			item := protocol.CompletionItem{
				Label:    string(name),
				SortText: "__3__" + string(name),
				Kind:     12, //Value
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: string(name),
				},
			}
			items = append(items, item)
		}
	}

	querys, err := doc.GetQueries()
	if err != nil {
		return nil, err
	}

	for _, q := range querys {
		if rec := q.Record; rec != "" && strings.HasPrefix(rec, metricName) {
			item := protocol.CompletionItem{
				Label:            rec,
				SortText:         "__2__" + rec,
				Kind:             3, //Value
				InsertTextFormat: 2, //Snippet
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: rec,
				},
			}
			items = append(items, item)
		}
	}

	return &protocol.CompletionList{
		IsIncomplete: true,
		Items:        items,
	}, nil
}
