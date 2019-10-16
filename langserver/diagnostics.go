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

	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/prometheus/promql"
)

func (s *Server) diagnostics(ctx context.Context, doc *document) {
	//	ctx, done := trace.StartSpan(ctx, "lsp:background-worker")
	//defer done()
	doc.Mu.RLock()
	uri := doc.doc.URI
	file := doc.posData
	content := doc.doc.Text
	version := doc.doc.Version
	doc.Mu.RUnlock()
	switch doc.doc.LanguageID {
	case "promql":
		_, err := promql.ParseFile(content, file)

		// Everything is fine
		if err == nil {
			return
		}

		s.client.ShowMessage(ctx, &protocol.ShowMessageParams{
			Type:    protocol.Error,
			Message: "Found Error",
		})
		parseErr, ok := err.(*promql.ParseErr)

		// TODO (slrtbtfs) Maybe give some more feedback here
		if !ok {
			return
		}
		line := parseErr.Position.Line - 1
		char := parseErr.Position.Column - 1

		message := protocol.Diagnostic{
			Range: protocol.Range{
				Start: protocol.Position{
					Line:      float64(line),
					Character: float64(char),
				},
				End: protocol.Position{
					Line:      float64(line) + 1,
					Character: 0,
				},
			},
			Severity: 1, // Error
			Source:   "promql-lsp",
			Message:  parseErr.Err.Error(),
			Code:     "promql-parseerr",
			//Tags:    []protocol.DiagnosticTag{},
		}
		diagnostics := &protocol.PublishDiagnosticsParams{
			URI:         uri,
			Version:     version,
			Diagnostics: []protocol.Diagnostic{message},
		}
		doc.Mu.RLock()
		newVersion := doc.doc.Version
		doc.Mu.RUnlock()
		// There is no point in publishing our diagnostics if they are already outdated
		if newVersion > version {
			return
		}
		s.client.PublishDiagnostics(ctx, diagnostics)

	}

}
