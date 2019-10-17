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
	"go/token"

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/go-tools-vendored/span"
	"github.com/slrtbtfs/prometheus/promql"
)

func (s *Server) diagnostics(ctx context.Context, doc *document) {
	doc.Mu.RLock()
	uri := doc.doc.URI
	file := doc.posData
	content := doc.doc.Text
	version := doc.doc.Version
	doc.Mu.RUnlock()
	var diagnostics *protocol.PublishDiagnosticsParams
	switch doc.doc.LanguageID {
	case "promql":
		ast, err := promql.ParseFile(content, file)

		var parseErr *promql.ParseErr = nil
		var ok bool
		if err != nil {
			parseErr, ok = err.(*promql.ParseErr)
			// TODO (slrtbtfs) Maybe give some more feedback here
			if !ok {
				return
			}
		}

		recent := doc.updateCompileData(version, ast, parseErr)
		if !recent {
			return
		}

		// Everything is fine
		diagnostics = &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Version:     version,
			Diagnostics: []protocol.Diagnostic{},
		}
		if err != nil {
			var pos protocol.Position
			pos, ok = doc.positionToProtocolPostion(version, parseErr.Position)
			if !ok {
				return
			}

			message := protocol.Diagnostic{
				Range: protocol.Range{
					Start: pos,
					End:   endOfLine(pos),
				},
				Severity: 1, // Error
				Source:   "promql-lsp",
				Message:  parseErr.Err.Error(),
				Code:     "promql-parseerr",
				//Tags:    []protocol.DiagnosticTag{},
			}
			diagnostics.Diagnostics = append(diagnostics.Diagnostics, message)
		}
		s.client.PublishDiagnostics(ctx, diagnostics)
	default:
		doc.updateCompileData(version, nil, nil)
	}

}

// Updates the compilation Results of a document. Returns true if the Results were still recent
func (doc *document) updateCompileData(version float64, ast promql.Node, err *promql.ParseErr) bool {
	doc.Mu.Lock()
	defer doc.Mu.Unlock()
	defer doc.compilers.Done()
	if doc.doc.Version > version {
		return false
	}
	doc.compileResult = compileResult{ast, err}
	return true
}

func (doc *document) positionToProtocolPostion(version float64, pos token.Position) (protocol.Position, bool) {
	doc.Mu.RLock()
	defer doc.Mu.RUnlock()
	if doc.doc.Version > version {
		return protocol.Position{}, false
	}
	line := pos.Line
	char := pos.Column
	// Convert to the Postions as described in the LSP Spec
	point := span.NewPoint(line, char, int(doc.posData.LineStart(line)))
	var err error
	char, err = span.ToUTF16Column(point, []byte(doc.doc.Text))
	char = char - 1
	if err != nil {
		return protocol.Position{}, false
	}
	line = line - 1
	return protocol.Position{
		Line:      float64(line),
		Character: float64(char),
	}, true
}

func endOfLine(p protocol.Position) protocol.Position {
	return protocol.Position{
		Line:      p.Line + 1,
		Character: 0,
	}
}
