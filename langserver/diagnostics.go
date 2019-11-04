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

	"github.com/slrtbtfs/promql-lsp/langserver/cache"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// nolint:funlen
func (s *Server) diagnostics(uri string) {
	d, ctx, err := s.cache.GetDocument(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Document %v doesn't exist any more", uri)
	}

	version, expired := d.GetVersion(ctx)
	if expired != nil {
		return
	}

	// Everything is fine
	diagnostics := &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Version:     version,
		Diagnostics: []protocol.Diagnostic{},
	}

	compileResult, err := d.GetCompileResult(ctx)
	if err != nil {
		return
	}

	fmt.Fprint(os.Stderr, "Got it")

	if compileResult.Err != nil {
		var pos protocol.Position

		if pos, err = d.PositionToProtocolPostion(ctx, compileResult.Err.Position); err != nil {
			fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
			return
		}

		message := protocol.Diagnostic{
			Range: protocol.Range{
				Start: pos,
				End:   cache.EndOfLine(pos),
			},
			Severity: 1, // Error
			Source:   "promql-lsp",
			Message:  compileResult.Err.Err.Error(),
			Code:     "promql-parseerr",
			//Tags:    []protocol.DiagnosticTag{},
		}
		diagnostics.Diagnostics = append(diagnostics.Diagnostics, message)
	}

	if err = s.client.PublishDiagnostics(ctx, diagnostics); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to publish diagnostics")
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
