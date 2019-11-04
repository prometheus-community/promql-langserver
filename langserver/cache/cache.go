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
	"sync"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
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

// Initializes a Document cache
func (c *DocumentCache) Init() {
	c.FileSet = token.NewFileSet()
	c.DocumentsMu.Lock()
	c.Documents = make(map[protocol.DocumentURI]*Document)
	c.DocumentsMu.Unlock()
}

// Add a Document to the cache
func (c *DocumentCache) AddDocument(doc *protocol.TextDocumentItem) (*Document, error) {
	file := c.FileSet.AddFile(doc.URI, -1, maxDocumentSize)

	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: %v", err)
		}
	}

	file.SetLinesForContent([]byte(doc.Text))

	d := &Document{
		PosData:    file,
		uri:        doc.URI,
		languageID: doc.LanguageID,
	}

	err := d.SetContent(doc.Text, doc.Version, true)
	if err != nil {
		return nil, err
	}

	c.DocumentsMu.Lock()
	c.Documents[doc.URI] = d
	c.DocumentsMu.Unlock()

	return d, nil
}

// Retrieve a Document from the cache
// Additionally returns a context that expires as soon as the document changes
func (c *DocumentCache) GetDocument(uri protocol.DocumentUri) (*Document, context.Context, error) {
	c.DocumentsMu.RLock()
	ret, ok := c.Documents[uri]
	c.DocumentsMu.RUnlock()

	if !ok {
		return nil, nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/getDocument: Document not found: %v", uri)
	}

	return ret, ret.versionCtx, nil
}

// Remove a Document from the cache
func (c *DocumentCache) RemoveDocument(uri protocol.DocumentURI) error {
	c.DocumentsMu.Lock()
	delete(c.Documents, uri)
	c.DocumentsMu.Unlock()

	return nil
}
