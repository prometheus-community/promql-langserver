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

// This File includes code from the go/tools project which is governed by the following license:
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package langserver

import (
	"bytes"
	"context"

	"github.com/slrtbtfs/go-tools-vendored/jsonrpc2"
	"github.com/slrtbtfs/go-tools-vendored/lsp/protocol"
	"github.com/slrtbtfs/go-tools-vendored/span"
)

func (s *Server) DidOpen(ctx context.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc, err := s.cache.addDocument(&params.TextDocument)
	if err != nil {
		return err
	}

	go s.diagnostics(context.Background(), doc)
	return nil
}

func (s *Server) DidClose(_ context.Context, params *protocol.DidCloseTextDocumentParams) error {
	return s.cache.removeDocument(params.TextDocument.URI)
}

func (s *Server) DidChange(ctx context.Context, params *protocol.DidChangeTextDocumentParams) error {
	//options := s.session.Options()
	if len(params.ContentChanges) < 1 {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "no content changes provided")
	}

	uri := params.TextDocument.URI

	doc, err := s.cache.getDocument(uri)
	if err != nil {
		return err
	}

	// Check if the client sent the full content of the file.
	// We accept a full content change even if the server expected incremental changes.
	text, isFullChange := fullChange(params.ContentChanges)

	if !isFullChange {
		// Determine the new file content.
		text, err = doc.applyIncrementalChanges(params.ContentChanges, params.TextDocument.Version)
		if err != nil {
			return err
		}
	}

	// Cache the new file content
	err = doc.setContent(text, params.TextDocument.Version)
	if err != nil {
		return err
	}

	go s.diagnostics(context.Background(), doc)
	return nil
}
func fullChange(changes []protocol.TextDocumentContentChangeEvent) (string, bool) {
	if len(changes) > 1 {
		return "", false
	}
	// The length of the changes must be 1 at this point.
	if changes[0].Range == nil && changes[0].RangeLength == 0 {
		return changes[0].Text, true
	}
	return "", false
}
func (d *document) applyIncrementalChanges(changes []protocol.TextDocumentContentChangeEvent, version float64) (string, error) {
	d.Mu.RLock()
	if version <= d.doc.Version {
		return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}
	content := []byte(d.doc.Text)
	uri := d.doc.URI

	d.Mu.RUnlock()
	for _, change := range changes {
		// Update column mapper along with the content.
		converter := span.NewContentConverter(uri, content)
		m := &protocol.ColumnMapper{
			URI:       span.URI(d.doc.URI),
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
		//fmt.Fprintf(os.Stderr, string(content))
	}
	return string(content), nil
}
