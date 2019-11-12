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
	"context"
	"go/token"

	"github.com/prometheus/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

func (d *Document) promQLErrToProtocolDiagnostic(ctx context.Context, promQLErr *promql.ParseErr) (*protocol.Diagnostic, error) { // nolint: lll
	var pos protocol.Position

	var err error

	if pos, err = d.PositionToProtocolPostion(ctx, promQLErr.Position); err != nil {
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

func (d *Document) warnQuotedYaml(ctx context.Context, start token.Pos, end token.Pos) error {
	// d.posData is syncronized itself. The results might be wrong if the document
	// changes in the meantime, but that does not matter since in this case the results
	// will be discarded later
	startPosition := d.posData.Position(start)
	endPosition := d.posData.Position(end)

	message := &protocol.Diagnostic{
		Severity: 2, // Warning
		Source:   "promql-lsp",
		Message:  "Quoted queries are not supported by the language server",
	}

	var err error

	if message.Range.Start, err = d.PositionToProtocolPostion(ctx, startPosition); err != nil {
		return err
	}

	if message.Range.End, err = d.PositionToProtocolPostion(ctx, endPosition); err != nil {
		return err
	}

	return d.AddDiagnostic(ctx, message)
}

// AddDiagnostic updates the compilation Results of a Document. Discards the Result if the context is expired
func (d *Document) AddDiagnostic(ctx context.Context, diagnostic *protocol.Diagnostic) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		d.diagnostics = append(d.diagnostics, *diagnostic)
		return nil
	}
}
