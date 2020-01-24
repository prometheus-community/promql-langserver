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
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/util/strutil"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *server) Completion(ctx context.Context, params *protocol.CompletionParams) (ret *protocol.CompletionList, err error) {
	location, err := s.find(&params.TextDocumentPositionParams)
	if err != nil {
		return nil, nil
	}

	ret = &protocol.CompletionList{}

	completions := &ret.Items

	switch n := location.node.(type) {
	case *promql.Call:
		var name string

		name, err = location.doc.GetSubstring(
			location.query.Pos+token.Pos(location.node.PositionRange().Start),
			location.query.Pos+token.Pos(location.node.PositionRange().End),
		)

		i := 0
		for j, c := range name {
			if 'a' <= c && c <= 'Z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
				i = j
			} else {
				break
			}
		}

		name = name[:i]

		if err != nil {
			return
		}
		if err = s.completeFunctionName(ctx, completions, location, name); err != nil {
			return
		}
	case *promql.VectorSelector:
		metricName := n.Name

		if posRange := n.PositionRange(); int(posRange.End-posRange.Start) == len(n.Name) {
			if err = s.completeFunctionName(ctx, completions, location, metricName); err != nil {
				return
			}
		}
		if location.query.Pos+token.Pos(location.node.PositionRange().Start)+token.Pos(len(metricName)) >= location.pos {
			if err = s.completeMetricName(ctx, completions, location, metricName); err != nil {
				return
			}
		} else {
			if err = s.completeLabels(ctx, completions, location, metricName); err != nil {
				return
			}
		}
	case *promql.AggregateExpr, *promql.BinaryExpr:
		if err = s.completeLabels(ctx, completions, location, ""); err != nil {
			return
		}
	}

	return //nolint: nakedret
}

// nolint:funlen
func (s *server) completeMetricName(ctx context.Context, completions *[]protocol.CompletionItem, location *location, metricName string) error { // nolint:lll
	api := s.getPrometheus()

	var allNames model.LabelValues

	if api != nil {
		var err error

		allNames, _, err = api.LabelValues(ctx, "__name__")
		if err != nil {
			// nolint: errcheck
			s.client.LogMessage(s.lifetime, &protocol.LogMessageParams{
				Type:    protocol.Error,
				Message: errors.Wrapf(err, "could not get metric data from Prometheus").Error(),
			})

			allNames = nil
		}
	}

	editRange, err := getEditRange(location, metricName)
	if err != nil {
		return err
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
			*completions = append(*completions, item)
		}
	}

	queries, err := location.doc.GetQueries()
	if err != nil {
		return err
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
			*completions = append(*completions, item)
		}
	}

	return nil
}

func (s *server) completeFunctionName(_ context.Context, completions *[]protocol.CompletionItem, location *location, metricName string) error { // nolint:lll
	var err error

	editRange, err := getEditRange(location, metricName)
	if err != nil {
		return err
	}

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
			*completions = append(*completions, item)
		}
	}

	for name, desc := range aggregators {
		if strings.HasPrefix(strings.ToLower(name), metricName) {
			item := protocol.CompletionItem{
				Label:            name,
				SortText:         "__1__" + name,
				Kind:             3, //Function
				InsertTextFormat: 2, //Snippet
				Detail:           desc,
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: name + "($1)",
				},
			}
			*completions = append(*completions, item)
		}
	}

	return nil
}

var aggregators = map[string]string{ // nolint:gochecknoglobals
	"sum":          "calculate sum over dimensions",
	"max":          "select maximum over dimensions",
	"min":          "select minimum over dimensions",
	"avg":          "calculate the average over dimensions",
	"stddev":       "calculate population standard deviation over dimensions",
	"stdvar":       "calculate population standard variance over dimensions",
	"count":        "count number of elements in the vector",
	"count_values": "count number of elements with the same value",
	"bottomk":      "smallest k elements by sample value",
	"topk":         "largest k elements by sample value",
	"quantile":     "calculate φ-quantile (0 ≤ φ ≤ 1) over dimensions",
}

