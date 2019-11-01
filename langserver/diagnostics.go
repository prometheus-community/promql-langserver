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
	"os"

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// nolint:funlen
func (s *Server) diagnostics(uri string) {
	d, ctx, err := s.cache.GetDocument(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Document %v doesn't exist any more", uri)
	}

	file := d.PosData

	content, expired := d.GetContent(ctx)
	if expired != nil {
		return
	}

	var version float64

	version, expired = d.GetVersion(ctx)
	if expired != nil {
		return
	}

	var diagnostics *protocol.PublishDiagnosticsParams

	switch d.GetLanguageID() {
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

		recent := d.UpdateCompileData(ctx, ast, parseErr)
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

			if pos, ok = d.PositionToProtocolPostion(version, parseErr.Position); !ok {
				fmt.Fprintf(os.Stderr, "Conversion failed\n")
				return
			}

			message := protocol.Diagnostic{
				Range: protocol.Range{
					Start: pos,
					End:   cache.EndOfLine(pos),
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
		d.UpdateCompileData(ctx, nil, nil)
	}
}
