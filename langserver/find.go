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
	"go/token"

	"github.com/prometheus-community/promql-langserver/langserver/cache"
	"github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/protocol"
	"github.com/prometheus/prometheus/promql"
)

type location struct {
	doc   *cache.DocumentHandle
	pos   token.Pos
	query *cache.CompiledQuery
	node  promql.Node
}

func (s *server) find(where *protocol.TextDocumentPositionParams) (there *location, err error) {
	there = &location{}

	if there.doc, err = s.cache.GetDocument(where.TextDocument.URI); err != nil {
		return
	}

	if there.pos, err = there.doc.ProtocolPositionToTokenPos(where.Position); err != nil {
		return
	}

	if there.query, err = there.doc.GetQuery(there.pos); err != nil {
		return
	}

	there.node = getSmallestSurroundingNode(there.query, there.pos)

	return
}
