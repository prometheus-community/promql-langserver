package langserver

import (
	"go/token"
	"sync"

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
)

// We need this so we can reserve a certain position range in the FileSet
// for each document.
// Anything that is larger than 1MB probably is an attempt to bring down the Server anyways
// The bad thing is, that it adds an 2000 file limit (no of open files per connection)
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
}

// Initializes a document cache
func (c *documentCache) init() {
	c.fileSet = token.NewFileSet()
	c.documentsMu.Lock()
	c.documents = make(map[protocol.DocumentURI]*document)
	c.documentsMu.Unlock()
}

// Add a document to the cache
func (c *documentCache) addDocument(doc *protocol.TextDocumentItem) error {
	if len(doc.Text) > maxDocumentSize {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: Provided document to large.")
	}

	file := c.fileSet.AddFile(doc.URI, -1, maxDocumentSize)
	if r := recover(); r != nil {
		var ok bool
		_, ok = r.(error)
		if !ok {
			return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/addDocument: %v", r)
		}
	}

	file.SetLinesForContent([]byte(doc.Text))

	c.documentsMu.Lock()
	c.documents[doc.URI] = &document{
		posData: file,
		doc:     doc,
	}
	c.documentsMu.Unlock()
	return nil
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
