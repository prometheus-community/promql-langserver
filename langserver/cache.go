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
	"go/token"
	"sync"

	"github.com/slrtbtfs/prometheus/promql"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
)

// We need this so we can reserve a certain position range in the FileSet
// for each document.
// Anything that is larger than 1MB would probably not work with reasonable performance anyway
// The bad thing is, that it adds an 2000 file limit (no of files per connection)
// on 32bit systems
const maxDocumentSize = 1000000

type documentCache struct {
	fileSet *token.FileSet

	documents   map[protocol.DocumentURI]*document
	documentsMu sync.RWMutex
}

type document struct {
	posData *token.File
	doc     *protocol.TextDocumentItem
	Mu      sync.RWMutex

	compileResult compileResult

	// Wait for this before accessing  compileResults
	compilers sync.WaitGroup
}

type compileResult struct {
	ast promql.Node
	err *promql.ParseErr
}

// Initializes a document cache
func (c *documentCache) init() {
	c.fileSet = token.NewFileSet()
	c.documentsMu.Lock()
	c.documents = make(map[protocol.DocumentURI]*document)
	c.documentsMu.Unlock()
}

// Add a document to the cache
func (c *documentCache) addDocument(doc *protocol.TextDocumentItem) (*document, error) {
	if len(doc.Text) > maxDocumentSize {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: Provided document to large.")
	}

	file := c.fileSet.AddFile(doc.URI, -1, maxDocumentSize)
	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: %v", err)
		}
	}

	file.SetLinesForContent([]byte(doc.Text))

	docu := &document{
		posData: file,
		doc:     doc,
	}

	c.documentsMu.Lock()
	c.documents[doc.URI] = docu
	c.documentsMu.Unlock()

	return docu, nil
}

// retrieve a document from the cache
func (c *documentCache) getDocument(uri protocol.DocumentUri) (*document, error) {
	c.documentsMu.RLock()
	ret, ok := c.documents[uri]
	c.documentsMu.RUnlock()

	if !ok {
		return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/getDocument: Document not found: %v", uri)
	}

	return ret, nil
}

// Remove a document from the cache
func (c *documentCache) removeDocument(uri protocol.DocumentURI) error {
	c.documentsMu.Lock()
	delete(c.documents, uri)
	c.documentsMu.Unlock()

	return nil
}

// Set the content after an update send by the client. Must increase the version number
func (d *document) setContent(content string, version float64) error {
	d.Mu.Lock()
	defer d.Mu.Unlock()
	if version <= d.doc.Version {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}
	d.doc.Text = content
	d.doc.Version = version
	d.posData.SetLinesForContent([]byte(content))

	return nil
}
