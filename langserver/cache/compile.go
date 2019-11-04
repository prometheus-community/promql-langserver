package cache

import (
	"context"

	"github.com/slrtbtfs/prometheus/promql"
)

type CompiledQuery struct {
	Ast promql.Node
	Err *promql.ParseErr
}

// Updates the compilation Results of a Document. Returns true if the Results were still recent
func (d *Document) UpdateCompileData(ctx context.Context, ast promql.Node, err *promql.ParseErr) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	defer d.Compilers.Done()

	select {
	case <-ctx.Done():
		return false
	default:
		d.compileResult = &CompiledQuery{ast, err}
		return true
	}
}
