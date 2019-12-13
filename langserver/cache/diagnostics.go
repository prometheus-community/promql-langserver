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

package cache

import (
	"go/token"

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

func (d *DocumentHandle) promQLErrToProtocolDiagnostic(promQLErr *promql.ParseErr) (*protocol.Diagnostic, error) { // nolint: lll
	var pos protocol.Position

	var err error

	if pos, err = d.PositionToProtocolPosition(promQLErr.Position); err != nil {
		return nil, err
	}

	message := &protocol.Diagnostic{
		Range: protocol.Range{
			Start: pos,
			End:   EndOfLine(pos),
		},
		Severity: 1, // Error
		Source:   "promql-lsp",
		Message:  promQLErr.Err.Error(),
	}

	return message, nil
}

func (d *DocumentHandle) warnQuotedYaml(start token.Pos, end token.Pos) error {
	var startPosition token.Position

	var endPosition token.Position

	var err error

	startPosition, err = d.TokenPosToTokenPosition(start)
	if err != nil {
		return err
	}

	endPosition, err = d.TokenPosToTokenPosition(end)
	if err != nil {
		return err
	}

	message := &protocol.Diagnostic{
		Severity: 2, // Warning
		Source:   "promql-lsp",
		Message:  "Quoted queries are not supported by the language server",
	}

	if message.Range.Start, err = d.PositionToProtocolPosition(startPosition); err != nil {
		return err
	}

	if message.Range.End, err = d.PositionToProtocolPosition(endPosition); err != nil {
		return err
	}

	return d.AddDiagnostic(message)
}

// AddDiagnostic updates the compilation Results of a Document. Discards the Result if the context is expired
func (d *DocumentHandle) AddDiagnostic(diagnostic *protocol.Diagnostic) error {
	d.doc.mu.Lock()
	defer d.doc.mu.Unlock()

	select {
	case <-d.ctx.Done():
		return d.ctx.Err()
	default:
		d.doc.diagnostics = append(d.doc.diagnostics, *diagnostic)
		return nil
	}
}
