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

// Document caches content, metadata and compile results of a document
// All exported access methods should be threadsafe
type Document struct {
	posData *token.File

	uri        string
	languageID string
	version    float64
	content    string

	mu sync.RWMutex

	versionCtx      context.Context
	obsoleteVersion context.CancelFunc

	compileResult *CompiledQuery

	// Wait for this before accessing  compileResults
	compilers sync.WaitGroup
}

// ApplyIncrementalChanges applies giver changes to a given Document Content
func (d *Document) ApplyIncrementalChanges(changes []protocol.TextDocumentContentChangeEvent, version float64) (string, error) { //nolint:lll
	d.mu.RLock()

	if version <= d.version {
		return "", jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	content := []byte(d.content)
	uri := d.uri

	d.mu.RUnlock()

	for _, change := range changes {
		// Update column mapper along with the content.
		converter := span.NewContentConverter(uri, content)
		m := &protocol.ColumnMapper{
			URI:       span.URI(d.uri),
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

// SetContent sets the content of a document
func (d *Document) SetContent(content string, version float64, new bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !new && version <= d.version {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInvalidParams, "Update to file didn't increase version number")
	}

	if len(content) > maxDocumentSize {
		return jsonrpc2.NewErrorf(jsonrpc2.CodeInternalError, "cache/SetContent: Provided.document to large.")
	}

	if !new {
		d.obsoleteVersion()
	}

	d.versionCtx, d.obsoleteVersion = context.WithCancel(context.Background())

	d.content = content
	d.version = version
	d.posData.SetLinesForContent([]byte(content))

	d.compilers.Add(1)

	go d.compile(d.versionCtx)

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
		return d.content, nil
	}
}

// GetCompileResult returns the Compilation Results of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
// It blocks until all compile tasks are finished
func (d *Document) GetCompileResult(ctx context.Context) (*CompiledQuery, error) {
	d.compilers.Wait()

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
		return d.version, nil
	}
}

// GetURI returns the content of a document
// Since the URI never changes, it does not block or return errors
func (d *Document) GetURI() string {
	return d.uri
}

// GetLanguageID returns the content of a document
// Since the URI never changes, it does not block or return errors
func (d *Document) GetLanguageID() string {
	return d.languageID
}
