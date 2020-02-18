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
	"fmt"
	"testing"

	"github.com/prometheus-community/promql-langserver/vendored/go-tools/lsp/protocol"
)

func TestCache(t *testing.T) { // nolint:funlen
	c := &DocumentCache{}

	c.Init()

	doc, err := c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{
			URI:        "test_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       "test_text",
		})
	if err != nil {
		panic("Failed to AddDocument() to cache")
	}

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "test_file",
			LanguageID: "yaml",
			Version:    1,
			Text:       "test_text",
		})
	if err == nil {
		panic("Should not be able to add same document twice")
	}

	doc1, err := c.GetDocument("test_file")
	if err != nil {
		panic("Failed to GetDocument() from cache")
	}

	if doc1.doc != doc.doc {
		panic("Cache returned wrong document")
	}

	err = c.RemoveDocument("test_file")
	if err != nil {
		panic("Failed to RemoveDocument() from cache")
	}

	err = c.RemoveDocument("test_file")
	if err == nil {
		panic("should have failed to RemoveDocument() twice")
	}

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "test_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       "test_text",
		})
	if err != nil {
		panic("Should be able to readd document after removing it")
	}

	tooLongString := string(make([]byte, maxDocumentSize+1))

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "long_test_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       tooLongString,
		})
	if err == nil {
		panic("Shouldn't be able to add overlong document")
	}

	wrongYaml := "asdf["

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "wrong_yaml_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       wrongYaml,
		})
	if err != nil {
		panic("Should be able to handle yaml with syntax errors")
	}

	rulesFile := `
groups:
  - name: example
    rules:
    - record: job:http_inprogress_requests:sum
      expr: sum(http_inprogress_requests) by (job)
    - record: job:http_inprogress_requests:sum:wrong
      expr: 
          sum(http_inprogress_requests) by (job
    - record: job:http_inprogress_requests:sum:quoted
      expr: "sum(http_inprogress_requests) by (job)"
`

	_, err = c.AddDocument(
		context.Background(),
		&protocol.TextDocumentItem{

			URI:        "rules_file",
			LanguageID: "yaml",
			Version:    0,
			Text:       rulesFile,
		})
	if err != nil {
		panic("adding rules file failed")
	}

	doc, err = c.GetDocument("rules_file")
	if err != nil {
		panic("failed to get rules file")
	}

	diagnostics, err := doc.GetDiagnostics()
	if err != nil {
		panic("failed to get diagnostics for rules file")
	}

	if len(diagnostics) != 3 {
		fmt.Println(diagnostics)
		panic("expected exactly 3 error messages for rules file got " + fmt.Sprint(len(diagnostics)))
	}
}
