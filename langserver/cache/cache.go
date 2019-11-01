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
	"bytes"
	"context"
	"go/token"
	"sync"

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/span"
)

// We need this so we can reserve a certain position range in the FileSet
// for each Document.
// Anything that is larger than 1MB would probably not work with reasonable performance anyway
// The bad thing is, that it adds an 2000 file limit (no of files per connection)
// on 32bit systems
const maxDocumentSize = 1000000

type DocumentCache struct {
	FileSet *token.FileSet

	Documents   map[protocol.DocumentURI]*Document
	DocumentsMu sync.RWMutex
}

type Document struct {
	PosData *token.File
	Doc     *protocol.TextDocumentItem
	Mu      sync.RWMutex

	versionCtx      context.Context
	obsoleteVersion context.CancelFunc

	CompileResult CompileResult

	// Wait for this before accessing  compileResults
	Compilers sync.WaitGroup
}

type CompileResult struct {
	Ast promql.Node
	Err *promql.ParseErr
}

// Initializes a Document cache
func (c *DocumentCache) Init() {
	c.FileSet = token.NewFileSet()
	c.DocumentsMu.Lock()
	c.Documents = make(map[protocol.DocumentURI]*Document)
	c.DocumentsMu.Unlock()
}

// Add a Document to the cache
func (c *DocumentCache) AddDocument(doc *protocol.TextDocumentItem) (*Document, error) {
	if len(doc.Text) > maxDocumentSize {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: Provided Document to large.")
	}

	file := c.FileSet.AddFile(doc.URI, -1, maxDocumentSize)

	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: %v", err)
		}
	}

	file.SetLinesForContent([]byte(doc.Text))

	d := &Document{
		PosData: file,
		Doc:     doc,
	}

	d.versionCtx, d.obsoleteVersion = context.WithCancel(context.Background())

	c.DocumentsMu.Lock()
	c.Documents[doc.URI] = d
	c.DocumentsMu.Unlock()

	return d, nil
}

// retrieve a Document from the cache
func (c *DocumentCache) GetDocument(uri protocol.DocumentUri) (*Document, error) {
	c.DocumentsMu.RLock()
	ret, ok := c.Documents[uri]
	c.DocumentsMu.RUnlock()

	if !ok {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/getDocument: Document not found: %v", uri)
	}

	return ret, nil
}

// Remove a Document from the cache
func (c *DocumentCache) RemoveDocument(uri protocol.DocumentURI) error {
	c.DocumentsMu.Lock()
	delete(c.Documents, uri)
	c.DocumentsMu.Unlock()

	return nil
}

// Set the content after an update send by the client. Must increase the version number
func (d *Document) SetContent(content string, version float64) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	if version <= d.Doc.Version {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	if len(content) > maxDocumentSize {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: Provided Document to large.")
	}

	d.obsoleteVersion()

	d.versionCtx, d.obsoleteVersion = context.WithCancel(context.Background())

	d.Doc.Text = content
	d.Doc.Version = version
	d.PosData.SetLinesForContent([]byte(content))

	return nil
}

func (d *Document) ApplyIncrementalChanges(changes []protocol.TextDocumentContentChangeEvent, version float64) (string, error) { //nolint:lll
	d.Mu.RLock()

	if version <= d.Doc.Version {
		return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	content := []byte(d.Doc.Text)
	uri := d.Doc.URI

	d.Mu.RUnlock()

	for _, change := range changes {
		// Update column mapper along with the content.
		converter := span.NewContentConverter(uri, content)
		m := &protocol.ColumnMapper{
			URI:       span.URI(d.Doc.URI),
			Converter: converter,
			Content:   content,
		}

		spn, err := m.RangeSpan(*change.Range)

		if err != nil {
			return "", err
		}

		if !spn.HasOffset() {
			return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "invalid range for content change")
		}

		start, end := spn.Start().Offset(), spn.End().Offset()
		if end < start {
			return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "invalid range for content change")
		}

		var buf bytes.Buffer

		buf.Write(content[:start])
		buf.WriteString(change.Text)
		buf.Write(content[end:])

		content = buf.Bytes()
	}

	return string(content), nil
}
