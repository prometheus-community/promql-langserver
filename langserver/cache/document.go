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
	Doc     *protocol.TextDocumentItem
	Mu      sync.RWMutex

	versionCtx      context.Context
	obsoleteVersion context.CancelFunc

	CompileResult *CompileResult

	// Wait for this before accessing  compileResults
	Compilers sync.WaitGroup
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

// GetContent returns the content of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
func (d *Document) GetContent(ctx context.Context) (string, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return d.Doc.Text, nil
	}
}

// GetCompileResult returns the Compilation Results of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
// It blocks until all compile tasks are finished
func (d *Document) GetCompileResult(ctx context.Context) (*CompileResult, error) {
	d.Compilers.Wait()

	d.Mu.RLock()
	defer d.Mu.RUnlock()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return d.CompileResult, nil
	}
}

// GetVersion returns the content of a document
// It expects a context.Context retrieved using cache.GetDocument
// and returns an error if that context has expired, i.e. the Document
// has changed since
func (d *Document) GetVersion(ctx context.Context) (float64, error) {
	d.Mu.RLock()
	defer d.Mu.RUnlock()

	select {
	case <-ctx.Done():
		return 0, ctx.Err()
	default:
		return d.Doc.Version, nil
	}
}
