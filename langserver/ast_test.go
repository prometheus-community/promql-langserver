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
	"fmt"
	"go/token"
	"reflect"
	"testing"

	"github.com/prometheus-community/promql-langserver/langserver/cache"
	"github.com/prometheus/prometheus/promql"
)

func TestSmallestSurroundingNode(t *testing.T) {
	shouldMatchFull := []struct {
		input string
		pos   token.Pos
	}{
		{
			input: "1",
			pos:   1,
		}, {
			input: "+1 + -2 * 1",
			pos:   4,
		},
	}
	for _, test := range shouldMatchFull {
		parseResult, err := promql.ParseExpr(test.input)
		if err != nil {
			panic("Parser should not have failed on " + test.input)
		}

		node := getSmallestSurroundingNode(&cache.CompiledQuery{Ast: parseResult}, test.pos)

		if !reflect.DeepEqual(node, parseResult) {
			panic("Whole Expression should have been matched for " + test.input)
		}
	}

	for _, test := range testExpressions {
		parseResult, err := promql.ParseExpr(test)
		if err != nil {
			// We're currently only interested in correct expressions
			continue
		}

		for pos := 1; pos <= len(test); pos++ {
			node := getSmallestSurroundingNode(&cache.CompiledQuery{Ast: parseResult}, token.Pos(pos))
			// If we are outside the outermost Expression, nothing should be matched
			if node == nil && (int(parseResult.PositionRange().Start) > pos || int(parseResult.PositionRange().End) >= pos) {
				continue
			}

			if int(node.PositionRange().Start) > pos || int(node.PositionRange().End) < pos {
				panic("The smallestSurroundingNode is not actually surrounding for input " + test +
					" and pos " + fmt.Sprintln(pos) + "Got: " + fmt.Sprintln(node) +
					"Pos: " + fmt.Sprintln(node.PositionRange().Start) + "EndPos: " + fmt.Sprintln(node.PositionRange().End))
			}
		}
	}
}
