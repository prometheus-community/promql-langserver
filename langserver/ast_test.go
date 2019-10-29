package langserver

import (
	"fmt"
	"go/token"
	"reflect"
	"testing"

	"github.com/slrtbtfs/prometheus/promql"
)

func TestSmallestSurroundingNode(t *testing.T) {
	// TODO: #32
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
		node := getSmallestSurroundingNode(parseResult, test.pos)

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
			node := getSmallestSurroundingNode(parseResult, token.Pos(pos))
			// If we are outside the outermost Expression, nothing should be matched
			if node == nil && (int(parseResult.Pos()) > pos || int(parseResult.EndPos()) >= pos) {
				continue
			}
			if int(node.Pos()) > pos || int(node.EndPos()) <= pos {
				panic("The smallestSurroundingNode is not actually surrounding for input " + test +
					" and pos " + fmt.Sprintln(pos) + "Got: " + fmt.Sprintln(node) +
					"Pos: " + fmt.Sprintln(node.Pos()) + "EndPos: " + fmt.Sprintln(node.EndPos()))
			}

		}
	}
}
