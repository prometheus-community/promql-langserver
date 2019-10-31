package cache

import "github.com/slrtbtfs/prometheus/promql"

// Updates the compilation Results of a Document. Returns true if the Results were still recent
func (d *Document) UpdateCompileData(version float64, ast promql.Node, err *promql.ParseErr) bool {
	d.Mu.Lock()
	defer d.Mu.Unlock()

	defer d.Compilers.Done()

	if d.Doc.Version > version {
		return false
	}

	d.CompileResult = CompileResult{ast, err}

	return true
}
