package cache

import (
	"bytes"
	"context"
	"go/token"
	"sync"

	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/jsonrpc2"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/lsp/protocol"
	"github.com/slrtbtfs/promql-lsp/vendored/go-tools/span"
)

type Document struct {
	PosData *token.File
	doc     *protocol.TextDocumentItem
	mu      sync.RWMutex

	versionCtx      context.Context
	obsoleteVersion context.CancelFunc

	compileResult *CompiledQuery

	// Wait for this before accessing  compileResults
	Compilers sync.WaitGroup
}

func (d *Document) ApplyIncrementalChanges(changes []protocol.TextDocumentContentChangeEvent, version float64) (string, error) { //nolint:lll
	d.mu.RLock()

	if version <= d.doc.Version {
		return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	content := []byte(d.doc.Text)
	uri := d.doc.URI

	d.mu.RUnlock()

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
	}

	return string(content), nil
}

// Set the content after an update send by the client. Must increase the version number
func (d *Document) SetContent(content string, version float64) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if version <= d.doc.Version {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	if len(content) > maxDocumentSize {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/ad.document: Provided.document to large.")
	}

	d.obsoleteVersion()

	d.versionCtx, d.obsoleteVersion = context.WithCancel(context.Background())

	d.doc.Text = content
	d.doc.Version = version
	d.PosData.SetLinesForContent([]byte(content))

	return nil
}

// GetContent returns the content of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
func (d *Document) GetContent(ctx context.Context) (string, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return d.doc.Text, nil
	}
}

// GetCompileResult returns the Compilation Results of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
// It blocks until all compile tasks are finished
func (d *Document) GetCompileResult(ctx context.Context) (*CompiledQuery, error) {
	d.Compilers.Wait()

	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return d.compileResult, nil
	}
}

// GetVersion returns the content of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
func (d *Document) GetVersion(ctx context.Context) (float64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return d.doc.Version, nil
	}
}

// GetURI returns the content of a document
// Since the URI never changes, it does not block or return errors
func (d *Document) GetURI() string {
	return d.doc.URI
}

// GetLanguageID returns the content of a document
// Since the URI never changes, it does not block or return errors
func (d *Document) GetLanguageID() string {
	return d.doc.LanguageID
}