// nolint: funlen
func (s *server) completeLabels(ctx context.Context, completions *[]protocol.CompletionItem, location *location, metricName string) error { // nolint:lll
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

	for token.Pos(item.Pos)+token.Pos(len(item.Val))+token.Pos(offset)+location.query.Pos < location.pos {
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
			return nil
		}
	}

	item.Pos += offset

	loc := *location

	if isLabel {
		loc.node = &item
		return s.completeLabel(ctx, completions, &loc, metricName)
	}

	if item.Typ == promql.COMMA || item.Typ == promql.LEFT_PAREN || item.Typ == promql.LEFT_BRACE {
		loc.node = &promql.Item{Pos: item.Pos + 1}
		return s.completeLabel(ctx, completions, &loc, metricName)
	}

	if isValue && lastLabel != "" {
		loc.node = &item
		return s.completeLabelValue(ctx, completions, &loc, lastLabel)
	}

	if item.Typ == promql.EQL || item.Typ == promql.NEQ {
		loc.node = &promql.Item{Pos: item.Pos + promql.Pos(len(item.Val))}
		return s.completeLabelValue(ctx, completions, &loc, lastLabel)
	}

	return nil
}

// nolint:funlen, unparam
func (s *server) completeLabel(ctx context.Context, completions *[]protocol.CompletionItem, location *location, metricName string) error { // nolint:lll
	api := s.getPrometheus()

	var allNames []string

	if api != nil {
		var err error

		allNames, _, err = api.LabelNames(ctx)
		if err != nil {
			// nolint: errcheck
			s.client.LogMessage(s.lifetime, &protocol.LogMessageParams{
				Type:    protocol.Error,
				Message: errors.Wrapf(err, "could not get label data from prometheus").Error(),
			})

			allNames = nil
		}
	}

	editRange, err := getEditRange(location, "")
	if err != nil {
		return err
	}

	for _, name := range allNames {
		if strings.HasPrefix(name, location.node.(*promql.Item).Val) {
			item := protocol.CompletionItem{
				Label: name,
				Kind:  12, //Value
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: name,
				},
			}
			*completions = append(*completions, item)
		}
	}

	return nil
}

// nolint: funlen
func (s *server) completeLabelValue(ctx context.Context, completions *[]protocol.CompletionItem, location *location, labelName string) error { // nolint:lll
	var allNames model.LabelValues

	api := s.getPrometheus()

	if api != nil {
		var err error

		allNames, _, err = api.LabelValues(ctx, labelName)
		if err != nil {
			// nolint: errcheck
			s.client.LogMessage(s.lifetime, &protocol.LogMessageParams{
				Type:    protocol.Error,
				Message: errors.Wrapf(err, "could not get label value data from Prometheus").Error(),
			})
		}
	}

	editRange, err := getEditRange(location, "")
	if err != nil {
		return err
	}

	quoted := location.node.(*promql.Item).Val

	var quote byte

	var unquoted string

	if len(quoted) != 0 {
		quote = quoted[0]

		unquoted, err = strutil.Unquote(quoted)
		if err != nil {
			return nil
		}
	} else {
		quote = '"'
	}

	for _, name := range allNames {
		if strings.HasPrefix(string(name), unquoted) {
			var quoted string

			if quote == '`' {
				if strings.ContainsRune(string(name), '`') {
					quote = '"'
				} else {
					quoted = fmt.Sprint("`", name, "`")
				}
			}

			if quoted == "" {
				quoted = strconv.Quote(string(name))
			}

			if quote == '\'' {
				quoted = quoted[1 : len(quoted)-1]

				quoted = strings.ReplaceAll(quoted, `\"`, `"`)
				quoted = strings.ReplaceAll(quoted, `'`, `\'`)
				quoted = fmt.Sprint("'", quoted, "'")
			}

			item := protocol.CompletionItem{
				Label: quoted,
				Kind:  12, //Value
				TextEdit: &protocol.TextEdit{
					Range:   editRange,
					NewText: quoted,
				},
			}
			*completions = append(*completions, item)
		}
	}

	return nil
}

// getEditRange computes the editRange for a completion. In case the completion area is shorter than
// the node, the oldname of the token to be completed must be provided. The latter mechanism only
// works if oldname is an ASCII string, which can be safely assumed for metric and function names.
func getEditRange(location *location, oldname string) (editRange protocol.Range, err error) {
	editRange.Start, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().Start))
	if err != nil {
		return
	}

	if oldname == "" {
		editRange.End, err = location.doc.PosToProtocolPosition(location.query.Pos + token.Pos(location.node.PositionRange().End))
		if err != nil {
			return
		}
	} else {
		editRange.End = editRange.Start
		editRange.End.Character += float64(len(oldname))
	}

	return //nolint: nakedret
}
