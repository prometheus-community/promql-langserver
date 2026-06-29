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
	"errors"

	promql "github.com/prometheus/prometheus/promql/parser"
	"github.com/prometheus/prometheus/promql/parser/posrange"
	"go.lsp.dev/protocol"
)

// Parameter labels reused across multiple PromQL function signatures.
const (
	paramInstantVector     = "v instant-vector"
	paramRangeVector       = "v range-vector"
	paramVectorTimeInstant = "v=vector(time()) instant-vector"
)

// rateFuncName is the name of the PromQL rate function. It is also referenced
// from the tests.
const rateFuncName = "rate"

// SignatureHelp is required by the protocol.Server interface.
func (s *server) SignatureHelp(ctx context.Context, params *protocol.SignatureHelpParams) (*protocol.SignatureHelp, error) {
	location, err := s.cache.Find(&params.TextDocumentPositionParams)
	if err != nil {
		return nil, nil
	}

	call, ok := location.Node.(*promql.Call)
	if !ok {
		return nil, nil
	}

	signature, err := getSignature(call.Func.Name)
	if err != nil {
		return nil, nil
	}

	activeParameter := 0.

	for i, arg := range call.Args {
		if arg != nil && arg.PositionRange().End < posrange.Pos(location.Pos-location.Query.Pos) {
			activeParameter = float64(i) + 1
		}
	}

	// For the label_join function, which has a variable number of arguments,
	// the "..." should be highlighted at some point.
	// For reference, the signature is:
	// label_join(v instant-vector, dst_label string, separator string, src_label_1 string, src_label_2 string, ...)
	if call.Func.Name == "label_join" && activeParameter >= 5 {
		activeParameter = 5
	}

	response := &protocol.SignatureHelp{
		Signatures:      []protocol.SignatureInformation{signature},
		ActiveParameter: activeParameter,
	}

	return response, nil
}

// nolint: funlen
func getSignature(name string) (protocol.SignatureInformation, error) {
	signatures := map[string]protocol.SignatureInformation{
		"abs": {
			Label: "abs(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"absent": {
			Label: "absent(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"ceil": {
			Label: "ceil(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"clamp_max": {
			Label: "clamp_max(v instant-vector, max scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
				{Label: "max scalar"},
			},
		},
		"clamp_min": {
			Label: "clamp_min(v instant-vector, min scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
				{Label: "min scalar"},
			},
		},
		"day_of_month": {
			Label: "day_of_month(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"day_of_week": {
			Label: "day_of_week(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"day_in_month": {
			Label: "day_in_month(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"delta": {
			Label: "delta(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"deriv": {
			Label: "deriv(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"exp": {
			Label: "exp(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"floor": {
			Label: "floor(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"histogram_quantile": {
			Label: "histogram_quantile(φ float, b instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: "φ float"},
				{Label: "b instant-vector"},
			},
		},
		"holt_winters": {
			Label: "holt_winters(v range-vector, sf scalar, tf scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
				{Label: "sf scalar"},
				{Label: "tf scalar"},
			},
		},
		"hour": {
			Label: "hour(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"idelta": {
			Label: "idelta(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"increase": {
			Label: "increase(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"irate": {
			Label: "irate(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"label_join": {
			Label: "label_join(v instant-vector, dst_label string, separator string, src_label_1 string, src_label_2 string, ...)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
				{Label: "dst_label string"},
				{Label: "separator string"},
				{Label: "src_label_1 string"},
				{Label: "src_label_2 string"},
				{Label: "..."},
			},
		},
		"label_replace": {
			Label: "label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
				{Label: "dst_label string"},
				{Label: "replacement string"},
				{Label: "src_label string"},
				{Label: "regex string"},
			},
		},
		"ln": {
			Label: "ln(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"log2": {
			Label: "log2(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"log10": {
			Label: "log10(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"minute": {
			Label: "minute(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"month": {
			Label: "month(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"predict_linear": {
			Label: "predict_linear(v range-vector, t scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
				{Label: "t scalar"},
			},
		},
		"rate": {
			Label: "rate(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"resets": {
			Label: "resets(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"round": {
			Label: "round(v instant-vector, to_nearest=1 scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
				{Label: " to_nearest=1 scalar"},
			},
		},
		"scalar": {
			Label: "scalar(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"sort": {
			Label: "sort(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"sort_desc": {
			Label: "sort_desc(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"time": {
			Label:      "time()",
			Parameters: []protocol.ParameterInformation{},
		},
		"timestamp": {
			Label: "timestamp(v instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramInstantVector},
			},
		},
		"vector": {
			Label: "vector(s scalar)",
			Parameters: []protocol.ParameterInformation{
				{Label: "s scalar"},
			},
		},
		"year": {
			Label: "year(v=vector(time()) instant-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramVectorTimeInstant},
			},
		},
		"avg_over_time": {
			Label: "avg_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"sum_over_time": {
			Label: "sum_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"min_over_time": {
			Label: "min_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"max_over_time": {
			Label: "max_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"count_over_time": {
			Label: "count_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"stddev_over_time": {
			Label: "stddev_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"stdvar_over_time": {
			Label: "stdvar_over_time(v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: paramRangeVector},
			},
		},
		"qunatile_over_time": {
			Label: "quantile_over_time(s scalar, v range-vector)",
			Parameters: []protocol.ParameterInformation{
				{Label: "s scalar"},
				{Label: paramRangeVector},
			},
		},
	}

	ret, ok := signatures[name]
	if !ok {
		return protocol.SignatureInformation{}, errors.New("no signature found")
	}

	return ret, nil
}
