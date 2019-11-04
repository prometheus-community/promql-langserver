package cache

import (
	"context"
	"fmt"
	"os"

	"github.com/slrtbtfs/prometheus/promql"
)

// CompiledQuery stores the results of compiling one query
type CompiledQuery struct {
	Ast promql.Node
	Err *promql.ParseErr
}

func (d *Document) compile(ctx context.Context) {
	content, expired := d.GetContent(ctx)
	if expired != nil {
		return
	}

	file := d.posData

	switch d.GetLanguageID() {
	case "promql":
		ast, err := promql.ParseFile(content, file)

		var parseErr *promql.ParseErr

		var ok bool

		if parseErr, ok = err.(*promql.ParseErr); !ok {
			parseErr = nil
		}

		d.AddCompileResult(ctx, ast, parseErr)
	default:
		d.AddCompileResult(ctx, nil, nil)
	}
}

// AddCompileResult updates the compilation Results of a Document. Discards the Result if the context is expired
func (d *Document) AddCompileResult(ctx context.Context, ast promql.Node, err *promql.ParseErr) {
	d.mu.Lock()
	defer d.mu.Unlock()

	defer d.compilers.Done()

	select {
	case <-ctx.Done():
		fmt.Fprint(os.Stderr, "Context expired\n")
	default:
		d.compileResult = &CompiledQuery{ast, err}
		fmt.Fprintf(os.Stderr, "Added compileResult: %v\n", d.compileResult)
	}
}
