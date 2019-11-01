package cache

import (
	"context"

	"github.com/slrtbtfs/prometheus/promql"
)

// Updates the compilation Results of a Document. Returns true if the Results were still recent
func (d *Document) UpdateCompileData(ctx context.Context, ast promql.Node, err *promql.ParseErr) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	defer d.Compilers.Done()

	select {
	case <-ctx.Done():
		return false
	default:
		d.compileResult = &CompileResult{ast, err}
		return true
	}
}
