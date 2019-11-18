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
	"errors"
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

// DocumentCache caches the documents and compile Results associated with one server-client connection
type DocumentCache struct {
	fileSet *token.FileSet

	documents map[protocol.DocumentURI]*Document
	mu        sync.RWMutex
}

// Init Initializes a Document cache
func (c *DocumentCache) Init() {
	c.fileSet = token.NewFileSet()
	c.mu.Lock()
	c.documents = make(map[protocol.DocumentURI]*Document)
	c.mu.Unlock()
}

// AddDocument adds a Document to the cache
func (c *DocumentCache) AddDocument(doc *protocol.TextDocumentItem) (*Document, error) {
	if _, ok := c.documents[doc.URI]; ok {
		return nil, errors.New("document already exists")
	}

	file := c.fileSet.AddFile(doc.URI, -1, maxDocumentSize)

	if r := recover(); r != nil {
		if err, ok := r.(error); !ok {
			return nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: %v", err)
		}
	}

	file.SetLinesForContent([]byte(doc.Text))

	d := &Document{
		posData:    file,
		uri:        doc.URI,
		languageID: doc.LanguageID,
	}

	err := d.SetContent(doc.Text, doc.Version, true)
	if err != nil {
		return nil, err
	}

	c.mu.Lock()
	c.documents[doc.URI] = d
	c.mu.Unlock()

	return d, nil
}

// GetDocument retrieve a Document from the cache
// Additionally returns a context that expires as soon as the document changes
func (c *DocumentCache) GetDocument(uri protocol.DocumentUri) (*Document, context.Context, error) {
	c.mu.RLock()
	ret, ok := c.documents[uri]
	c.mu.RUnlock()

	if !ok {
		return nil, nil, jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/getDocument: Document not found: %v", uri)
	}

	return ret, ret.versionCtx, nil
}

// RemoveDocument removes a Document from the cache
func (c *DocumentCache) RemoveDocument(uri protocol.DocumentURI) error {
	c.mu.Lock()
	delete(c.documents, uri)
	c.mu.Unlock()

	return nil
}
