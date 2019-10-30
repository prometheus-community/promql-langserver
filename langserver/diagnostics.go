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
	"os"

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// nolint:funlen
func (s *Server) diagnostics(ctx context.Context, d *document) {
	d.Mu.RLock()
	uri := d.doc.URI
	file := d.posData
	content := d.doc.Text
	version := d.doc.Version
	d.Mu.RUnlock()

	var diagnostics *protocol.PublishDiagnosticsParams

	switch d.doc.LanguageID {
	case "promql":
		ast, err := promql.ParseFile(content, file)

		var parseErr *promql.ParseErr

		var ok bool

		if err != nil {
			parseErr, ok = err.(*promql.ParseErr)

			fmt.Fprintf(os.Stderr, "Failed to convert %v to a promql.Parserr", err)

			if !ok {
				return
			}
		}

		recent := d.updateCompileData(version, ast, parseErr)
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

			if pos, ok = d.positionToProtocolPostion(version, parseErr.Position); !ok {
				fmt.Fprintf(os.Stderr, "Conversion failed\n")
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

		if err = s.client.PublishDiagnostics(ctx, diagnostics); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to publish diagnostics")
			fmt.Fprintln(os.Stderr, err.Error())
		}
	default:
		d.updateCompileData(version, nil, nil)
	}
}

// Updates the compilation Results of a document. Returns true if the Results were still recent
func (d *document) updateCompileData(version float64, ast promql.Node, err *promql.ParseErr) bool {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	defer d.compilers.Done()

	if d.doc.Version > version {
		return false
	}

	d.compileResult = compileResult{ast, err}

	return true
}
