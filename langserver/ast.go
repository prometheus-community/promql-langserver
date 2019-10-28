package langserver

import (
	"go/token"

	"github.com/slrtbtfs/prometheus/promql"
)

func getSmallestSurroundingNode(ast promql.Node, pos token.Pos) promql.Node {
	if pos < ast.Pos() || pos >= ast.EndPos() {
		return nil
	}
	ret := ast
BIG_LOOP:
	for {
		for _, child := range ret.Children() {
			if child.Pos() <= pos && child.EndPos() > pos {
				ret = child
				continue BIG_LOOP
			}
		}
		break
	}

	return ret
}
