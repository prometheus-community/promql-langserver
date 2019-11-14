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
	"testing"

	"github.com/slrtbtfs/prometheus/promql"
)

// Call the (* Document) Functions with an expired context. Expected behaviour is that all
// of these calls return an error
func TestDocumentContext(t *testing.T) { //nolint: funlen
	d := &Document{}

	expired, cancel := context.WithCancel(context.Background())

	cancel()

	// From diagnostics.go

	if _, err := d.promQLErrToProtocolDiagnostic(expired, &promql.ParseErr{}); err == nil {
		panic("Expected promQLErrToProtocolDiagnostic to fail with expired context")
	}

	if err := d.warnQuotedYaml(expired, token.NoPos, token.NoPos); err == nil {
		panic("Expected warnQuotedYaml to fail with expired context")
	}

	if err := d.AddDiagnostic(expired, nil); err == nil {
		panic("Expected AddDiagnostic to fail with expired context")
	}

	// From document.go

	if _, err := d.GetContent(expired); err == nil {
		panic("Expected GetContent to fail with expired context")
	}

	if _, err := d.GetSubstring(expired, token.NoPos, token.NoPos); err == nil {
		panic("Expected GetSubstring to fail with expired context")
	}

	if _, err := d.GetQueries(expired); err == nil {
		panic("Expected GetQueries to fail with expired context")
	}

	if _, err := d.GetQuery(expired, token.NoPos); err == nil {
		panic("Expected GetQuery to fail with expired context")
	}

	if _, err := d.GetVersion(expired); err == nil {
		panic("Expected GetVersion to fail with expired context")
	}

	if _, err := d.GetYamls(expired); err == nil {
		panic("Expected GetYamls to fail with expired context")
	}

	if _, err := d.GetDiagnostics(expired); err == nil {
		panic("Expected GetDiagnostics to fail with expired context")
	}

	// From position.go

	if _, err := d.PositionToProtocolPosition(expired, token.Position{}); err == nil {
		panic("Expected PositionToProtocolPosition to fail with expired context")
	}

	if _, err := d.PosToProtocolPosition(expired, token.NoPos); err == nil {
		panic("Expected PosToProtocolPosition to fail with expired context")
	}

	if _, err := d.yamlPositionToTokenPos(expired, 0, 0, 0); err == nil {
		panic("Expected yamlPositionToTokenPos to fail with expired context")
	}

	if _, err := d.TokenPosToTokenPosition(expired, token.NoPos); err == nil {
		panic("Expected TokenPosToTokenPosition to fail with expired context")
	}

	if _, err := d.GetVersion(expired); err == nil {
		panic("Expected GetVersion to fail with expired context")
	}

	if _, err := d.GetYamls(expired); err == nil {
		panic("Expected GetYamls to fail with expired context")
	}

	if _, err := d.GetDiagnostics(expired); err == nil {
		panic("Expected GetContent to fail with expired context")
	}
}
