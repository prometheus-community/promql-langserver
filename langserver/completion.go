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

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// Completion is required by the protocol.Server interface
// nolint: wsl
func (s *Server) Completion(ctx context.Context, params *protocol.CompletionParams) (*protocol.CompletionList, error) {
	fmt.Fprintln(os.Stderr, "0")
	doc, docCtx, err := s.cache.GetDocument(params.TextDocument.URI)
	if err != nil {
		fmt.Fprintln(os.Stderr, "1")
		return nil, err
	}

	var pos token.Pos

	pos, err = doc.ProtocolPositionToTokenPos(docCtx, params.TextDocumentPositionParams.Position)
	if err != nil {
		fmt.Fprintln(os.Stderr, "2")
		return nil, nil
	}

	var query *cache.CompiledQuery

	query, err = doc.GetQuery(docCtx, pos-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, "3")
		return nil, nil
	}

	node := getSmallestSurroundingNode(query.Ast, pos-1)
	if node == nil {
		fmt.Fprintln(os.Stderr, "3")
		return nil, nil
	}

	if pos != node.EndPos() {
		// Not at the end of the expression to be completed
		fmt.Fprintln(os.Stderr, "4")
		return nil, nil
	}

	return s.getCompletions(ctx, node)
}

func (s *Server) getCompletions(_ context.Context, node promql.Node) (*protocol.CompletionList, error) {
	fmt.Fprintf(os.Stderr, "getCompletions called: %+v\n\n", node)
	return nil, notImplemented("Completion")
}
