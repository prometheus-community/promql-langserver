package langserver

import (
	"go/token"

	"github.com/slrtbtfs/prometheus/promql"
)

func getSmallestSourroundingNode(ast promql.Node, pos token.Pos) promql.Node {
	ret := ast
BIG_LOOP:
	for {
		for _, child := range ret.Childs() {
			if child.Pos() <= pos && child.EndPos() > pos {
				ret = child
				continue BIG_LOOP
			}
		}
		break
	}

	return ret
}
