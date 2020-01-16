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

// nolint: lll
package langserver

import (
	"context"
	"fmt"
	"go/token"
	"os"
	"strings"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *server) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	location, err := s.find(&params.TextDocumentPositionParams)
	if err != nil {
		return nil, nil
	}

	var metricName string

	var noLabelSelectors bool

	switch n := location.node.(type) {
	case *promql.VectorSelector:
		metricName = n.Name

		if posRange := n.PositionRange(); int(posRange.End-posRange.Start) == len(n.Name) {
			noLabelSelectors = true
		}
		if location.query.Pos+token.Pos(location.node.PositionRange().Start)+token.Pos(len(metricName)) >= location.pos {
			return s.completeMetricName(ctx, location, metricName, noLabelSelectors)
		}
		return s.completeLabels(ctx, location, metricName)
	default:
		return nil, nil
	}
}

// nolint:funlen
func (s *server) completeMetricName(ctx context.Context, location *location, metricName string, noLabelSelectors bool) (*protocol.CompletionList, error) { // nolint:lll
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

	editRange.Start, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().Start))
	if err != nil {
		return nil, err
	}

	editRange.End, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().Start) + token.Pos(len(metricName)))
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

	queries, err := location.doc.GetQueries()
	if err != nil {
		return nil, err
	}

	for _, q := range queries {
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

// nolint: funlen
func (s *server) completeLabels(ctx context.Context, location *location, metricName string) (*protocol.CompletionList, error) { // nolint:lll
	offset := location.node.PositionRange().Start
	l := promql.Lex(location.query.Content[offset:])

	var (
		item         promql.Item
		lastLabel    string
		insideParen  bool
		insideBraces bool
		isLabel      bool
		isValue      bool
		wantValue    bool
	)

	for token.Pos(item.Pos)+token.Pos(offset)+location.query.Pos < location.pos {
		isLabel = false
		isValue = false

		l.NextItem(&item)

		switch item.Typ {
		case promql.AVG, promql.BOOL, promql.BOTTOMK, promql.BY, promql.COUNT, promql.COUNT_VALUES, promql.GROUP_LEFT, promql.GROUP_RIGHT, promql.IDENTIFIER, promql.IGNORING, promql.LAND, promql.LOR, promql.LUNLESS, promql.MAX, promql.METRIC_IDENTIFIER, promql.MIN, promql.OFFSET, promql.QUANTILE, promql.STDDEV, promql.STDVAR, promql.SUM, promql.TOPK:
			if insideParen || insideBraces {
				lastLabel = item.Val
				isLabel = true
			}
		case promql.EQL, promql.NEQ:
			wantValue = true
		case promql.EQL_REGEX, promql.NEQ_REGEX:
			wantValue = false
		case promql.LEFT_PAREN:
			insideParen = true
			lastLabel = ""
		case promql.RIGHT_PAREN:
			insideParen = false
			lastLabel = ""
		case promql.LEFT_BRACE:
			insideBraces = true
			lastLabel = ""
		case promql.RIGHT_BRACE:
			insideBraces = false
			lastLabel = ""
		case promql.STRING:
			if wantValue {
				isValue = true
			}
		case promql.COMMA:
			lastLabel = ""
		case 0:
			return nil, nil
		}
	}

	item.Pos += offset

	loc := *location

	if isLabel {
		return s.completeLabel(ctx, &loc, metricName)
	}

	if isValue && lastLabel != "" {
		return s.completeLabelValue(ctx, &loc, metricName)
	}

	return nil, nil
}

// nolint:funlen
func (s *server) completeLabel(ctx context.Context, location *location, metricName string) (*protocol.CompletionList, error) { // nolint:lll
	if s.prometheus == nil {
		return nil, nil
	}

	api := v1.NewAPI(s.prometheus)

	allNames, _, err := api.LabelNames(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get label data from prometheus: %s", err.Error())

		allNames = nil
	}

	var editRange protocol.Range

	editRange.Start, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().Start))
	if err != nil {
		return nil, err
	}

	editRange.End, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().Start) + token.Pos(len(metricName)))
	if err != nil {
		return nil, err
	}

	var items []protocol.CompletionItem

	for _, name := range allNames {
		if strings.HasPrefix(name, metricName) {
			item := protocol.CompletionItem{
				Label: name,
				Kind:  12, //Value
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: name,
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

func (s *server) completeLabelValue(ctx context.Context, location *location, labelName string) (*protocol.CompletionList, error) { // nolint:lll
	return nil, nil
}
