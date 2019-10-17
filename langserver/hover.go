package langserver

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/prometheus/promql"
)

func (s *Server) Hover(_ context.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	doc, err := s.cache.getDocument(params.TextDocumentPositionParams.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// FIXME: This is still a bit racy
	doc.compilers.Wait()
	doc.Mu.RLock()
	defer doc.Mu.RUnlock()
	pos, err := doc.protocolPositionToTokenPos(params.TextDocumentPositionParams.Position)
	if err != nil {
		return nil, err
	}
	node := getSmallestSourroundingNode(doc.compileResult.ast, pos)

	markdown := nodeToDocMarkdown(node)

	if markdown != "" {
		return &protocol.Hover{
			Contents: protocol.MarkupContent{
				Kind:  "markdown",
				Value: markdown,
			},
		}, nil
	} else {
		return nil, nil
	}
}

func nodeToDocMarkdown(node promql.Node) string {
	var ret bytes.Buffer
	expr, ok := node.(promql.Expr)
	if ok {
		_, err := ret.WriteString(fmt.Sprintf("Type: %v\n", expr.Type()))
		if err != nil {
			return ""
		}
	}

	var call *promql.Call
	call, ok = node.(*promql.Call)
	if ok {
		doc := funcDocStrings(strings.ToLower(call.Func.Name))
		_, err := ret.WriteString(doc)
		if err != nil {
			return ""
		}
	}
	return string(ret.Bytes())
}

func funcDocStrings(name string) string {
	// TODO
	/*
		matcher = regexp.MustCompile("^## `"+name+"(?m).*?^##")
		file :=
		documentation = string(ioutil.Read)
	*/
	//stub
	if name == "label_replace" {
		return "## `label_replace()`\n\nFor each timeseries in `v`, `label_replace(v instant-vector, dst_label string,\nreplacement string, src_label string, regex string)` matches the regular\nexpression `regex` against the label `src_label`.  If it matches, then the\ntimeseries is returned with the label `dst_label` replaced by the expansion of\n`replacement`. `$1` is replaced with the first matching subgroup, `$2` with the\nsecond etc. If the regular expression doesn't match then the timeseries is\nreturned unchanged.\n\nThis example will return a vector with each time series having a `foo`\nlabel with the value `a` added to it:\n\n```\nlabel_replace(up{job=\"api-server\",service=\"a:c\"}, \"foo\", \"$1\", \"service\", \"(.*):.*\")\n```\n"
	} else {
		return ""
	}

}
