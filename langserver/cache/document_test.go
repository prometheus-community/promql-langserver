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

	"github.com/prometheus/prometheus/promql"
	"gopkg.in/yaml.v3"
)

// Call the (* Document) Functions with an expired context. Expected behaviour is that all
// of these calls return an error
func TestDocumentContext(t *testing.T) { //nolint: funlen
	doc := &document{}

	doc.posData = token.NewFileSet().AddFile("", -1, 0)

	doc.compilers.initialize()

	expired, cancel := context.WithCancel(context.Background())

	cancel()

	d := &DocumentHandle{doc, expired}

	// From compile.go

	// Necessary since compile() will call d.compilers.Done()
	d.doc.compilers.Add(1)

	d.doc.languageID = "promql"

	if err := d.compile(); err == nil {
		panic("Expected compile to fail with expired context (languageID: promql)")
	}

	// Necessary since compile() will call d.compilers.Done()
	d.doc.compilers.Add(1)

	d.doc.languageID = "yaml"

	if err := d.compile(); err == nil {
		panic("Expecexpiredted compile to fail with expired context (languageID: promql)")
	}

	// Necessary since compileQuery() will call d.compilers.Done()
	d.doc.compilers.Add(1)

	if err := d.compileQuery(true, token.NoPos, token.NoPos, ""); err == nil {
		panic("Expected compileQuery to fail with expired context (fullFile: true)")
	}

	// Necessary since compileQuery() will call d.compilers.Done()
	d.doc.compilers.Add(1)

	if err := d.compileQuery(false, token.NoPos, token.NoPos, ""); err == nil {
		panic("Expected compileQuery to fail with expired context (fullFile: false)")
	}

	if err := d.addCompileResult(token.NoPos, &promql.MatrixSelector{}, nil, "", ""); err == nil {
		panic("Expected AddCompileResult to fail with expired context")
	}

	// From diagnostics.go

	if _, err := d.promQLErrToProtocolDiagnostic(token.NoPos, &promql.ParseErr{}); err == nil {
		panic("Expected promQLErrToProtocolDiagnostic to fail with expired context")
	}

	if err := d.warnQuotedYaml(token.NoPos, token.NoPos); err == nil {
		panic("Expected warnQuotedYaml to fail with expired context")
	}

	if err := d.addDiagnostic(nil); err == nil {
		panic("Expected AddDiagnostic to fail with expired context")
	}

	// From document.go

	if _, err := d.GetContent(); err == nil {
		panic("Expected GetContent to fail with expired context")
	}

	if _, err := d.GetSubstring(token.NoPos, token.NoPos); err == nil {
		panic("Expected GetSubstring to fail with expired context")
	}

	if _, err := d.GetQueries(); err == nil {
		panic("Expected GetQueries to fail with expired context")
	}

	if _, err := d.getQuery(token.NoPos); err == nil {
		panic("Expected GetQuery to fail with expired context")
	}

	if _, err := d.GetVersion(); err == nil {
		panic("Expected GetVersion to fail with expired context")
	}

	if _, err := d.getYamlDocuments(); err == nil {
		panic("Expected GetYamls to fail with expired context")
	}

	if _, err := d.GetDiagnostics(); err == nil {
		panic("Expected GetDiagnostics to fail with expired context")
	}

	// From position.go

	if _, err := d.tokenPositionToProtocolPosition(token.Position{}); err == nil {
		panic("Expected PositionToProtocolPosition to fail with expired context")
	}

	if _, err := d.PosToProtocolPosition(token.NoPos); err == nil {
		panic("Expected PosToProtocolPosition to fail with expired context")
	}

	if _, err := d.yamlPositionToTokenPos(0, 0, 0); err == nil {
		panic("Expected YamlPositionToTokenPos to fail with expired context")
	}

	if _, err := d.tokenPosToTokenPosition(token.NoPos); err == nil {
		panic("Expected TokenPosToTokenPosition to fail with expired context")
	}

	if _, err := d.GetVersion(); err == nil {
		panic("Expected GetVersion to fail with expired context")
	}

	if _, err := d.getYamlDocuments(); err == nil {
		panic("Expected GetYamls to fail with expired context")
	}

	if _, err := d.GetDiagnostics(); err == nil {
		panic("Expected GetContent to fail with expired context")
	}

	// From yaml.go

	if err := d.parseYamls(); err == nil {
		panic("Expected ParseYamls to fail with expired context")
	}

	if err := d.addYaml(nil); err == nil {
		panic("Expected addYaml to fail with expired context")
	}

	// Necessary since scanYamlTree will call d.compilers.Done()
	d.doc.compilers.Add(1)

	if err := d.scanYamlTree(); err == nil {
		panic("Expected scanYamlTree to fail with expired context")
	}

	/*
		Excluded since it does not do anything on an empty document
		if err := d.scanYamlTreeRec(expired, &yaml.Node{}, token.NoPos, 0); err == nil {
			panic("Expected scanYamlTreeRec to fail with expired context")
		}
	*/
	if err := d.foundQuery(&yaml.Node{}, token.NoPos, nil, 0); err == nil {
		panic("Expected foundQuery to fail with expired context")
	}
}
